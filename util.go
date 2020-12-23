package main

import (
	"./container"
	"fmt"
	"log"
	"io/ioutil"
	"path"
	"os"
	"strings"
)


func GetContainerInfoByContent(content []byte) *container.ContainerInfo {
	var _ci container.ContainerInfo
	contentList := strings.Split(string(content), ";")
	_ci.Id = contentList[0]; _ci.Pid = contentList[1]; _ci.Command = contentList[2]
	_ci.CreatedTime = contentList[3]; _ci.Status = contentList[4]; _ci.Name = contentList[5]; _ci.Volume = contentList[6]
	return &_ci
}


func GetContainerInfoByName(name string) (*container.ContainerInfo, error) {
	configFilePath := path.Join(fmt.Sprintf(container.DefaultInfoLocation, name), container.ConfigName)
	content, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		log.Panicf("[GetContainerPidByName] cannot read content in the config file: %v", err)
		return nil, fmt.Errorf("[GetContainerPidByName] cannot read content in the config file: %v", err)
	}
	return GetContainerInfoByContent(content), nil
}

func WriteContainerInfo(ci *container.ContainerInfo) error {
	recordDirPath := fmt.Sprintf(container.DefaultInfoLocation, ci.Name)
    if err := os.MkdirAll(recordDirPath, 0622); err != nil {
        log.Panicf("[WriteContainerInfo] Make all directories failed: %v", err)
        return fmt.Errorf("[WriteContainerInfo] Make all directories failed: %v", err)
    }

    recordFilePath := path.Join(recordDirPath, container.ConfigName)
    recordFile, err := os.OpenFile(recordFilePath, os.O_RDWR|os.O_CREATE, 0755); defer recordFile.Close()
    if err != nil {
        log.Panicf("[WriteContainerInfo] open record file failed: %v", err)
        return fmt.Errorf("[WriteContainerInfo] open record file failed: %v", err)
    }

    writeString := ci.Id+";"+ci.Pid+";"+ci.Command+";"+ci.CreatedTime+";"+ci.Status+";"+ci.Name+";"+ci.Volume
    if _, err := recordFile.WriteString(writeString); err != nil {
        log.Panicf("[WriteContainerInfo] write record failed: %v", err)
        return fmt.Errorf("[WriteContainerInfo] write record failed: %v", err)
    }
    return nil
}

func DeleteContainerInfo(name string) {
    log.Printf("** DeleteContainerInfo START **")
    defer log.Printf("** DeleteContainerInfo END **")

    recordDirPath := fmt.Sprintf(container.DefaultInfoLocation, name)
    if err := os.RemoveAll(recordDirPath); err != nil {
        log.Panicf("[DeleteContainerInfo] remove all files in a directory failed: %v", err)
    }
}