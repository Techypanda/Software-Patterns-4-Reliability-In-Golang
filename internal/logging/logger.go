package logging

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

func ReadLogContents() (*string, error) {
	logFile, _ := filepath.Abs("../../data/log")
	contents, err := ioutil.ReadFile(logFile)
	if err != nil {
		return nil, err
	}
	parsed := string(contents)
	return &parsed, nil
}

func writeLogLogic(contents string) {
	logFile, _ := filepath.Abs("../../data/log")
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("Error, failed to create/append logfile, will retry in a bit", err.Error())
		time.Sleep(time.Millisecond * 100)
		go writeLogLogic(contents)
		return
	}
	_, err = file.WriteString(contents)
	if err != nil {
		log.Println("Error, failed to write to logfile, will retry in a bit", err.Error())
		time.Sleep(time.Millisecond * 100)
		go writeLogLogic(contents)
		return
	}
	log.Println("Successfully wrote to logfile")
}

func WriteLog(contents string) {
	log.Printf("A Request to log %s to logfile has been issued, will silently try to keep logging this to a file async\n", contents)
	go writeLogLogic(contents)
}

func WriteLogLn(contents string) {
	log.Printf("A Request to log %s to logfile has been issued, will silently try to keep logging this to a file async\n", contents)
	go writeLogLogic(contents + "\n")
}
