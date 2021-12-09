package main

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"path/filepath"
	"strings"

	"flight-ticket-aggregator/domain"
	"flight-ticket-aggregator/domain/logging"
	"flight-ticket-aggregator/mail/rpc/mail"
	"flight-ticket-aggregator/space/rpc/space"
	"flight-ticket-aggregator/ticket/rpc/ticket"
)

func (s *Server) ProcessFlightRecord(ctx context.Context, input *ticket.ProcessFlightRecordInput) (*ticket.ProcessFlightRecordResponse, error) {
	log := logging.Log("ProcessFlightRecord")
	log.Trace("Start")
	defer log.Trace("End")

	log.Tracef("Input: %+v", input)

	response := &ticket.ProcessFlightRecordResponse{}

	if input == nil {
		log.Debugf("Empty input")
		return response, domain.ErrInvalidInput
	}

	// Check if uploaded flightRecord is valid
	flightRecordInfoResp, err := IsUploadedFlightRecordValid(input.Filename)
	if err != nil {
		return response, err
	}

	// Parsing the file
	reader := csv.NewReader(bytes.NewReader(input.FlightRecord))
	reader.LazyQuotes = true
	reader.Comma = ','

	// Read all records
	records, err := reader.ReadAll()
	if err != nil {
		return response, err
	} else if len(records) == 0 {
		return response, domain.ErrEmptyFile
	}

	_, err = spaceService.SaveFile(context.Background(), &space.SaveFileInput{Filename: flightRecordInfoResp.Filename, Filepath: flightRecordInfoResp.Filepath, File: input.FlightRecord})
	if err != nil {
		log.Errorf("spaceService.SaveFile error : %s", err)
		return response, err
	}

	response, err = ProcessRecords(records, flightRecordInfoResp)
	if err != nil {
		log.Errorf("Error while processing records: %s", err)
		return response, err
	}
	response.AirlineName = flightRecordInfoResp.AirlineName

	if response.PassedRecordFileName == domain.Empty && response.FailedRecordFileName == domain.Empty {
		log.Infof("No passed or failed record generated from the uploaded file %s, This shouldn't be the case", flightRecordInfoResp.Filename)
		return response, err
	}

	go func() {
		_, err := mailService.SendProcessedFlightRecordsMail(context.Background(), &mail.SendProcessedFlightRecordsMailInput{UploadedFileName: flightRecordInfoResp.Filename, UploadedFilePath: flightRecordInfoResp.Filepath, Processedfiles: []string{response.PassedRecordFileName, response.FailedRecordFileName}, AirlineName: flightRecordInfoResp.AirlineName})
		if err != nil {
			log.Errorf("Error while sending flightRecord %s/%s mail to %s airline, error: %s", flightRecordInfoResp.Filepath, flightRecordInfoResp.Filename, flightRecordInfoResp.AirlineName, err)
		}
	}()
	return response, nil
}

func ProcessRecords(records [][]string, flightRecordInfo *domain.FightRecordInfo) (*ticket.ProcessFlightRecordResponse, error) {
	log := logging.Log("ProcessRecords")
	response := &ticket.ProcessFlightRecordResponse{}

	flightRecordName := strings.TrimSuffix(flightRecordInfo.Filename, filepath.Ext(flightRecordInfo.Filename))

	// Declare flightRecord writers
	flightRecordWriters := map[string]*domain.FlightRecordWriter{}
	flightRecordWriters[domain.TypeRecordPassed] = &domain.FlightRecordWriter{Filename: fmt.Sprintf("%s_passedRecord.csv", flightRecordName), Filepath: flightRecordInfo.Filepath, TypeRecord: domain.TypeRecordPassed, File: &bytes.Buffer{}, FileWriter: &csv.Writer{}}
	flightRecordWriters[domain.TypeRecordFailed] = &domain.FlightRecordWriter{Filename: fmt.Sprintf("%s_failedRecord.csv", flightRecordName), Filepath: flightRecordInfo.Filepath, TypeRecord: domain.TypeRecordFailed, File: &bytes.Buffer{}, FileWriter: &csv.Writer{}}

	recordedPNRs := map[string]bool{}
	// Loop over all records
	for i := range records {
		recordEntry := records[i]
		if i == 0 {
			err := IsRecordsColumnValid(records[i])
			if err != nil {
				return response, err
			}
			continue
		}
		flightRecordEntry := PrepareFlightRecord(recordEntry)
		flightRecordWriter := flightRecordWriters[domain.TypeRecordPassed]
		recordErr := IsRecordValid(recordedPNRs, flightRecordEntry)
		if recordErr != nil {
			recordEntry = append(recordEntry, recordErr.Error())
			flightRecordWriter = flightRecordWriters[domain.TypeRecordFailed]
		} else {
			recordEntry = append(recordEntry, GetDiscountCode(flightRecordEntry.FareClass))
		}
		err := WriteFlightRecord(flightRecordWriter, recordEntry)
		if err != nil {
			log.Errorf("Error while writing %s flight record of filename %s with recordEntry %+v , error: %s", flightRecordWriter.TypeRecord, flightRecordWriter.Filename, recordEntry, err)
		}
		recordedPNRs[flightRecordEntry.PNR] = true
	}
	SaveFlightRecord(flightRecordWriters, response, log)
	return response, nil
}

func WriteFlightRecord(flightRecordWriter *domain.FlightRecordWriter, recordEntry []string) error {
	if flightRecordWriter.TypeRecord != domain.TypeRecordPassed && flightRecordWriter.TypeRecord != domain.TypeRecordFailed {
		return domain.ErrInvalidInput
	}
	if !flightRecordWriter.IsFileCreated {
		recordWriterResp, err := InitWriter(flightRecordWriter.File, domain.TypeRecordPassed)
		if err != nil {
			return err
		}
		flightRecordWriter.FileWriter = recordWriterResp
		flightRecordWriter.IsFileCreated = true
	}
	return WriteRecord(flightRecordWriter.FileWriter, recordEntry)
}

func InitWriter(buffer *bytes.Buffer, typeRecord string) (*csv.Writer, error) {
	recordWriter := csv.NewWriter(buffer)
	newEntry := domain.RecordError
	if typeRecord == domain.TypeRecordPassed {
		newEntry = domain.DiscountCode
	}
	err := WriteRecord(recordWriter, append(domain.RecordColumns, newEntry))
	if err != nil {
		return nil, err
	}
	return recordWriter, nil
}

func WriteRecord(writer *csv.Writer, recordEntry []string) error {
	if len(recordEntry) == 0 {
		return domain.ErrInvalidInput
	}
	err := writer.Write(recordEntry)
	if err != nil {
		return err
	}
	return nil
}

func SaveFlightRecord(flightRecordWriters map[string]*domain.FlightRecordWriter, response *ticket.ProcessFlightRecordResponse, log *logging.Logger) {
	for _, flightRecordWriter := range flightRecordWriters {
		if !flightRecordWriter.IsFileCreated {
			return
		}
		flightRecordWriter.FileWriter.Flush()
		spacePathUrl := ""
		_, err := spaceService.SaveFile(context.Background(), &space.SaveFileInput{Filename: flightRecordWriter.Filename, Filepath: flightRecordWriter.Filepath, File: flightRecordWriter.File.Bytes()})
		if err != nil {
			log.Errorf("Error while saving file to s3 path: %s/%s, error: %s", flightRecordWriter.Filepath, flightRecordWriter.Filename, err)
		} else {
			spacePathUrl = fmt.Sprintf("%s/%s/%s/%s", *webServerHostName, domain.DownloadFlightRecordPath, flightRecordWriter.Filepath, flightRecordWriter.Filename)
		}
		if flightRecordWriter.TypeRecord == domain.TypeRecordPassed {
			response.PassedRecordFileName = flightRecordWriter.Filename
			response.PassedRecordFilePathUrl = spacePathUrl
		}
		if flightRecordWriter.TypeRecord == domain.TypeRecordFailed {
			response.FailedRecordFileName = flightRecordWriter.Filename
			response.FailedRecordFilePathUrl = spacePathUrl
		}
	}
}

func GetDiscountCode(class string) string {
	code := strings.ToUpper(class)
	asciiValue := []rune(code)[0]
	offerCode := domain.Empty
	switch {
	case asciiValue >= domain.FareClassA && asciiValue <= domain.FareClassB:
		offerCode = domain.DiscountCode20
	case asciiValue >= domain.FareClassF && asciiValue <= domain.FareClassK:
		offerCode = domain.DiscountCode30
	case asciiValue >= domain.FareClassL && asciiValue <= domain.FareClassR:
		offerCode = domain.DiscountCode25
	}
	return offerCode
}

func PrepareFlightRecord(recordEntry []string) domain.FlightRecord {
	return domain.FlightRecord{
		FirstName:     recordEntry[0],
		LastName:      recordEntry[1],
		PNR:           recordEntry[2],
		FareClass:     recordEntry[3],
		TravelDate:    recordEntry[4],
		Pax:           recordEntry[5],
		TicketingDate: recordEntry[6],
		Email:         recordEntry[7],
		MobilePhone:   recordEntry[8],
		BookedCabin:   recordEntry[9],
	}
}
