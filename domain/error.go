package domain

import "errors"

var (
	/* General errors */

	ErrMissingArgument      = errors.New(MissingArgument)
	ErrInvalidArgument      = errors.New(InvalidArgument)
	ErrNotFound             = errors.New(NotFound)
	ErrInternal             = errors.New(InternalError)
	ErrInvalidRequest       = errors.New(InvalidRequest)
	ErrInvalidContentType   = errors.New(InvalidContentType)
	ErrInvalidFile          = errors.New(InvalidFile)
	ErrInvalidFilename      = errors.New(InvalidFilename)
	ErrEmptyFile            = errors.New(EmptyFile)
	ErrFileNotFound         = errors.New(FileNotFound)
	ErrInvalidTypeAssetion  = errors.New(InvalidTypeAssetion)
	ErrInvalidInput         = errors.New(InvalidInput)
	ErrEmptyInput           = errors.New(EmptyInput)
	ErrInvalidAirline       = errors.New(InvalidAirline)
	ErrTryAgainLater        = errors.New(TryAgainLater)
	ErrInvalidRequestMethod = errors.New(InvalidRequestMethod)

	/* Record Validation errors */

	ErrInvalidMail      = errors.New(InvalidMail)
	ErrInvalidPhone     = errors.New(InvalidPhone)
	ErrInvalidCabin     = errors.New(InvalidCabin)
	ErrInvalidPNR       = errors.New(InvalidPNR)
	ErrInvalidFareClass = errors.New(InvalidFareClass)
	ErrTicketingDate    = errors.New(InvalidTicketingDate)
	ErrTravelDate       = errors.New(InvalidTravelDate)
	ErrUploadedDate     = errors.New(InvalidUploadedDate)
	ErrInvalidBooking   = errors.New(InvalidBooking)
	ErrPNRRepeated      = errors.New(PNRRepeated)
)
