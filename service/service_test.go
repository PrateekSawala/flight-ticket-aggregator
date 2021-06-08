package service

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"

	"flight-ticket-aggregator/domain"
	"github.com/stretchr/testify/assert"
)

var (
	UploadResponseRecords *domain.Record
)

const (
	TestFlightRecord               = "flightRecord.csv"
	TestFailedFlightRecord         = "testFailedFlightRecord.csv"
	TestEmptyFlightRecord          = "testEmptyFlightRecord.csv"
	ExpectedFlightPassedRecordName = "expectedFlightPassedRecords.csv"
	ExpectedFlightFailedRecordName = "expectedFlightFailedRecords.csv"
)

func TestInputFlightRecordSuccess(t *testing.T) {
	// Open the test file
	filename := TestFlightRecord
	importfile := fmt.Sprintf("../%s/%s", domain.TemplateFolder, filename)
	uploadFolderPath = fmt.Sprintf("../%s", domain.UploadFolder)

	fileBuffer, err := ioutil.ReadFile(importfile)
	if err != nil {
		t.Errorf("Error occured while looking for file %s, error: %s", filename, err)
	}

	// Find document name
	_, documentName := filepath.Split(filename)
	_, err = UploadFlightRecord(documentName, fileBuffer)
	if err != nil {
		t.Errorf("Error while uploading flight record %s, error: %s", filename, err)
		return
	}
	assert.NoError(t, err)
}

func TestInputFlightRecordFailure(t *testing.T) {
	// Upload folder path
	uploadFolderPath = fmt.Sprintf("../%s", domain.UploadFolder)

	t.Run("Should return error of invalid file", func(t *testing.T) {
		// Open the test file
		filename := TestFailedFlightRecord
		importfile := fmt.Sprintf("../%s/%s", domain.TemplateFolder, filename)

		fileBuffer, err := ioutil.ReadFile(importfile)
		if err != nil {
			t.Errorf("Error occured while looking for file %s, error: %s", filename, err)
		}
		// Find document name
		_, documentName := filepath.Split(filename)
		_, err = UploadFlightRecord(documentName, fileBuffer)
		assert.EqualError(t, err, domain.ErrInvalidFile.Error())
	})

	t.Run("Should return error of empty file", func(t *testing.T) {
		// Open the test file
		filename := TestEmptyFlightRecord
		importfile := fmt.Sprintf("../%s/%s", domain.TemplateFolder, filename)

		fileBuffer, err := ioutil.ReadFile(importfile)
		if err != nil {
			t.Errorf("Error occured while looking for file %s, error: %s", filename, err)
		}
		// Find document name
		_, documentName := filepath.Split(filename)
		_, err = UploadFlightRecord(documentName, fileBuffer)
		assert.EqualError(t, err, domain.ErrEmptyFile.Error())
	})
}
