package endpoint

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"flight-ticket-aggregator/domain"
	"flight-ticket-aggregator/domain/logging"
	"flight-ticket-aggregator/space/rpc/space"
	"flight-ticket-aggregator/ticket/rpc/ticket"

	"github.com/stretchr/testify/assert"
)

func TestMakeUploadFlightRecordsHandler(t *testing.T) {
	ticketService := ticket.NewTicketProtobufClient("http://localhost:3003", &http.Client{})
	spaceService := space.NewSpaceProtobufClient("http://localhost:3005", &http.Client{})

	endpointLogger := logging.Log("TestMakeUploadFlightRecordsHandler").LogrusEntry
	service := WebService{SpaceService: spaceService, TicketService: ticketService, Logger: endpointLogger}

	testHandler := MakeUploadFlightRecordsHandler(service, endpointLogger)

	// Declare test file
	testFilePath := fmt.Sprintf("../templates/airline1_2020-10-12_flightRecord.csv")
	// Open file
	testfile, err := os.Open(testFilePath)
	if err != nil {
		t.Error(err)
		return
	}
	defer testfile.Close()

	// Create a writer
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	// Write the domain
	part, err := writer.CreateFormFile(domain.FlightRecordKey, filepath.Base(testFilePath))
	if err != nil {
		writer.Close()
		t.Error(err)
		return
	}
	// Upload the file
	io.Copy(part, testfile)

	writer.Close()

	request := httptest.NewRequest(http.MethodPost, "/upload/flightRecord", body)
	request.Header.Add("Content-Type", writer.FormDataContentType())
	response := httptest.NewRecorder()

	testHandler.ServeHTTP(response, request)

	result := response.Result()
	assert.Equal(t, http.StatusOK, result.StatusCode)
	err = result.Body.Close()
	assert.NoError(t, err)
}

func TestMakeDownloadFlightRecordsHandler(t *testing.T) {
	spaceService := space.NewSpaceProtobufClient("http://localhost:3005", &http.Client{})

	endpointLogger := logging.Log("TestMakeDownloadFlightRecordsHandler").LogrusEntry
	service := WebService{SpaceService: spaceService, Logger: endpointLogger}

	testHandler := http.StripPrefix("/download/", MakeDownloadFlightRecordsHandler(service, endpointLogger))

	request := httptest.NewRequest(http.MethodGet, "/download/flightRecord/airline1/2020/10/12/flightRecord.csv", nil)
	response := httptest.NewRecorder()

	testHandler.ServeHTTP(response, request)

	result := response.Result()
	assert.Equal(t, http.StatusOK, result.StatusCode)
	err := result.Body.Close()

	// Check result StatusCode
	if result.StatusCode != http.StatusOK {
		return
	}

	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		t.Error(err)
		return
	}
	err = ioutil.WriteFile("../templates/output.csv", body, 0777)
	if err != nil {
		t.Error(err)
		return
	}
	assert.NoError(t, err)
}
