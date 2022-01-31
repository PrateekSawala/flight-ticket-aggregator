package endpoint

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"testing"

	"github.com/PrateekSawala/flight-ticket-aggregator/domain/logging"
	"github.com/PrateekSawala/flight-ticket-aggregator/space/rpc/space"
	"github.com/stretchr/testify/assert"

	"github.com/PrateekSawala/flight-ticket-aggregator/domain"
)

func TestUploadFlightRecordError(t *testing.T) {
	endpointLogger := logging.Log("TestUploadFlightRecordError").LogrusEntry
	service := WebService{Logger: endpointLogger}

	files := []*multipart.FileHeader{&multipart.FileHeader{Filename: "test.csv"}}

	response, err := service.UploadFlightRecords(files)
	expectedResponse := []*domain.FileStatus{&domain.FileStatus{Upload: false, Filename: "test.csv"}}

	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, response)
}

func TestDownloadFlightRecordError(t *testing.T) {
	spaceService := space.NewSpaceProtobufClient("http://localhost:3005", &http.Client{})

	endpointLogger := logging.Log("TestDownloadFlightRecordError").LogrusEntry
	service := WebService{SpaceService: spaceService, Logger: endpointLogger}

	_, err := service.DownloadFlightRecord("test", "test.csv")
	expectedErr := fmt.Errorf("twirp error internal: %s", domain.FileNotFound)
	assert.EqualError(t, expectedErr, err.Error())
}
