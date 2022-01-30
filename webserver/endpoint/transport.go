package endpoint

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	kitlogrus "github.com/go-kit/kit/log/logrus"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/sirupsen/logrus"

	"flight-ticket-aggregator/domain"
)

func MakeUploadFlightRecordsHandler(service WebService, logger *logrus.Entry) http.Handler {
	options := []httptransport.ServerOption{
		httptransport.ServerBefore(ExtractAcceptHeaderIntoContext),
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(kitlogrus.NewLogger(logger))),
		httptransport.ServerErrorEncoder(EncodeErrorFunc(logger)),
	}

	mw := endpoint.Chain(
		AcceptContentTypeFormDataValidationMiddleware(),
	)

	endpointHandler := httptransport.NewServer(
		mw(makeValidateUploadedFlightRecordEndpoint(service)),
		decodeUploadFileRequest,
		encodeUploadFileResponse,
		options...,
	)

	return endpointHandler
}

func MakeDownloadFlightRecordsHandler(service WebService, logger *logrus.Entry) http.Handler {
	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(kitlogrus.NewLogger(logger))),
		httptransport.ServerErrorEncoder(EncodeErrorFunc(logger)),
	}

	endpointHandler := httptransport.NewServer(
		makeValidateDownloadFlightRecordEndpoint(service),
		decodeDownloadFileRequest,
		encodeDownloadFileResponse,
		options...,
	)

	return endpointHandler
}

func decodeUploadFileRequest(_ context.Context, request *http.Request) (interface{}, error) {
	if err := IsMethodPost(request.Method); err != nil {
		return nil, err
	}
	if request.ContentLength <= 0 {
		return nil, domain.ErrInvalidFile
	}
	return request, nil
}

func decodeDownloadFileRequest(_ context.Context, request *http.Request) (interface{}, error) {
	if err := IsMethodGet(request.Method); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeUploadFileResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func encodeDownloadFileResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	recordFile, ok := response.([]byte)
	if !ok {
		return domain.ErrInvalidFile
	}
	if _, err := w.Write(recordFile); err != nil {
		return json.NewEncoder(w).Encode(response)
	}
	return nil
}

func EncodeErrorFunc(logger *logrus.Entry) httptransport.ErrorEncoder {
	return func(ctx context.Context, err error, w http.ResponseWriter) {
		if err == nil {
			logger.Errorf("No error found %s", err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(HTTPStatusCode(err))
		json.NewEncoder(w).Encode(map[string]interface{}{
			"errorMessage": err.Error(),
		})
	}
}

func HTTPStatusCode(err error) int {
	switch err {
	case domain.ErrInvalidRequest, domain.ErrInvalidRequestMethod, domain.ErrInvalidInput, domain.ErrInvalidFile:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

func ExtractAcceptHeaderIntoContext(ctx context.Context, r *http.Request) context.Context {
	if acceptHeaderValue := r.Header.Get(domain.HttpContentTypeHeaderContentType); acceptHeaderValue != "" {
		return context.WithValue(ctx, AcceptHeader, acceptHeaderValue)
	}
	return context.WithValue(ctx, AcceptHeader, domain.MimeTypeJsonFormData)
}
