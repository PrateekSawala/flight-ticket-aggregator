package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/PrateekSawala/flight-ticket-aggregator/domain/logging"
	"github.com/PrateekSawala/flight-ticket-aggregator/space/rpc/space"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var (
	service = flag.String("service", os.Getenv("FTA_SERVICE_NAME"), "Service name")
)

func main() {
	logging.InitializeLogging()
	log := logging.Log("Initializing space server")

	s3Config := &aws.Config{
		Credentials: credentials.NewStaticCredentials(os.Getenv("FTA_SPACE_ACCESS_KEY"), os.Getenv("FTA_SPACE_SECRET_KEY"), ""),
		Endpoint:    aws.String(os.Getenv("FTA_SPACE_Endpoint")),
		Region:      aws.String("FTA_SPACE_Region"),
	}

	newSession := session.New(s3Config)
	s3Client := s3.New(newSession)
	bucket := os.Getenv("FTA_SPACE_BUCKET")

	server := Server{client: s3Client, bucket: bucket}

	handler := space.NewSpaceServer(&server, nil)
	serverFQDNandPort := os.Getenv("FTA_FQDN") + ":" + os.Getenv("FTA_PORT")

	log.Debugf("%s server started: %v", *service, serverFQDNandPort)
	log.Warnf("Server exited: %s", http.ListenAndServe(serverFQDNandPort, handler))
}

type Server struct {
	client *s3.S3
	bucket string
}
