package main

import (
	"syscall"
	"log"
	"./cgroups"
	"./container"
	"strconv"
	"fmt"
)

func StopContainer(name string) {
	ci, err := GetContainerInfoByName(name)
	if err != nil {
		log.Panicf("[StopContainer] GetContainerInfoByName throws error: %v", err)
	}

	pidInt, err := strconv.Atoi(ci.Pid)
	log.Printf("[StopContainer] prepare to kill process %d", pidInt)
	if err = syscall.Kill(pidInt, syscall.SIGTERM); err != nil {
		log.Panicf("[StopContainer] cannot kill the container: %v", err)
	}

	ci.Status = container.STOP; ci.Pid = " "
	WriteContainerInfo(ci)
}

func RemoveContainer(name string) {
	ci, err := GetContainerInfoByName(name)
	if err != nil {
		log.Panicf("[RemoveContainer] GetContainerInfoByName throws error: %v", err)
	}

	if ci.Status != container.STOP {
		fmt.Printf("[RemoveContainer] The container is running: %v", err)
	}

	DeleteContainerInfo(name)
	container.DeleteWorkSpace(ci.Volume, name)
	cgroupManager := cgroups.NewCgroupManager("yyfdocker-cgroup")
	cgroupManager.Destroy()

}
