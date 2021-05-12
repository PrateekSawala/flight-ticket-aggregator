package domain

import "errors"

var (
	// General errors
	ErrMissingArgument    = errors.New("missing_argument")
	ErrInvalidArgument    = errors.New("invalid_argument")
	ErrNotFound           = errors.New("not_found")
	ErrInternal           = errors.New("internal_error")
	ErrInvalidRequest     = errors.New("invalid_request")
	ErrInvalidContentType = errors.New("invalid_content-type")
	ErrInvalidFile        = errors.New("invalid_file")
	ErrEmptyFile          = errors.New("empty_file")
	ErrTypeAssetion       = errors.New("invalid_type_assetion")
	ErrInvalidInput       = errors.New("invalid_input")

	// Validation errors
	ErrInvalidMail    = errors.New("Email invalid")
	ErrInvalidPhone   = errors.New("Phone invalid")
	ErrInvalidCabin   = errors.New("Cabin invalid")
	ErrInvalidPNR     = errors.New("PNR invalid")
	ErrTicketingDate  = errors.New("Ticketing Date invalid")
	ErrTravelDate     = errors.New("Travel Date invalid")
	ErrInvalidBooking = errors.New("Booking invalid")
	ErrPNRRepeated    = errors.New("PNR Must Be Unique")

	// Records errors
	ErrPassedRecord = errors.New("")
	ErrFailedRecord = errors.New("")
)
