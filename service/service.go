package service

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"flight-ticket-aggregator/domain"
	"flight-ticket-aggregator/domain/logging"
	"flight-ticket-aggregator/utility"
)

var (
	PNRs map[string]bool
)

func UploadFlightRecord(filename string, flightRecord []byte) (*domain.Record, error) {
	log := logging.Log("UploadFlightRecord")
	log.Tracef("Processing file %s", filename)
	defer log.Tracef("Processed file %s", filename)

	// Declare return response
	response := &domain.Record{}

	currentTime := time.Now().UnixNano()

	// Prepare records
	recordName := strings.TrimSuffix(filename, filepath.Ext(filename))
	passedRecordFileName := fmt.Sprintf("%s_%d_passedRecord.csv", recordName, currentTime)
	failedRecordFileName := fmt.Sprintf("%s_%d_failedRecord.csv", recordName, currentTime)

	// Parsing the file
	reader := csv.NewReader(bytes.NewReader(flightRecord))
	reader.LazyQuotes = true
	reader.Comma = ','

	records, err := reader.ReadAll()
	if err != nil {
		return response, err
	} else if len(records) == 0 {
		return response, domain.ErrEmptyFile
	}

	response, err = ProcessRecords(records, passedRecordFileName, failedRecordFileName)
	if err != nil {
		log.Errorf("Error while processing records: %s", err)
		return response, err
	}

	return response, nil
}

func ProcessRecords(records [][]string, passedRecordName string, failedRecordName string) (*domain.Record, error) {
	log := logging.Log("ProcessRecords")

	passedRecord := &os.File{}
	failedRecord := &os.File{}
	passedRecordWriter := &csv.Writer{}
	failedRecordWriter := &csv.Writer{}
	passedRecordError := domain.ErrPassedRecord
	failedRecordError := domain.ErrFailedRecord
	passedRecordCreated, failedRecordCreated := false, false
	PNRs = map[string]bool{}

	// Declare return response
	response := &domain.Record{}

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

		flightRecord := PrepareFlightRecord(recordEntry)
		recordErr := IsRecordValid(flightRecord)
		if recordErr == nil {
			if !passedRecordCreated {
				passedRecord, passedRecordWriter, passedRecordError = InitWriter(passedRecordName, domain.TypeRecordPassed)
				if passedRecordError != nil {
					return response, passedRecordError
				}
				passedRecordCreated = true
			}
			if offerCode, passedRecordError := GetDiscountCode(flightRecord.FareClass); passedRecordError == nil {
				passedRecordError = WriteRecord(passedRecordWriter, append(recordEntry, offerCode))
				if passedRecordError == nil {
					PNRs[flightRecord.PNR] = true
					continue
				}
				log.Errorf("Error while writing flight passed record of file %s,  %+v , error: %s", passedRecordName, recordEntry, passedRecordError)
			}
			recordErr = passedRecordError
		}
		if !failedRecordCreated {
			failedRecord, failedRecordWriter, failedRecordError = InitWriter(failedRecordName, domain.TypeRecordFailed)
			if failedRecordError != nil {
				return response, failedRecordError
			}
			failedRecordCreated = true
		}
		err := WriteRecord(failedRecordWriter, append(recordEntry, recordErr.Error()))
		if err != nil {
			log.Errorf("Error while writing flight failed record of file %s, %+v , error: %s", failedRecordName, recordEntry, err)
		}
	}

	// Release resources
	if passedRecordCreated {
		releaseResources(passedRecord, passedRecordWriter)
		response.PassedRecordFileName = passedRecordName
	}
	if failedRecordCreated {
		releaseResources(failedRecord, failedRecordWriter)
		response.FailedRecordFileName = failedRecordName
	}
	return response, nil
}

func releaseResources(file *os.File, writer *csv.Writer) {
	go func() {
		writer.Flush()
		file.Close()
	}()
}

func CreateRecords(recordPath string) (*os.File, error) {

	newRecordPath := fmt.Sprintf("../%s/%s", domain.UploadFolder, recordPath)
	// Create new record file
	record, err := os.Create(newRecordPath)
	if err != nil {
		return nil, err
	}
	return record, nil
}

func InitWriter(recordPath string, typeRecord string) (*os.File, *csv.Writer, error) {
	record, err := CreateRecords(recordPath)
	if err != nil {
		return nil, nil, err
	}

	recordWriter := csv.NewWriter(record)
	newEntry := domain.Error
	if typeRecord == domain.TypeRecordPassed {
		newEntry = domain.DiscountCode
	}

	err = WriteRecord(recordWriter, append(domain.RecordColumns, newEntry))
	if err != nil {
		return nil, nil, err
	}
	return record, recordWriter, nil
}

func WriteRecord(writer *csv.Writer, recordEntry []string) error {
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

	if !utility.IsAlphabetic(class) {
		return domain.Empty, domain.ErrInvalidArgument
	}

	code := strings.ToUpper(class)
	asciiValue := []rune(code)[0]

	offerCode := domain.Empty
	switch {
	case asciiValue >= 65 && asciiValue <= 69:
		offerCode = domain.DiscountCode20
	case asciiValue >= 70 && asciiValue <= 75:
		offerCode = domain.DiscountCode30
	case asciiValue >= 76 && asciiValue <= 82:
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
