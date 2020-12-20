package main

import (
	"log"
	"os/exec"
)

func CommitContainer(imageName string) {
	mntPath := "/root/mnt"
	imageTarPath := "/root/" + imageName + ".tar"

	// exec "tar -czf /root/{imageName}.tar -C /root/mnt ."
	_, err := exec.Command("tar", "-czf", imageTarPath, "-C", mntPath, ".").CombinedOutput()
	if err != nil {
		log.Panicf("[CommitContainer] Create image Error: %v", err)
	}
}