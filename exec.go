package main

import (
	"os/exec"
	"os"
	"log"
	"strings"
	"fmt"
	_ "./nsenter"
)

const ENV_EXEC_PID = "yyfdocker_pid"
const ENV_EXEC_CMD = "yyfdocker_cmd"

func ExecContainer(name string, cmdArray []string) {
	ci, err := GetContainerInfoByName(name)
	if err != nil {
		log.Panicf("[ExecContainer] GetContainerInfoByName throws error: %v", err)
	}
	pid := ci.Pid
	cmdStr := strings.Join(cmdArray, " ")

	cmd := exec.Command("/proc/self/exe", "exec")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("[ExecContainer] prepare to enter into a container. (pid: %v)\n", pid)

	os.Setenv(ENV_EXEC_PID, pid)
	os.Setenv(ENV_EXEC_CMD, cmdStr)

	fmt.Printf("[ExecContainer] cmdStr: %v\n", cmdStr)

	if err := cmd.Run(); err != nil {
		log.Panicf("[ExecContainer] Execute container %s error %v", name, err)
	}
}