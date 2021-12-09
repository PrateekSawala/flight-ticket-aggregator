package main

import (
	"context"
	"io/ioutil"
	"net/http"
	"testing"

	"flight-ticket-aggregator/space/rpc/space"
	"github.com/stretchr/testify/assert"
)

var (
	testclient = space.NewSpaceProtobufClient("http://localhost:3005", &http.Client{})
)

func TestSaveFileSuccess(t *testing.T) {
	// Read file
	file, err := ioutil.ReadFile("../../templates/example.txt")
	if err != nil {
		t.Errorf("Error while reading file, Error: %s", err)
		return
	}
	// Save file
	_, err = testclient.SaveFile(context.Background(), &space.SaveFileInput{File: file, Filename: "test.txt", Filepath: "/test/"})
	if err != nil {
		t.Errorf("Error while saving file, Error: %s", err)
		return
	}
	assert.NoError(t, err)
}

func TestGetFileSuccess(t *testing.T) {
	result, err := testclient.GetFile(context.Background(), &space.GetFileInput{Filename: "test.txt", Filepath: "/test/"})
	if err != nil {
		t.Errorf("Error while fetching file, Error: %s", err)
		return
	}
	assert.NoError(t, err)
	ioutil.WriteFile("../../templates/output.txt", result.File, 777)
}

func TestDeleteFileSuccess(t *testing.T) {
	_, err := testclient.DeleteFile(context.Background(), &space.DeleteFileInput{Filename: "test.txt", Filepath: "/test/"})
	if err != nil {
		t.Errorf("Error while deleting file, Error: %s", err)
		return
	}
	assert.NoError(t, err)
}
