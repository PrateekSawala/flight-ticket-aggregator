package domain

type FlightRecord struct {
	FirstName     string
	LastName      string
	PNR           string
	FareClass     string
	TravelDate    string
	Pax           string
	TicketingDate string
	Email         string
	MobilePhone   string
	BookedCabin   string
}

type FileStatus struct {
	Upload   bool
	Filename string
	Records  *Record
}

type Record struct {
	PassedRecordFileName    string
	FailedRecordFileName    string
	PassedRecordFilePathUrl string
	FailedRecordFilePathUrl string
}
