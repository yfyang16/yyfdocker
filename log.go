package main

import (
	"fmt"
	"./container"
	"os"
	"io/ioutil"
	"path"
	"log"
)

func LogContainer(containerName string) {
	logFilePath := path.Join(fmt.Sprintf(container.DefaultInfoLocation, containerName), container.ContainerLog)
	logFile, err := os.Open(logFilePath)
	defer logFile.Close()

	if err != nil {
		log.Panicf("[logContainer] open log file failed: %v", err); return
	}

	content, err := ioutil.ReadAll(logFile)
	if err != nil {
		log.Panicf("[logContainer] read file error: %v", err); return
	}

	fmt.Fprintf(os.Stdout, string(content))
}