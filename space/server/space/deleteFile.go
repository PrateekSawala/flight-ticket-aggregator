package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/PrateekSawala/flight-ticket-aggregator/domain"
	"github.com/PrateekSawala/flight-ticket-aggregator/domain/logging"
	"github.com/PrateekSawala/flight-ticket-aggregator/space/rpc/space"

	"github.com/aws/aws-sdk-go/service/s3"
)

func (s *Server) DeleteFile(ctx context.Context, input *space.DeleteFileInput) (*space.DeleteFileResponse, error) {
	if input == nil {
		return nil, domain.ErrInvalidInput
	}

	log := logging.Log("DeleteFile")
	log.Tracef("Start")
	defer log.Tracef("End")

	response := &space.DeleteFileResponse{}

	if input.Filename != "" {
		if strings.Count(input.Filename, ".") == 0 {
			return response, errors.New("Filename needs to contain '.' character")
		}
	}

	if input.Filepath != "" {
		if !strings.HasPrefix(input.Filepath, "/") {
			input.Filepath = "/" + input.Filepath
		}
	}

	if os.Getenv("FTA_ENVIRONMENT") == domain.LocalEnv {
		localFilePath := strings.TrimSpace(strings.Replace(input.Filepath, "/", " ", -1))
		filepath := fmt.Sprintf("/files/local/%s*%s", localFilePath, input.Filename)
		// Removing file from the local filepath if exist in the local filepath
		if _, err := os.Stat(filepath); !os.IsNotExist(err) {
			err = os.Remove(filepath)
			if err != nil {
				log.Tracef("Error while removing the file %s from the local file system: %s", filepath, err)
			}
			return response, err
		}
		return response, fmt.Errorf("File not available locally")
	}
	fullpath := input.Filepath

	if input.Filename != "" {
		if !strings.HasSuffix(fullpath, "/") {
			fullpath = fullpath + "/"
		}
		fullpath = fullpath + input.Filename
	}

	object := s3.DeleteObjectInput{
		Bucket: &s.bucket,
		Key:    &fullpath,
	}
	_, err := s.client.DeleteObject(&object)
	if err != nil {
		log.Infof("Error while deleting file %s in cloud storage: %s", fullpath, err)
		return response, err
	}
	return response, err
}
