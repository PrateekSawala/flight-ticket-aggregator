package endpoint

import (
	"bufio"
	"context"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"mime/multipart"
	"sync"

	"flight-ticket-aggregator/domain"
	"flight-ticket-aggregator/domain/system"
	"flight-ticket-aggregator/space/rpc/space"
	"flight-ticket-aggregator/ticket/rpc/ticket"
)

type Service interface {
	UploadFlightRecords(Files map[string][]*multipart.FileHeader) ([]*domain.FileStatus, error)
	DownloadFlightRecord(Path string, Filename string) ([]byte, error)
}

type WebService struct {
	SpaceService  space.Space
	TicketService ticket.Ticket
	Logger        *logrus.Entry
}

// NewClient returns a new service Client
func NewServiceClient(webService WebService) *WebService {
	return &webService
}

func (s *WebService) UploadFlightRecords(files []*multipart.FileHeader) ([]*domain.FileStatus, error) {
	fileResp := []*domain.FileStatus{}
	log := s.Logger

	var wg sync.WaitGroup

	for i, _ := range files {
		wg.Add(1)
		go func(file *multipart.FileHeader) {
			// Decrement the counter when the go routine completes
			defer wg.Done()

			fileStatus := &domain.FileStatus{Upload: false}
			fileResp = append(fileResp, fileStatus)

			filename := file.Filename
			if filename == "" {
				log.Debugf("filename not found in the uploaded file")
				return
			}
			fileStatus.Filename = filename

			fileread, err := file.Open()
			defer fileread.Close()
			if err != nil {
				log.Debugf("Error while checking file %s error: %s", filename, err)
				return
			}

			// Read the file content
			bufferReader := bufio.NewReader(fileread)

			err = system.IsFileCSV(bufferReader)
			if err != nil {
				log.Debugf("Error while checking file %s content-type, error: %s", filename, err)
				return
			}

			fileBuffer, err := ioutil.ReadAll(bufferReader)
			if err != nil {
				log.Debugf("ioutil.ReadAll error: %s", err)
				return
			}

			records, err := s.TicketService.ProcessFlightRecord(context.Background(), &ticket.ProcessFlightRecordInput{Filename: filename, FlightRecord: fileBuffer})
			if err != nil {
				log.Debugf("ticketService.ProcessFlightRecord error: %s", err)
				return
			}

			// Find processed records path
			processedRecords := []string{}
			if records.PassedRecordFileName != domain.Empty {
				processedRecords = append(processedRecords, records.PassedRecordFileName)
			}
			if records.FailedRecordFileName != domain.Empty {
				processedRecords = append(processedRecords, records.FailedRecordFileName)
			}
			fileStatus.Records = records
			fileStatus.Upload = true
		}(files[i])
	}
	// Wait for all go routines to finish
	wg.Wait()
	return fileResp, nil
}

func (s *WebService) DownloadFlightRecord(path string, filename string) ([]byte, error) {
	log := s.Logger
	getFileResp, err := s.SpaceService.GetFile(context.Background(), &space.GetFileInput{Filepath: path, Filename: filename})
	if err != nil {
		log.Tracef("spaceService.GetFile: %s", err)
		return nil, err
	}
	recordFile := getFileResp.File
	log.Tracef("Received Record file %s with a []byte count of: %d", filename, len(recordFile))
	return recordFile, nil
}
