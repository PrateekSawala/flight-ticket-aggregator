package main

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"flight-ticket-aggregator/domain/logging"
	"flight-ticket-aggregator/domain/system"
	"flight-ticket-aggregator/ticket/rpc/ticket"

	"github.com/fsnotify/fsnotify"
)

func fileWatcher() {
	log := logging.Log("fileWatcher")

	watcherFolderPath := os.Getenv("FTA_WEBSERVER_WATCH_DIR")

	log.Debugf("Start watcher on %s folder", watcherFolderPath)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Errorf("watcher-error %s", err)
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
				log.Errorf("watcher error: %s", err)
			}
		}
	}()

	err = watcher.Add(watcherFolderPath)
	if err != nil {
		log.Errorf("watcher.Add error %s", err)
	}
	<-done
}

func ReadFlightRecord(filename string) error {
	log := logging.Log("ReadFlightRecord")

	importfile := fmt.Sprintf("./%s", filename)
	file, err := os.Open(importfile)
	if err != nil {
		log.Errorf("os.Open Error %s", err)
		return err
	}
	defer file.Close()

	bufferReader := bufio.NewReader(file)

	if err := system.IsFileCSV(bufferReader); err != nil {
		log.Errorf("IsFileCSV Error %s", err)
		return err
	}

	fileBuffer, err := ioutil.ReadAll(bufferReader)
	if err != nil {
		log.Errorf("ioutil.ReadAll error : %s", err)
		return err
	}

	_, documentName := filepath.Split(filename)
	_, err = ticketService.ProcessFlightRecord(context.Background(), &ticket.ProcessFlightRecordInput{Filename: documentName, FlightRecord: fileBuffer})
	if err != nil {
		log.Errorf("ticketService.ProcessFlightRecord error : %s", err)
		return err
	}
	return nil
}
