package endpoint

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"sync"

	"flight-ticket-aggregator/domain"
	"flight-ticket-aggregator/domain/logging"
	"flight-ticket-aggregator/service"
	"flight-ticket-aggregator/utility"
)

func UploadFlightRecords(writer http.ResponseWriter, request *http.Request) {
	log := logging.Log("UploadFlightRecords")
	log.Tracef("Start File Upload")
	defer log.Tracef("End File Upload")

	// Check request method
	if err := utility.IsMethodPost(request.Method); err != nil {
		http.Error(writer, domain.ErrInvalidRequest.Error(), http.StatusBadRequest)
		return
	}

	// Check request payload content length
	if request.ContentLength <= 0 {
		log.Debugf("Request found with no content. Not processing")
		http.Error(writer, domain.ErrInvalidRequest.Error(), http.StatusBadRequest)
		return
	}

	// Limit upload size
	request.Body = http.MaxBytesReader(writer, request.Body, domain.UploadFileSizeLimit)

	// Upload Folder Route
	uploadFolderRoute := fmt.Sprintf("http://%s/%s", request.Host, domain.UploadFolder)

	// Parse the multipart form in the request
	err := request.ParseMultipartForm(100000)
	if err != nil {
		log.Debugf("MultipartForm error: %s", err)
		http.Error(writer, domain.ErrInvalidRequest.Error(), http.StatusInternalServerError)
		return
	}

	// Get the reference to the parsed multipart form
	multipartForm := request.MultipartForm
	files := multipartForm.File["flightRecord"]

	// Response fileStatus
	fileResp := []*domain.FileStatus{}

	var wg sync.WaitGroup

	// Loop over all files
	for i, _ := range files {
		wg.Add(1)
		go func(file *multipart.FileHeader) {
			// Decrement the counter when the go routine completes
			defer wg.Done()

			// for each fileheader, get a handle to the actual file
			fileread, err := file.Open()
			defer fileread.Close()
			if err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}

			filename := file.Filename
			// Check if fileName is present or not
			if filename == "" {
				log.Debugf("filename is not found in the uploaded file")
				return
			}

			// Declare fileStatus
			fileStatus := &domain.FileStatus{Upload: false, Filename: filename}
			fileResp = append(fileResp, fileStatus)

			// Read the file content
			bufferReader := bufio.NewReader(fileread)

			// Plausibility check for detecting the file content-type
			err = utility.IsFileCSV(bufferReader)
			if err != nil {
				log.Debugf("Error while checking file %s content-type, error: %s", filename, err)
				return
			}

			// Read the file
			fileBuffer, err := ioutil.ReadAll(bufferReader)
			if err != nil {
				log.Debugf("ioutil.ReadAll error: %s", err)
				http.Error(writer, domain.ErrInternal.Error(), http.StatusInternalServerError)
				return
			}

			// Process file
			records, err := service.UploadFlightRecord(filename, fileBuffer)
			if err != nil {
				log.Debugf("service.UploadFlightRecord error: %s", err)
				return
			}

			// Find processed records path
			if records.PassedRecordFileName != domain.Empty {
				records.PassedRecordFilePathUrl = fmt.Sprintf("%s/%s", uploadFolderRoute, records.PassedRecordFileName)
			}
			if records.FailedRecordFileName != domain.Empty {
				records.FailedRecordFilePathUrl = fmt.Sprintf("%s/%s", uploadFolderRoute, records.FailedRecordFileName)
			}

			fileStatus.Records = records
			fileStatus.Upload = true
		}(files[i])
	}

	// Wait for all go routines to finish
	wg.Wait()

	// Convert object to json
	jsonResp, err := json.Marshal(fileResp)
	if err != nil {
		log.Debugf("json.Marshal error : %s", err)
		http.Error(writer, domain.ErrInternal.Error(), http.StatusInternalServerError)
		return
	}

	// Set response header
	writer.Header().Set("Content-Type", "application/json")

	// Returning json response
	_, err = writer.Write(jsonResp)
	if err != nil {
		log.Debugf("writer.Write error : %s", err)
		http.Error(writer, domain.ErrInternal.Error(), http.StatusInternalServerError)
	}
}
