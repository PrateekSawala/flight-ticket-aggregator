package domain

import (
	"regexp"
)

var (
	ValidFlightCabins = map[string]bool{
		Economy:        true,
		PremiumEconomy: true,
		Business:       true,
		First:          true,
	}

	ValidAirlines = map[string]bool{
		Airline1: true,
		Airline2: true,
		Airline3: true,
	}

	AirlinesMails = map[string]string{
		Airline1: "airline1@mail.com",
		Airline2: "airline2@mail.com",
		Airline3: "airline3@mail.com",
	}

	RecordColumns = []string{"First_name", "Last_name", "PNR", "Fare_class", "Travel_date", "Pax", "Ticketing_date", "Email", "Mobile_phone", "Booked_cabin"}

	ValidRecordIndexes = map[int]string{
		0: FirstName,
		1: LastName,
		2: PNR,
		3: FareClass,
		4: TravelDate,
		5: Pax,
		6: TicketingDate,
		7: Email,
		8: MobilePhone,
		9: BookedCabin,
	}

	EmailRegex        = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	AlphanumericRegex = regexp.MustCompile("^[a-zA-Z0-9_]*$")

	UploadFileSizeLimit int64 = 20 << 20
)

const (
	Empty = ""

	/* Discount offers */

	DiscountCode20 = "OFFER_20"
	DiscountCode30 = "OFFER_30"
	DiscountCode25 = "OFFER_25"

	/* Cabin Types */

	Economy        = "Economy"
	PremiumEconomy = "Premium Economy"
	Business       = "Business"
	First          = "First"

	/* FlightRecord Time Formats */

	FlightRecordTimeFormat = "2006-01-02"

	/* Phone Dial Codes */

	PhoneDialCode = "+91"

	/* http Request constants */

	HttpContentTypeHeaderContentType = "content-type"
	MimeTypeFormData                 = "multipart/form-data"
	MimeTypeJsonFormData             = "application/json"
	HttpAcceptHeader                 = "accept"

	FlightRecordKey = "flightRecord"

	DownloadFlightRecordPath = "download"

	/* Allowed file types */

	TypeFileContentText = "text/plain"

	/* Application folder names */

	TemplateFolder = "templates"
	UploadFolder   = "uploads"
	ProjectFolder  = "Flight-Ticket-Aggregator"

	/* Flight Record Entries */

	FirstName     = "First_name"
	LastName      = "Last_name"
	PNR           = "PNR"
	FareClass     = "Fare_class"
	TravelDate    = "Travel_date"
	Pax           = "Pax"
	TicketingDate = "Ticketing_date"
	Email         = "Email"
	MobilePhone   = "Mobile_phone"
	BookedCabin   = "Booked_cabin"
	DiscountCode  = "Discount_code"
	RecordError   = "Error"

	TypeRecordPassed = "passed"
	TypeRecordFailed = "failed"

	TestFlightRecord      = "airline1_2020-10-30_flightRecord.csv"
	TestEmptyFlightRecord = "airline1_2020-10-30_testEmptyFlightRecord.csv"

	/* Environments */

	LocalEnv = "local"

	/* Airlines */

	Airline1 = "airline1"
	Airline2 = "airline2"
	Airline3 = "airline3"

	/* Fixed Length */

	FlightRecordNameEntriesLength       = 3
	FlightRecordEntriesLength           = 10
	FlightRecordPNREntryLength          = 6
	FlightRecordEmailEntryMinimumLength = 3
	FlightRecordEmailEntryMaximumLength = 254

	/* Fare classes */

	FareClassA = 65
	FareClassB = 69
	FareClassF = 70
	FareClassK = 75
	FareClassL = 76
	FareClassR = 82

	/* Return messages */

	MissingArgument      = "missing_argument"
	InvalidArgument      = "invalid_argument"
	NotFound             = "not_found"
	InternalError        = "internal_error"
	InvalidRequest       = "invalid_request"
	InvalidContentType   = "invalid_contentType"
	InvalidFile          = "invalid_file"
	FileNotFound         = "file_not_found"
	InvalidFilename      = "invalid_filename"
	EmptyFile            = "empty_file"
	InvalidTypeAssetion  = "invalid_type_assetion"
	InvalidInput         = "invalid_input"
	EmptyInput           = "empty_input"
	InvalidAirline       = "invalid_airline"
	TryAgainLater        = "try_Again_later"
	InvalidRequestMethod = "invalid_request_method"

	InvalidMail          = "email_invalid"
	InvalidPhone         = "phone_invalid"
	InvalidCabin         = "cabin_invalid"
	InvalidPNR           = "PNR_invalid"
	InvalidFareClass     = "fare_class_invalid"
	InvalidTicketingDate = "ticketing_date_invalid"
	InvalidTravelDate    = "travel_date_invalid"
	InvalidUploadedDate  = "uploaded_date_invalid"
	InvalidBooking       = "booking_invalid"
	PNRRepeated          = "PNR_must_be_unique"
)
