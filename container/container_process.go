package container

import (
	"syscall"
	"os/exec"
	"os"
	"log"
)

/* 
 * NewParentProcess: fork an isolated process which will execute the "yyfdocker init"
 *   para: 
 *     - tty: if connect the stdin, stdout, stderr between my process and the os
 *     - command: argument of "yyfdocker init"
 */
func NewParentProcess(tty bool, command string) *exec.Cmd {
	log.Printf("** NewParentProcess START **\n")

	args := []string{"init", command}
	cmd := exec.Command("/proc/self/exe", args...)
	
    cmd.SysProcAttr = &syscall.SysProcAttr{
        Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS |
		syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
    }
	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	log.Printf("** NewParentProcess END **\n")
	return cmd
}