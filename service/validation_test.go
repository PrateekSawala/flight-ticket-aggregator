package service

import (
	"flight-ticket-aggregator/domain"
	"testing"
)

func TestIsRecordsColumnValid(t *testing.T) {
	err := IsRecordsColumnValid(domain.RecordColumns)
	if err != nil {
		t.Errorf("Validation error %s", err)
	}
}

func TestIsRecordValid(t *testing.T) {
	record := domain.FlightRecord{
		FirstName:     "Abhishek",
		LastName:      "Kumar",
		PNR:           "ABC123",
		FareClass:     "F",
		TravelDate:    "2019-07-31",
		Pax:           "2",
		TicketingDate: "2019-05-21",
		Email:         "abhishek@zzz.com",
		MobilePhone:   "9876543210",
		BookedCabin:   "Economy",
	}
	err := IsRecordValid(record)
	if err != nil {
		t.Errorf("Validation error %s", err)
	}
}

func TestIsRecordsColumnValidationFailure(t *testing.T) {
	record := domain.RecordColumns
	record[0] = "FirstName"
	err := IsRecordsColumnValid(record)
	if err != nil {
		t.Errorf("Validation error %s", err)
	}
}

func TestIsRecordsColumnIndexesValidationFailure(t *testing.T) {
	record := []string{}
	err := IsRecordsColumnValid(record)
	if err != nil {
		t.Errorf("Validation error %s", err)
	}
}

func TestIsRecordPNRValidationFailure(t *testing.T) {
	record := domain.FlightRecord{
		FirstName:     "Abhishek",
		LastName:      "Kumar",
		PNR:           "ABC",
		FareClass:     "F",
		TravelDate:    "2019-07-31",
		Pax:           "2",
		TicketingDate: "2019-05-21",
		Email:         "abhishek@zzz.com",
		MobilePhone:   "9876543210",
		BookedCabin:   "Economy",
	}
	err := IsRecordValid(record)
	if err != nil {
		t.Errorf("Validation error %s", err)
	}
}

func TestIsRecordMailValidationFailure(t *testing.T) {
	record := domain.FlightRecord{
		FirstName:     "Abhishek",
		LastName:      "Kumar",
		PNR:           "ABC123",
		FareClass:     "F",
		TravelDate:    "2019-07-31",
		Pax:           "2",
		TicketingDate: "2019-05-21",
		Email:         "abhishek@zzz",
		MobilePhone:   "9876543210",
		BookedCabin:   "Economy",
	}
	err := IsRecordValid(record)
	if err != nil {
		t.Errorf("Validation error %s", err)
	}
}

func TestIsRecordPhoneValidationFailure(t *testing.T) {
	record := domain.FlightRecord{
		FirstName:     "Abhishek",
		LastName:      "Kumar",
		PNR:           "ABC123",
		FareClass:     "F",
		TravelDate:    "2019-07-31",
		Pax:           "2",
		TicketingDate: "2019-05-21",
		Email:         "abhishek@zzz.com",
		MobilePhone:   "000111",
		BookedCabin:   "Economy",
	}
	err := IsRecordValid(record)
	if err != nil {
		t.Errorf("Validation error %s", err)
	}
}

func TestIsRecordTicketingDateBeforeTravelDateValidationFailure(t *testing.T) {
	record := domain.FlightRecord{
		FirstName:     "Abhishek",
		LastName:      "Kumar",
		PNR:           "ABC123",
		FareClass:     "F",
		TravelDate:    "2019-09-21",
		Pax:           "2",
		TicketingDate: "2019-09-25",
		Email:         "abhishek@zzz.com",
		MobilePhone:   "9876543210",
		BookedCabin:   "Economy",
	}
	err := IsRecordValid(record)
	if err != nil {
		t.Errorf("Validation error %s", err)
	}
}

func TestIsRecordCabinValidValidationFailure(t *testing.T) {
	record := domain.FlightRecord{
		FirstName:     "Abhishek",
		LastName:      "Kumar",
		PNR:           "ABC123",
		FareClass:     "F",
		TravelDate:    "2019-07-31",
		Pax:           "2",
		TicketingDate: "2019-05-21",
		Email:         "abhishek@zzz.com",
		MobilePhone:   "9876543210",
		BookedCabin:   "Basic Economy",
	}
	err := IsRecordValid(record)
	if err != nil {
		t.Errorf("Validation error %s", err)
	}
}
