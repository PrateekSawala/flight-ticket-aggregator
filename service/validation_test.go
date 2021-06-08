package service

import (
	"testing"

	"flight-ticket-aggregator/domain"
	"github.com/stretchr/testify/assert"
)

func TestIsRecordsSuccessfulValidation(t *testing.T) {
	// Declare test record
	testRecord := domain.FlightRecord{
		FirstName:     "Ab",
		LastName:      "Kumar",
		PNR:           "ABC123",
		FareClass:     "F",
		TravelDate:    "2019-07-31",
		Pax:           "2",
		TicketingDate: "2019-05-21",
		Email:         "ab@zzz.com",
		MobilePhone:   "9876543210",
		BookedCabin:   "Economy",
	}
	err := IsRecordValid(testRecord)
	assert.NoError(t, err)
}

func TestIsRecordsValidationFailure(t *testing.T) {

	t.Run("Should return error of invalid email", func(t *testing.T) {
		// Declare test record
		testRecord := domain.FlightRecord{
			FirstName:     "Ab",
			LastName:      "Kumar",
			PNR:           "ABC123",
			FareClass:     "F",
			TravelDate:    "2019-07-31",
			Pax:           "2",
			TicketingDate: "2019-05-21",
			Email:         "ab@zzz",
			MobilePhone:   "9876543210",
			BookedCabin:   "Economy",
		}
		err := IsRecordValid(testRecord)
		assert.EqualError(t, err, domain.ErrInvalidMail.Error())
	})

	t.Run("Should return error of invalid phone", func(t *testing.T) {
		// Declare testRecord
		testRecord := domain.FlightRecord{
			FirstName:     "Ab",
			LastName:      "Kumar",
			PNR:           "ABC123",
			FareClass:     "F",
			TravelDate:    "2019-07-31",
			Pax:           "2",
			TicketingDate: "2019-05-21",
			Email:         "ab@zzz.com",
			MobilePhone:   "000111",
			BookedCabin:   "Economy",
		}
		err := IsRecordValid(testRecord)
		assert.EqualError(t, err, domain.ErrInvalidPhone.Error())
	})

	t.Run("Should return error of invalid booking", func(t *testing.T) {
		// Declare testRecord
		testRecord := domain.FlightRecord{
			FirstName:     "Ab",
			LastName:      "Kumar",
			PNR:           "ABC123",
			FareClass:     "F",
			TravelDate:    "2019-05-21",
			Pax:           "2",
			TicketingDate: "2019-05-22",
			Email:         "ab@zzz.com",
			MobilePhone:   "9876543210",
			BookedCabin:   "Economy",
		}
		err := IsRecordValid(testRecord)
		assert.EqualError(t, err, domain.ErrInvalidBooking.Error())
	})

	t.Run("Should return error of invalid PNR", func(t *testing.T) {
		// Declare testRecord
		testRecord := domain.FlightRecord{
			FirstName:     "Ab",
			LastName:      "Kumar",
			PNR:           "ABC",
			FareClass:     "F",
			TravelDate:    "2019-07-31",
			Pax:           "2",
			TicketingDate: "2019-05-21",
			Email:         "ab@zzz.com",
			MobilePhone:   "9876543210",
			BookedCabin:   "Economy",
		}
		err := IsRecordValid(testRecord)
		assert.EqualError(t, err, domain.ErrInvalidPNR.Error())
	})

	t.Run("Should return error of invalid cabin", func(t *testing.T) {
		// Declare testRecord
		testRecord := domain.FlightRecord{
			FirstName:     "Ab",
			LastName:      "Kumar",
			PNR:           "ABC123",
			FareClass:     "F",
			TravelDate:    "2019-07-31",
			Pax:           "2",
			TicketingDate: "2019-05-21",
			Email:         "ab@zzz.com",
			MobilePhone:   "9876543210",
			BookedCabin:   "Basic Economy",
		}
		err := IsRecordValid(testRecord)
		assert.EqualError(t, err, domain.ErrInvalidCabin.Error())
	})
}
