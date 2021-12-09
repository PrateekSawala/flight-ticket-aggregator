package main

import (
	"flag"
	"net/http"
	"os"

	"flight-ticket-aggregator/domain/logging"
	"flight-ticket-aggregator/mail/rpc/mail"
	"flight-ticket-aggregator/space/rpc/space"
	"flight-ticket-aggregator/ticket/rpc/ticket"
)

var (
	service           = flag.String("service", os.Getenv("FTA_SERVICE_NAME"), "Service Name")
	spaceService      = space.NewSpaceProtobufClient(os.Getenv("FTA_SPACE"), &http.Client{})
	mailService       = mail.NewMailProtobufClient(os.Getenv("FTA_MAIL"), &http.Client{})
	webServerHostName = flag.String("host", os.Getenv("FTA_WEBSERVER_HOSTNAME"), "Web Server Host Name")
)

func main() {
	logging.InitializeLogging()
	log := logging.Log("Initializing ticket server")

	server := Server{}
	handler := ticket.NewTicketServer(&server, nil)
	serverFQDNandPort := os.Getenv("FTA_FQDN") + ":" + os.Getenv("FTA_PORT")

	log.Debugf("%s server started: %v", *service, serverFQDNandPort)
	log.Warnf("Server exited: %s", http.ListenAndServe(serverFQDNandPort, handler))
}

// Server implements the twirp specs
type Server struct {
}
