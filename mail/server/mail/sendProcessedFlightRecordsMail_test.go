package main

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/PrateekSawala/flight-ticket-aggregator/domain"
	"github.com/PrateekSawala/flight-ticket-aggregator/mail/rpc/mail"
	"github.com/stretchr/testify/assert"
)

var (
	testclient = mail.NewMailProtobufClient("http://localhost:3004", &http.Client{})
)

func TestSendProcessedFlightRecordsMailError(t *testing.T) {
	_, err := testclient.SendProcessedFlightRecordsMail(context.Background(), &mail.SendProcessedFlightRecordsMailInput{})
	expectedErr := fmt.Errorf("twirp error internal: %s", domain.ErrInvalidInput)
	assert.EqualError(t, expectedErr, err.Error())
}
