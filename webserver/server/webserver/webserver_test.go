package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"flight-ticket-aggregator/domain/logging"

	"github.com/stretchr/testify/assert"
)

func TestWebPageSuccess(t *testing.T) {
	logging.InitializeLogging()

	testHandler := setupServer()

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "http://localhost:3002/", nil)
	testHandler.ServeHTTP(response, request)

	assert.Equal(t, http.StatusOK, response.Code)
}
