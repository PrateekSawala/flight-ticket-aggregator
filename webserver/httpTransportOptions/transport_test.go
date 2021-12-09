package httpTransportOptions

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeUploadFileRequestSuccess(t *testing.T) {
	body := fmt.Sprintf("\"%s\"", "DecodeUploadFileRequest")

	httpReq := httptest.NewRequest(http.MethodPost, "https://localhost:3002/upload/flightRecord", strings.NewReader(body))

	respose, err := DecodeUploadFileRequest(nil, httpReq)

	assert.NoError(t, err)
	assert.Equal(t, httpReq, respose)
}

func TestEncodeUploadFileResponseSuccess(t *testing.T) {
	bodyValue := "EncodeUploadFileResponse"
	expected := fmt.Sprintf("\"%s\"", bodyValue)

	r := httptest.NewRecorder()

	err := EncodeUploadFileResponse(nil, r, &bodyValue)
	assert.NoError(t, err)

	actual, err := ioutil.ReadAll(r.Body)
	assert.NoError(t, err)
	assert.JSONEq(t, expected, string(actual))
}

func TestDecodeDownloadFileRequestSuccess(t *testing.T) {
	body := fmt.Sprintf("\"%s\"", "DecodeDownloadFileRequest")

	httpReq := httptest.NewRequest(http.MethodGet, "https://localhost:3002/download/flightRecord", strings.NewReader(body))

	respose, err := DecodeDownloadFileRequest(nil, httpReq)

	assert.NoError(t, err)
	assert.Equal(t, httpReq, respose)
}

func TestEncodeDownloadFileResponseSuccess(t *testing.T) {
	bodyValue := "EncodeDownloadFileResponse"

	r := httptest.NewRecorder()
	err := EncodeDownloadFileResponse(nil, r, []byte(bodyValue))
	assert.NoError(t, err)
}
