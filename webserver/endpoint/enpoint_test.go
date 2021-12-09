package endpoint

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"flight-ticket-aggregator/domain"
	"flight-ticket-aggregator/domain/logging"

	"github.com/stretchr/testify/assert"
)

func TestMakeValidateUploadedFlightRecordEndpointError(t *testing.T) {
	endpointLogger := logging.Log("TestMakeValidateUploadedFlightRecordEndpointError").LogrusEntry
	service := WebService{Logger: endpointLogger}

	endpoint := makeValidateUploadedFlightRecordEndpoint(service)
	request := httptest.NewRequest(http.MethodPost, "/upload/flightRecord", nil)
	_, err := endpoint(context.Background(), request)

	expectedErr := domain.ErrInvalidRequest
	assert.EqualError(t, expectedErr, err.Error())
}

func TestMakeValidateDownloadFlightRecordEndpointError(t *testing.T) {
	endpointLogger := logging.Log("TestMakeValidateDownloadFlightRecordEndpointError").LogrusEntry
	service := WebService{Logger: endpointLogger}

	endpoint := makeValidateDownloadFlightRecordEndpoint(service)
	request := httptest.NewRequest(http.MethodGet, "/download/flightRecord/", nil)
	_, err := endpoint(context.Background(), request)

	expectedErr := domain.ErrInvalidInput
	assert.EqualError(t, expectedErr, err.Error())
}
