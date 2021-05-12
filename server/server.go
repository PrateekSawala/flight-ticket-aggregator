package server

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"flight-ticket-aggregator/domain/logging"
	"flight-ticket-aggregator/endpoint"
)

var (
	serverPort = "localhost:3002"
)

func Run() {
	/* Initialize Logging */
	logging.InitializeLogging()

	log := logging.Log("Run")

	// Check if port is provided in environment configuration
	if os.Getenv("TEST_PORT") != "" {
		serverPort = os.Getenv("TEST_PORT")
	}

	handler := setUpServer()
	srv := &http.Server{Addr: serverPort, Handler: handler}
	go func() {
		log.Infof("Starting server at port: %s", serverPort)

		// initiate fileWatcher
		go fileWatcher()

		// Serv
		err := srv.ListenAndServe()
		if err != nil {
			if err == http.ErrServerClosed {
				log.Info("Server shut down.")
			} else {
				log.WithError(err).WithField("server_port", srv.Addr).Fatal("failed to start server")
			}
		}
	}()

	// Wait for an interrupt
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM) // interrupt signal sent from terminal, system
	<-sigint

	log.Info("Shutting down server")
}

func setUpServer() http.Handler {
	mux := http.NewServeMux()

	// Handlers that serve static and uploaded folder files
	staticFileServer := http.FileServer(http.Dir("./static"))
	uploadedFileServer := http.FileServer(http.Dir("./uploads"))

	mux.Handle("/", staticFileServer)
	mux.Handle("/uploads/", http.StripPrefix("/uploads/", uploadedFileServer))
	mux.Handle("/upload/flightRecord", AcceptContentTypeValidationMiddleware(http.HandlerFunc(endpoint.UploadFlightRecords)))
	return mux
}
