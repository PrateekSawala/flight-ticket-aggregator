package domain

import (
	"bytes"
	"encoding/csv"
)

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

type FightRecordInfo struct {
	Filename    string
	Filepath    string
	AirlineName string
}

type FileStatus struct {
	Upload   bool
	Filename string
	Records  interface{}
}

type FlightRecordWriter struct {
	Filename      string
	Filepath      string
	File          *bytes.Buffer
	FileWriter    *csv.Writer
	IsFileCreated bool
	TypeRecord    string
}
