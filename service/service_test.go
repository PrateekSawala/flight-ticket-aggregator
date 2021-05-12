package service

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"

	"flight-ticket-aggregator/domain"
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

func TestUploadFlightRecord(t *testing.T) {
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
	records, err := UploadFlightRecord(documentName, fileBuffer)
	if err != nil {
		t.Errorf("Error while uploading flight record %s, error: %s", filename, err)
		return
	}
	UploadResponseRecords = records
}

func TestComparePassedRecordWithExpectedPassedRecord(t *testing.T) {
	ExpectedFlightPassedRecordPath := fmt.Sprintf("../%s/%s", domain.TemplateFolder, ExpectedFlightPassedRecordName)
	expected, err := ioutil.ReadFile(ExpectedFlightPassedRecordPath)
	if err != nil {
		t.Errorf("ioutil.ReadFile error: %s", err)
		return
	}

	if UploadResponseRecords.PassedRecordFileName == domain.Empty {
		t.Errorf("No passed record uploaded response exist")
		return
	}

	ResponsePassedRecordPath := fmt.Sprintf("../%s/%s", domain.UploadFolder, UploadResponseRecords.PassedRecordFileName)
	found, err := ioutil.ReadFile(ResponsePassedRecordPath)
	if err != nil {
		t.Errorf("ioutil.ReadFile error: %s", err)
		return
	}

	if !bytes.Equal(expected, found) {
		t.Errorf("Expected Passed record %s not matched with found passed record %s", ExpectedFlightPassedRecordName, UploadResponseRecords.PassedRecordFileName)
	}
}

func TestCompareFailedRecordWithExpectedFailedRecord(t *testing.T) {
	ExpectedFlightFailedRecordPath := fmt.Sprintf("../%s/%s", domain.TemplateFolder, ExpectedFlightFailedRecordName)
	expected, err := ioutil.ReadFile(ExpectedFlightFailedRecordPath)
	if err != nil {
		t.Errorf("ioutil.ReadFile error: %s", err)
		return
	}

	if UploadResponseRecords.FailedRecordFileName == domain.Empty {
		t.Errorf("No failed record uploaded response exist")
		return
	}

	ResponseFailedRecordPath := fmt.Sprintf("../%s/%s", domain.UploadFolder, UploadResponseRecords.FailedRecordFileName)
	found, err := ioutil.ReadFile(ResponseFailedRecordPath)
	if err != nil {
		t.Errorf("ioutil.ReadFile error: %s", err)
		return
	}

	if !bytes.Equal(expected, found) {
		t.Errorf("Expected failed record %s not matched with found failed record %s", ExpectedFlightFailedRecordName, UploadResponseRecords.FailedRecordFileName)
	}
}

func TestUploadFlightRecordFailure(t *testing.T) {
	// Open the test file
	filename := TestFailedFlightRecord
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
}

func TestUploadEmptyFlightRecordFailure(t *testing.T) {
	// Open the test file
	filename := TestEmptyFlightRecord
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
}
