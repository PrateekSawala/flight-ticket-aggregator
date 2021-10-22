package server

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"flight-ticket-aggregator/domain"
	"flight-ticket-aggregator/domain/logging"
	"flight-ticket-aggregator/service"
	"flight-ticket-aggregator/utility"

	"github.com/fsnotify/fsnotify"
)

// fileWatcher...
func fileWatcher() {
	log := logging.Log("fileWatcher")

	watcherFolderPath := fmt.Sprintf("./%s", domain.WatcherFolder)
	// Read/watch files.
	log.Tracef("Start watcher on %s folder", watcherFolderPath)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Tracef("watcher-error %s", err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
					log.Tracef("created or modified %s", event.Name)
					go func() {
						err := ReadFlightRecord(event.Name)
						if err != nil {
							log.Errorf("Error while reading flight record %s, error: %s", event.Name, err)
						}
					}()
				}
			case err := <-watcher.Errors:
				log.Tracef("watcher error: %s", err)
			}
		}
	}()

	err = watcher.Add(watcherFolderPath)
	if err != nil {
		log.Tracef("watcher.Add error %s", err)
	}
	<-done
}

func ReadFlightRecord(filename string) error {
	log := logging.Log("ReadFlightRecord")

	// Open the file
	importfile := fmt.Sprintf("./%s", filename)
	file, err := os.Open(importfile)
	if err != nil {
		log.Debugf("os.Open Error %s", err)
		return err
	}
	defer file.Close()

	// Read the file content
	bufferReader := bufio.NewReader(file)

	// Plausibility check for detecting the file content-type
	if err := utility.IsFileCSV(bufferReader); err != nil {
		log.Debugf("IsFileCSV Error %s", err)
		return err
	}

	fileBuffer, err := ioutil.ReadAll(bufferReader)
	if err != nil {
		log.Debugf("ioutil.ReadAll error : %s", err)
		return err
	}

	// Find document name
	_, documentName := filepath.Split(filename)
	_, err = service.UploadFlightRecord(documentName, fileBuffer)
	if err != nil {
		log.Debugf("UploadFlightRecord error : %s", err)
		return err
	}
	return nil
}
