package main

import (
	"log"
	"os/exec"
	"path"
	"./container"
)

func CommitContainer(containerName string, imageName string) {
	containerMntPath := path.Join(container.MntPath, containerName)
	imageTarPath := path.Join(container.RootPath, imageName + ".tar")

	// exec "tar -czf /root/{imageName}.tar -C /root/mnt ."
	_, err := exec.Command("tar", "-czf", imageTarPath, "-C", containerMntPath, ".").CombinedOutput()
	if err != nil {
		log.Panicf("[CommitContainer] Create image Error: %v", err)
	}
}