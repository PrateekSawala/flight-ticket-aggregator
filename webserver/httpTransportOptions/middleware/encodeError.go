package middleware

import (
	"context"
	"encoding/json"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/sirupsen/logrus"

	"flight-ticket-aggregator/domain"
)

func MakeEncodeErrorFunc(logger *logrus.Entry) httptransport.ErrorEncoder {
	return func(ctx context.Context, err error, w http.ResponseWriter) {
		if err == nil {
			logger.Errorf("No error found %s", err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(HTTPStatusCode(err))
		json.NewEncoder(w).Encode(err.Error())
	}
}

func HTTPStatusCode(err error) int {
	switch err {
	case domain.ErrInvalidRequest, domain.ErrInvalidRequestMethod, domain.ErrInvalidInput, domain.ErrInvalidFile, domain.ErrInvalidContentType:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
