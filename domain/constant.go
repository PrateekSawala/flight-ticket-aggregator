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

	/* FlightRecord Time Format */
	FlightRecordTimeFormat = "2006-01-02"

	PhoneDialCode = "+91"

	/* http Request constants */
	ContentType  = "content-type"
	MimeFormData = "multipart/form-data"

	FlightRecordFormKey = "flightRecord"

	/* Allowed file type */
	TypeFileContentText = "text/plain"

	/* Application folder names */
	WatcherFolder  = "import"
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
	Error         = "Error"

	SMTP_PORT     = 0
	SMTP_HOST     = ""
	SMTP_USER     = ""
	SMTP_PASSWORD = ""

	TypeRecordPassed = "passed"
	TypeRecordFailed = "failed"
)
