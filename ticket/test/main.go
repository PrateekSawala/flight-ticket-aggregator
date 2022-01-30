package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"flight-ticket-aggregator/ticket/rpc/ticket"
)

var (
	client = ticket.NewTicketProtobufClient("http://localhost:3003", &http.Client{})
)

func main() {
	processFlightRecord()
}

func processFlightRecord() {
	filename := "airline1_2020-10-30_flightRecord.csv"

	importfile := fmt.Sprintf("../../ticket/templates/templates/airline1_2020-10-30_flightRecord.csv")
	fileBuffer, err := ioutil.ReadFile(importfile)
	if err != nil {
		log.Println("Error", err)
		return
	}

	_, err = client.ProcessFlightRecord(context.Background(), &ticket.ProcessFlightRecordInput{Filename: filename, FlightRecord: fileBuffer})
	if err != nil {
		log.Println("Error", err)
	}
}
