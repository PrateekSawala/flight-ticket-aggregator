package httpTransportOptions

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/PrateekSawala/flight-ticket-aggregator/domain"
)

func DecodeUploadFileRequest(_ context.Context, request *http.Request) (interface{}, error) {
	if request.Method != http.MethodPost {
		return nil, domain.ErrInvalidRequestMethod
	}
	if request.ContentLength <= 0 {
		return nil, domain.ErrInvalidFile
	}
	return request, nil
}

func EncodeUploadFileResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func DecodeDownloadFileRequest(_ context.Context, request *http.Request) (interface{}, error) {
	if request.Method != http.MethodGet {
		return nil, domain.ErrInvalidRequestMethod
	}
	return request, nil
}

func EncodeDownloadFileResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	recordFile, ok := response.([]byte)
	if !ok {
		return domain.ErrInvalidFile
	}
	if _, err := w.Write(recordFile); err != nil {
		return json.NewEncoder(w).Encode(response)
	}
	return nil
}
