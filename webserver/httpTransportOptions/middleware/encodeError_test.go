package middleware

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PrateekSawala/flight-ticket-aggregator/domain"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestEncodeError(t *testing.T) {
	testCases := []struct {
		name               string
		err                error
		expectedOutput     string
		expectedStatusCode int
	}{
		{
			name:               "invalid content type",
			err:                domain.ErrInvalidContentType,
			expectedOutput:     fmt.Sprintf("\"%s\"", domain.InvalidContentType),
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "invalid request",
			err:                domain.ErrInvalidRequest,
			expectedOutput:     fmt.Sprintf("\"%s\"", domain.InvalidRequest),
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, testcase := range testCases {
		t.Run(testcase.name, func(t *testing.T) {
			r := httptest.NewRecorder()
			ctx := context.WithValue(context.Background(), domain.HttpAcceptHeader, "application/json")
			MakeEncodeErrorFunc(logrus.NewEntry(logrus.StandardLogger()))(ctx, testcase.err, r)

			resp := r.Result()
			defer resp.Body.Close()

			actual, _ := ioutil.ReadAll(resp.Body)
			assert.JSONEq(t, testcase.expectedOutput, string(actual))
			assert.Equal(t, testcase.expectedStatusCode, r.Code)
		})
	}
}
