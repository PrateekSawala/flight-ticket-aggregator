package main

import (
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"flight-ticket-aggregator/webserver/endpoint"

	"flight-ticket-aggregator/domain/logging"
	"flight-ticket-aggregator/space/rpc/space"
	"flight-ticket-aggregator/ticket/rpc/ticket"
)

var (
	service       = flag.String("service", os.Getenv("FTA_SERVICE_NAME"), "Service name")
	ticketService = ticket.NewTicketProtobufClient(os.Getenv("FTA_TICKET"), &http.Client{})
	spaceService  = space.NewSpaceProtobufClient(os.Getenv("FTA_SPACE"), &http.Client{})
)

func main() {
	logging.InitializeLogging()
	log := logging.Log("main")

	serverFQDNandPort := os.Getenv("FTA_FQDN") + ":" + os.Getenv("FTA_PORT")
	handler := setupServer()
	srv := &http.Server{Addr: serverFQDNandPort, Handler: handler}
	go func() {
		log.Debugf("%s server started: %s", *service, serverFQDNandPort)
		// initiate fileWatcher
		go fileWatcher()
		err := srv.ListenAndServe()
		if err != nil {
			if err == http.ErrServerClosed {
				log.Info("Server shut down.")
			} else {
				log.Error("failed to start server", err)
			}
		}
	}()
	// Wait for an interrupt
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM) // interrupt signal sent from terminal, system
	<-sigint
	log.Info("Shutting down server")
}

func setupServer() http.Handler {
	mux := http.NewServeMux()
	// Handlers that serve static and uploaded folder files
	staticFileServer := http.FileServer(http.Dir(os.Getenv("FTA_WEBSERVER_STATIC_DIR")))
	localUploadedFileServer := http.FileServer(http.Dir(os.Getenv("FTA_SPACE_LOCAL_DIR")))
	mux.Handle("/", staticFileServer)
	mux.Handle("/localfiles/", http.StripPrefix("/localfiles/", localUploadedFileServer))

	endpointLogger := logging.Log("endpoint").LogrusEntry
	endpointService := endpoint.NewServiceClient(endpoint.WebService{SpaceService: spaceService, TicketService: ticketService, Logger: endpointLogger})
	mux.Handle("/upload/flightRecord", endpoint.MakeUploadFlightRecordsHandler(*endpointService, endpointLogger))
	mux.Handle("/download/flightRecord/", http.StripPrefix("/download/", endpoint.MakeDownloadFlightRecordsHandler(*endpointService, endpointLogger)))

	return mux
}
