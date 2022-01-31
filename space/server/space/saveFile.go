package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/PrateekSawala/flight-ticket-aggregator/domain"
	"github.com/PrateekSawala/flight-ticket-aggregator/domain/logging"
	"github.com/PrateekSawala/flight-ticket-aggregator/space/rpc/space"

	"github.com/aws/aws-sdk-go/service/s3"
)

func (s *Server) SaveFile(ctx context.Context, input *space.SaveFileInput) (*space.SaveFileResponse, error) {
	if input == nil {
		return nil, domain.ErrInvalidInput
	}

	log := logging.Log("SaveFile")
	log.Tracef("Start")
	defer log.Tracef("End")

	response := &space.SaveFileResponse{}

	if strings.Count(input.Filename, ".") == 0 {
		return response, errors.New("Filename needs to contain '.' character")
	}

	if input.Filepath != "" {
		if !strings.HasPrefix(input.Filepath, "/") {
			input.Filepath = "/" + input.Filepath
		}
		if !strings.HasSuffix(input.Filepath, "/") {
			input.Filepath = input.Filepath + "/"
		}
	}

	if os.Getenv("FTA_ENVIRONMENT") == domain.LocalEnv {
		localFilePath := strings.TrimSpace(strings.Replace(input.Filepath, "/", " ", -1))
		filepath := fmt.Sprintf("/files/local/%s*%s", localFilePath, input.Filename)
		// Checking if the file already exist in the local filepath
		if _, err := os.Stat(filepath); err == nil {
			log.Debugf("file %s already exist in the local file system", input.Filename)
		}
		err := ioutil.WriteFile(filepath, input.File, os.ModePerm)
		if err != nil {
			log.Errorf("ioutil.WriteFile error while storing the file %s in the local file system: %s", filepath, err)
		}
		return response, err
	}

	file := bytes.NewReader(input.File)

	fullpath := input.Filepath + input.Filename

	object := s3.PutObjectInput{
		Body:   file,
		Bucket: &s.bucket,
		Key:    &fullpath,
	}
	_, err := s.client.PutObject(&object)
	if err != nil {
		log.Infof("Error while saving file %s in cloud storage: %s", fullpath, err)
		return response, err
	}

	return response, err
}
