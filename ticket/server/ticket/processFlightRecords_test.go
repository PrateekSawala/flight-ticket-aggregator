package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"testing"

	"flight-ticket-aggregator/domain"
	"flight-ticket-aggregator/ticket/rpc/ticket"

	"github.com/stretchr/testify/assert"
)

var (
	testclient = ticket.NewTicketProtobufClient("http://localhost:3003", &http.Client{})
)

func TestInputFlightRecordSuccess(t *testing.T) {
	// Open the test file
	filename := domain.TestFlightRecord

	importfile := fmt.Sprintf("../../templates/%s", filename)

	// Read file
	fileBuffer, err := ioutil.ReadFile(importfile)
	if err != nil {
		t.Errorf("Error occured while looking for file %s, Error: %s", filename, err)
		return
	}

	// Find document name
	_, documentName := filepath.Split(filename)

	_, err = testclient.ProcessFlightRecord(context.Background(), &ticket.ProcessFlightRecordInput{Filename: documentName, FlightRecord: fileBuffer})
	if err != nil {
		t.Errorf("Error while uploading flight record %s, Error: %s", filename, err)
		return
	}
	assert.NoError(t, err)
}

func TestInputFlightRecordError(t *testing.T) {
	// Open the test file
	filename := domain.TestEmptyFlightRecord
	importfile := fmt.Sprintf("../../templates/%s", filename)

	fileBuffer, err := ioutil.ReadFile(importfile)
	if err != nil {
		t.Errorf("Error occured while looking for file %s, error: %s", filename, err)
		return
	}
	// Find document name
	_, documentName := filepath.Split(filename)

	_, err = testclient.ProcessFlightRecord(context.Background(), &ticket.ProcessFlightRecordInput{Filename: documentName, FlightRecord: fileBuffer})

	expectedErr := fmt.Errorf("twirp error internal: %s", domain.ErrEmptyFile)
	assert.EqualError(t, expectedErr, err.Error())
}
