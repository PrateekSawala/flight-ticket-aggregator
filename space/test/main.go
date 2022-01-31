package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/PrateekSawala/flight-ticket-aggregator/space/rpc/space"
)

var (
	client = space.NewSpaceProtobufClient("http://localhost:3005", &http.Client{})
)

func main() {
	// savefile()
	// getfile()
	// deletefile()
}

func savefile() {
	file, err := ioutil.ReadFile("example.txt")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	_, err = client.SaveFile(context.Background(), &space.SaveFileInput{File: file, Filename: "test.txt", Filepath: "/test/"})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
}

func getfile() {
	result, err := client.GetFile(context.Background(), &space.GetFileInput{Filename: "test.txt", Filepath: "/test/"})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	ioutil.WriteFile("output.txt", result.File, 777)
}

func deletefile() {
	_, err := client.DeleteFile(context.Background(), &space.DeleteFileInput{Filename: "test.txt", Filepath: "/test/"})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
}
