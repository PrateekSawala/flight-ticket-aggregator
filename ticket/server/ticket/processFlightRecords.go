package main

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/PrateekSawala/flight-ticket-aggregator/domain"
	"github.com/PrateekSawala/flight-ticket-aggregator/domain/logging"
	"github.com/PrateekSawala/flight-ticket-aggregator/domain/system"
	"github.com/PrateekSawala/flight-ticket-aggregator/mail/rpc/mail"
	"github.com/PrateekSawala/flight-ticket-aggregator/space/rpc/space"
	"github.com/PrateekSawala/flight-ticket-aggregator/ticket/rpc/ticket"
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
	flightRecordInfoResp, err := system.IsUploadedFlightRecordValid(input.Filename)
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

	passedRecord, failedRecord := &bytes.Buffer{}, &bytes.Buffer{}
	passedRecordWriter, failedRecordWriter := &csv.Writer{}, &csv.Writer{}
	passedRecordCreated, failedRecordCreated := false, false
	recordedPNRs := map[string]bool{}

	recordName := strings.TrimSuffix(flightRecordInfo.Filename, filepath.Ext(flightRecordInfo.Filename))
	passedRecordFileName := fmt.Sprintf("%s_passedRecord.csv", recordName)
	failedRecordFileName := fmt.Sprintf("%s_failedRecord.csv", recordName)

	for i := range records {
		recordEntry := records[i]
		if i == 0 {
			err := IsRecordsColumnValid(records[i])
			if err != nil {
				return response, err
			}
			continue
		}
		flightRecord := PrepareFlightRecord(recordEntry)
		recordErr := IsRecordValid(recordedPNRs, flightRecord)
		if recordErr == nil {
			if !passedRecordCreated {
				passedRecordWriterResp, passedRecordError := InitWriter(passedRecord, domain.TypeRecordPassed)
				if passedRecordError != nil {
					return response, passedRecordError
				}
				passedRecordWriter = passedRecordWriterResp
				passedRecordCreated = true
			}
			if offerCode, passedRecordError := GetDiscountCode(flightRecord.FareClass); passedRecordError == nil {
				passedRecordError := WriteRecord(passedRecordWriter, append(recordEntry, offerCode))
				if passedRecordError != nil {
					log.Errorf("Error while writing flight passed record of file %s,  %+v , error: %s", passedRecordFileName, recordEntry, passedRecordError)
				}
				recordedPNRs[flightRecord.PNR] = true
			}
		} else {
			if !failedRecordCreated {
				failedRecordWriterResp, failedRecordError := InitWriter(failedRecord, domain.TypeRecordFailed)
				if failedRecordError != nil {
					return response, failedRecordError
				}
				failedRecordWriter = failedRecordWriterResp
				failedRecordCreated = true
			}
			err := WriteRecord(failedRecordWriter, append(recordEntry, recordErr.Error()))
			if err != nil {
				log.Errorf("Error while writing flight failed record of file %s, %+v , error: %s", failedRecordFileName, recordEntry, err)
			}
			recordedPNRs[flightRecord.PNR] = true
		}
	}
	if passedRecordCreated {
		response.PassedRecordFileName = passedRecordFileName
		passedRecordWriter.Flush()
		_, err := spaceService.SaveFile(context.Background(), &space.SaveFileInput{Filename: passedRecordFileName, Filepath: flightRecordInfo.Filepath, File: passedRecord.Bytes()})
		if err != nil {
			log.Errorf("Error while saving file to s3 path: %s/%s, error: %s", flightRecordInfo.Filepath, passedRecordFileName, err)
		} else {
			response.PassedRecordFilePathUrl = fmt.Sprintf("%s/%s/%s/%s", *webServerHostName, domain.DownloadFlightRecordPath, flightRecordInfo.Filepath, passedRecordFileName)
		}
	}
	if failedRecordCreated {
		response.FailedRecordFileName = failedRecordFileName
		failedRecordWriter.Flush()
		_, err := spaceService.SaveFile(context.Background(), &space.SaveFileInput{Filename: failedRecordFileName, Filepath: flightRecordInfo.Filepath, File: failedRecord.Bytes()})
		if err != nil {
			log.Errorf("Error while saving file to s3 path: %s/%s, error: %s", flightRecordInfo.Filepath, failedRecordFileName, err)
		} else {
			response.FailedRecordFilePathUrl = fmt.Sprintf("%s/%s/%s/%s", *webServerHostName, domain.DownloadFlightRecordPath, flightRecordInfo.Filepath, failedRecordFileName)
		}
	}
	return response, nil
}

func InitWriter(buffer *bytes.Buffer, typeRecord string) (*csv.Writer, error) {
	if typeRecord != domain.TypeRecordPassed && typeRecord != domain.TypeRecordFailed {
		return nil, domain.ErrInvalidInput
	}
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

func GetDiscountCode(class string) (string, error) {
	if class == "" {
		return domain.Empty, domain.ErrMissingArgument
	}
	if !system.IsAlphabetic(class) {
		return domain.Empty, domain.ErrInvalidArgument
	}
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
	return offerCode, nil
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
