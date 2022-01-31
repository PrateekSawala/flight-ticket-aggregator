package main

import (
	"context"
	"log"
	"net/http"

	"github.com/PrateekSawala/flight-ticket-aggregator/domain"
	"github.com/PrateekSawala/flight-ticket-aggregator/mail/rpc/mail"
)

var (
	client mail.Mail
)

func main() {
	client = mail.NewMailProtobufClient("http://localhost:3004", &http.Client{})
	sendProcessedFlightRecordsMail()
}

func sendProcessedFlightRecordsMail() {
	sendProcessedFlightRecordsMailInput := &mail.SendProcessedFlightRecordsMailInput{
		AirlineName:      domain.Airline1,
		Processedfiles:   []string{"flightRecord_passedRecord.csv", "flightRecord_failedRecord.csv"},
		UploadedFilePath: "/flightrecord/airline1/2020/10/12/",
		UploadedFileName: "flightRecord.csv",
	}
	_, err := client.SendProcessedFlightRecordsMail(context.Background(), sendProcessedFlightRecordsMailInput)
	if err != nil {
		log.Println("Error", err)
	}
}
