package mail

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"

	"flight-ticket-aggregator/domain"
	"flight-ticket-aggregator/domain/logging"
)

// SendProcessedFlightRecordsMail ...
func SendProcessedFlightRecordsMail(uploadedfile string, processedfiles []string) error {
	log := logging.Log("SendProcessedFlightRecordsMail")
	log.Tracef("Start SendProcessedFlightRecordsMail")
	defer log.Tracef("End SendProcessedFlightRecordsMail")

	log.Tracef("Input uploadedfile: %+v, processedfiles: %+v", uploadedfile, processedfiles)

	if len(processedfiles) == 0 || uploadedfile == domain.Empty {
		return domain.ErrInvalidInput
	}

	// Creating a new message
	m := NewMessage()

	//Setting headers
	m.SetHeader("From", "")
	m.SetHeader("To", "")
	m.SetHeader("Subject", "Processed Flight Records")

	//Set body content of the email template
	body := template.Must(template.New("Template-Body").Parse(`<p>Plese find processed flight records in the attachment from uploaded files:{{.File}}</p>`))

	bodyContent := new(bytes.Buffer)

	paramsBody := map[string]string{
		"File": uploadedfile,
	}

	// Set parameters which will be filled in template
	body.Execute(bodyContent, paramsBody)

	//Setting the body of the message
	m.SetBody("text/html", string(bodyContent.Bytes()))

	// Looping over the document
	for _, document := range processedfiles {
		uploadedFilePath := fmt.Sprintf("%s/%s", domain.UploadFolder, document)
		fileBuffer, err := ioutil.ReadFile(uploadedFilePath)
		if err != nil {
			log.Debugf("Error while reading the file %s", document)
			continue
		}

		attachment := bytes.NewReader(fileBuffer)
		// Attaching a file
		m.AttachReader(document, attachment)
	}

	// send mail
	err := sendMail(m)
	if err != nil {
		log.Debugf("Error while sending mail: %s", err)
		return err
	}
	return nil
}
