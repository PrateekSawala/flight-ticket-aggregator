package endpoint

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"flight-ticket-aggregator/domain"
	"github.com/stretchr/testify/assert"
	"mime/multipart"
)

func TestFlightRecordUploadSuccess(t *testing.T) {
	// Declare test file
	testFilePath := fmt.Sprintf("../%s/%s", domain.TemplateFolder, "flightRecord.csv")

	// Open file
	testfile, err := os.Open(testFilePath)
	if err != nil {
		t.Error(err)
	}
	defer testfile.Close()

	// Create a writer
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	// Write the domain
	part, err := writer.CreateFormFile("flightRecord", filepath.Base(testFilePath))
	if err != nil {
		writer.Close()
		t.Error(err)
	}

	// Upload the file
	io.Copy(part, testfile)
	writer.Close()

	request := httptest.NewRequest("POST", "http://localhost:3002/upload/flightRecord", body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	response := httptest.NewRecorder()

	UploadFlightRecords(response, request)
	assert.NoError(t, err)
}

func TestFlightRecordUploadFailure(t *testing.T) {
	// Declare test file
	testFilePath := fmt.Sprintf("../%s/%s", domain.TemplateFolder, "flightRecord.csv")

	// Open file
	testfile, err := os.Open(testFilePath)
	if err != nil {
		t.Error(err)
	}
	defer testfile.Close()

	// Create a writer
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	// Write the domain
	part, err := writer.CreateFormFile("flightRecord", filepath.Base(testFilePath))
	if err != nil {
		writer.Close()
		t.Error(err)
	}

	// Upload the file
	io.Copy(part, testfile)
	writer.Close()

	t.Run("Should return error of invalid request with response code 400", func(t *testing.T) {

		request := httptest.NewRequest("Get", "http://localhost:3002/upload/flightRecord", body)
		request.Header.Set("Content-Type", writer.FormDataContentType())
		response := httptest.NewRecorder()

		UploadFlightRecords(response, request)
		if response.Code != http.StatusBadRequest {
			t.Error("response code not 400")
		}

		assert.EqualError(t, errors.New(strings.TrimSpace(response.Body.String())), domain.ErrInvalidRequest.Error())
	})

	t.Run("Should return error of invalid request with response code 500", func(t *testing.T) {

		request := httptest.NewRequest("POST", "http://localhost:3002/upload/flightRecord", body)
		request.Header.Set("Content-Type", "multipart/formdata")
		response := httptest.NewRecorder()

		UploadFlightRecords(response, request)
		if response.Code != http.StatusInternalServerError {
			t.Error("response code not 500")
		}

		assert.EqualError(t, errors.New(strings.TrimSpace(response.Body.String())), domain.ErrInvalidRequest.Error())
	})
}
