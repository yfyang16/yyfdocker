package main

import (
	"fmt"
	"./container"
	"io/ioutil"
	"os"
	"text/tabwriter"
	"log"
	"path"
)

func ListContainers() {
	containerIndexDir := fmt.Sprintf(container.DefaultInfoLocation, "")
	containerIndexDir = containerIndexDir[:len(containerIndexDir) - 1]
	containerDirs, err := ioutil.ReadDir(containerIndexDir)
	if err != nil {
		log.Printf("[ListContainers] read directories failed: %v", err);
		containerDirs = nil
	}

	var containers []*container.ContainerInfo
	for _, containerDir := range containerDirs {
		_container, err := GetContainerInfo(containerDir)
		if err != nil {
			log.Printf("[ListContainers] list a container error: %v", err); continue
		}
		containers = append(containers, _container)
	}

	w := tabwriter.NewWriter(os.Stdout, 15, 0, 1, ' ', 0)
	fmt.Fprint(w, "ID\tNAME\tPID\tSTATUS\tCOMMAND\tCREATED\t\n")
	for _, c := range containers {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t\n", c.Id, c.Name, c.Pid, c.Status, c.Command, c.CreatedTime)
	}
	if err := w.Flush(); err != nil {
		log.Panicf("[ListContainers] Flush error %v", err)
		return
	}
}

func GetContainerInfo(containerDir os.FileInfo) (*container.ContainerInfo, error) {
	configFilePath := path.Join(fmt.Sprintf(container.DefaultInfoLocation, containerDir.Name()), container.ConfigName)
	content, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		log.Printf("[GetContainerInfo] read congif file error: %v", err)
		return nil, fmt.Errorf("[GetContainerInfo] read congif file error: %v", err)
	}

	return GetContainerInfoByContent(content), nil
}