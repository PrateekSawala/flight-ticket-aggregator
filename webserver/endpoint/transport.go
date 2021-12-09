package endpoint

import (
	"net/http"

	"github.com/go-kit/kit/endpoint"
	kitlogrus "github.com/go-kit/kit/log/logrus"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/sirupsen/logrus"

	"flight-ticket-aggregator/webserver/httpTransportOptions"
	"flight-ticket-aggregator/webserver/httpTransportOptions/middleware"
	"flight-ticket-aggregator/webserver/httpTransportOptions/serveroption"
)

func MakeUploadFlightRecordsHandler(service WebService, logger *logrus.Entry) http.Handler {
	options := []httptransport.ServerOption{
		httptransport.ServerBefore(serveroption.ExtractAcceptHeaderIntoContext),
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(kitlogrus.NewLogger(logger))),
		httptransport.ServerErrorEncoder(middleware.MakeEncodeErrorFunc(logger)),
	}

	mw := endpoint.Chain(
		middleware.MakeAcceptContentTypeFormDataValidationMiddleware(),
	)

	endpointHandler := httptransport.NewServer(
		mw(makeValidateUploadedFlightRecordEndpoint(service)),
		httpTransportOptions.DecodeUploadFileRequest,
		httpTransportOptions.EncodeUploadFileResponse,
		options...,
	)

	return endpointHandler
}

func MakeDownloadFlightRecordsHandler(service WebService, logger *logrus.Entry) http.Handler {
	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(kitlogrus.NewLogger(logger))),
		httptransport.ServerErrorEncoder(middleware.MakeEncodeErrorFunc(logger)),
	}

	endpointHandler := httptransport.NewServer(
		makeValidateDownloadFlightRecordEndpoint(service),
		httpTransportOptions.DecodeDownloadFileRequest,
		httpTransportOptions.EncodeDownloadFileResponse,
		options...,
	)

	return endpointHandler
}
