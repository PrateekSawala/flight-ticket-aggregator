package main

import (
	"testing"

	"flight-ticket-aggregator/domain"
	"github.com/stretchr/testify/assert"
)

func TestIsRecordsValidationSuccess(t *testing.T) {
	// Declare test record
	testRecord := domain.FlightRecord{
		FirstName:     "Abh",
		LastName:      "Kuj",
		PNR:           "ABC123",
		FareClass:     "F",
		TravelDate:    "2019-07-31",
		Pax:           "2",
		TicketingDate: "2019-05-21",
		Email:         "ab@zzz.com",
		MobilePhone:   "9876543210",
		BookedCabin:   "Economy",
	}
	err := IsRecordValid(nil, testRecord)
	assert.NoError(t, err)
}

func TestIsRecordsValidationError(t *testing.T) {
	// Declare test record
	testRecord := domain.FlightRecord{
		FirstName:     "Abh",
		LastName:      "Kuj",
		PNR:           "ABC123",
		FareClass:     "F",
		TravelDate:    "2019-07-31",
		Pax:           "2",
		TicketingDate: "2019-05-21",
		Email:         "ab@zzz",
		MobilePhone:   "9876543210",
		BookedCabin:   "Economy",
	}
	err := IsRecordValid(nil, testRecord)
	assert.EqualError(t, err, domain.ErrInvalidMail.Error())
}

func TestIsUploadedFlightRecordValidSuccess(t *testing.T) {
	filename := "airline1_2020-10-30_flightRecord.csv"

	response, err := IsUploadedFlightRecordValid(filename)

	expectedResponse := &domain.FightRecordInfo{AirlineName: "airline1", Filename: "flightRecord.csv", Filepath: "flightRecord/airline1/2020/10/30"}
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, response)
}
