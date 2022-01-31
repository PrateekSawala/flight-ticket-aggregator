package system

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
	"unicode"

	"github.com/PrateekSawala/flight-ticket-aggregator/domain"
)

func IsAlphabetic(value string) bool {
	for _, r := range value {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

func IsFileCSV(fileread io.Reader) error {
	bufferReader := bufio.NewReader(fileread)
	initialByte, _ := bufferReader.Peek(512)
	contentType := http.DetectContentType(initialByte)
	if !strings.Contains(contentType, domain.TypeFileContentText) {
		return domain.ErrInvalidFile
	}
	return nil
}

func IsUploadedFlightRecordValid(flightRecord string) (*domain.FightRecordInfo, error) {
	response := &domain.FightRecordInfo{}
	splitResponse := strings.Split(flightRecord, "_")
	if len(splitResponse) != domain.FlightRecordNameEntriesLength {
		return response, domain.ErrInvalidFilename
	}

	flightRecordAirline := splitResponse[0]
	flightRecordUploadedDate := splitResponse[1]
	flightRecordName := splitResponse[2]

	if !domain.ValidAirlines[flightRecordAirline] {
		return response, domain.ErrInvalidFilename
	}

	_, err := time.Parse(domain.FlightRecordTimeFormat, flightRecordUploadedDate)
	if err != nil {
		return response, domain.ErrUploadedDate
	}
	response.Filepath = fmt.Sprintf("%s/%s/%s", domain.FlightRecordKey, flightRecordAirline, strings.Replace(flightRecordUploadedDate, "-", "/", -1))
	response.Filename = flightRecordName
	response.AirlineName = flightRecordAirline
	return response, nil
}
