package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/PrateekSawala/flight-ticket-aggregator/domain"
	"github.com/PrateekSawala/flight-ticket-aggregator/domain/logging"
	"github.com/PrateekSawala/flight-ticket-aggregator/space/rpc/space"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func (s *Server) GetFile(ctx context.Context, input *space.GetFileInput) (*space.GetFileResponse, error) {
	if input == nil {
		return nil, domain.ErrInvalidInput
	}

	log := logging.Log("GetFile")
	log.Tracef("Start")
	defer log.Tracef("End")

	if strings.Count(input.Filename, ".") == 0 {
		return nil, errors.New("GetFilename needs to contain '.' character")
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
		// Reading file from the local filepath, if file exist in the local filepath
		if _, err := os.Stat(filepath); !os.IsNotExist(err) {
			file, err := ioutil.ReadFile(filepath)
			if err != nil {
				log.Tracef("ioutil.ReadFile error while reading the file %s from the local file system: %s", filepath, err)
				return &space.GetFileResponse{}, err
			}
			return &space.GetFileResponse{File: file}, nil
		}
		log.Tracef("File %s not available locally", filepath)
		return nil, domain.ErrFileNotFound
	}
	fullpath := input.Filepath + input.Filename

	result, err := s.client.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: &s.bucket,
		Key:    &fullpath,
	})
	if err != nil {
		log.Infof("Error while fetching file %s from cloud storage: %s", fullpath, err)
		return &space.GetFileResponse{}, nil
	}
	defer result.Body.Close()
	log.Tracef("Object Size: %d", aws.Int64Value(result.ContentLength))

	fileContent, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return &space.GetFileResponse{}, err
	}
	return &space.GetFileResponse{File: fileContent}, nil
}
