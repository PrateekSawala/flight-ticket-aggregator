package endpoint

import (
	"context"
	"net/http"
	"path/filepath"

	"github.com/go-kit/kit/endpoint"

	"github.com/PrateekSawala/flight-ticket-aggregator/domain"
)

func makeValidateUploadedFlightRecordEndpoint(s WebService) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		request := req.(*http.Request)
		var writer http.ResponseWriter
		// Validate upload size
		request.Body = http.MaxBytesReader(writer, request.Body, domain.UploadFileSizeLimit)

		// Parse the multipart form in the request
		err := request.ParseMultipartForm(100000)
		if err != nil {
			s.Logger.Debugf("MultipartForm error: %s", err)
			return nil, domain.ErrInvalidRequest
		}

		// Get the reference to the parsed multipart form
		multipartForm := request.MultipartForm
		files := multipartForm.File[domain.FlightRecordKey]

		uploadFlightRecordsResp, err := s.UploadFlightRecords(files)
		if err != nil {
			return nil, err
		}
		return uploadFlightRecordsResp, nil
	}
}

func makeValidateDownloadFlightRecordEndpoint(s WebService) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		request := req.(*http.Request)

		path, documentName := filepath.Split(request.URL.Path)
		if path == domain.Empty || documentName == domain.Empty {
			return nil, domain.ErrInvalidInput
		}

		downloadFlightRecordsResp, err := s.DownloadFlightRecord(path, documentName)
		if err != nil {
			return nil, err
		}
		return downloadFlightRecordsResp, nil
	}
}
