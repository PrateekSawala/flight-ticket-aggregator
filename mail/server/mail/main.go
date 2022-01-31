package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/PrateekSawala/flight-ticket-aggregator/domain/logging"
	"github.com/PrateekSawala/flight-ticket-aggregator/mail/rpc/mail"
	"github.com/PrateekSawala/flight-ticket-aggregator/space/rpc/space"
)

var (
	env              = flag.String("environment", os.Getenv("FTA_ENVIRONMENT"), "Development environment")
	service          = flag.String("service", os.Getenv("FTA_SERVICE_NAME"), "Service name")
	smtpAccountEmail = flag.String("SMTP account email", os.Getenv("FTA_SMTP_USER"), "The email account of mail sender")
	spaceService     = space.NewSpaceProtobufClient(os.Getenv("FTA_SPACE"), &http.Client{})
)

func main() {
	logging.InitializeLogging()
	log := logging.Log("Initializing mail server")

	goMailDialer := InitMailClient()

	server := Server{dialer: *goMailDialer}
	handler := mail.NewMailServer(&server, nil)

	serverFQDNandPort := os.Getenv("FTA_FQDN") + ":" + os.Getenv("FTA_PORT")

	log.Debugf("%s server started: %v", *service, serverFQDNandPort)
	// start the server.
	log.Warnf("Server exited: %s", http.ListenAndServe(serverFQDNandPort, handler))
}
