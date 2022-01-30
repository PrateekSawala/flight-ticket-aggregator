package main

import (
	"bytes"
	"context"
	"fmt"
	"html/template"

	"flight-ticket-aggregator/domain"
	"flight-ticket-aggregator/domain/logging"
	"flight-ticket-aggregator/mail/rpc/mail"
	"flight-ticket-aggregator/space/rpc/space"
)

func (s *Server) SendProcessedFlightRecordsMail(ctx context.Context, input *mail.SendProcessedFlightRecordsMailInput) (*mail.SendProcessedFlightRecordsMailResponse, error) {
	log := logging.Log("SendProcessedFlightRecordsMail")
	log.Tracef("Start SendProcessedFlightRecordsMail")
	defer log.Tracef("End SendProcessedFlightRecordsMail")

	log.Tracef("Input: %+v", input)

	response := &mail.SendProcessedFlightRecordsMailResponse{}

	if input.AirlineName == domain.Empty || len(input.Processedfiles) == 0 || input.UploadedFileName == domain.Empty || input.UploadedFilePath == domain.Empty {
		return response, domain.ErrInvalidInput
	}

	airlineSupportMail, ok := domain.AirlinesMails[input.AirlineName]
	if !ok {
		return response, domain.ErrInvalidAirline
	}

	m := s.NewMessage()

	m.SetHeader("From", *smtpAccountEmail)
	m.SetHeader("To", airlineSupportMail)
	m.SetHeader("Subject", "Processed Flight Records")

	// Set body content of the email template
	body := template.Must(template.New("Template-Body").Parse(`<p>Please find the processed flight records in the attachment from the uploaded file:{{.File}}</p>`))
	bodyContent := new(bytes.Buffer)

	// Prepare paramsBody
	paramsBody := map[string]string{
		"File": input.UploadedFileName,
	}

	// Set parameters which will be filled in template
	body.Execute(bodyContent, paramsBody)

	m.SetBody("text/html", string(bodyContent.Bytes()))

	for _, document := range input.Processedfiles {
		s3FilePath := fmt.Sprintf("%s", input.UploadedFilePath)
		getFileResponse, err := spaceService.GetFile(context.Background(), &space.GetFileInput{Filepath: s3FilePath, Filename: document})
		if err != nil {
			log.Errorf("spaceService.GetFile error while fetching file from path %s/%s", s3FilePath, document)
			continue
		}
		attachment := bytes.NewReader(getFileResponse.File)
		m.AttachReader(document, attachment)
	}

	err := s.SendMail(m)
	if err != nil {
		log.Debugf("Error while sending mail: %s", err)
		return response, err
	}
	return response, nil
}
