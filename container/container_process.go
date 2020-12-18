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
 */
func NewParentProcess(tty bool) (*exec.Cmd, *os.File) {
    log.Printf("** NewParentProcess START **\n")
    defer log.Printf("** NewParentProcess END **\n")

    readPipe, writePipe, err := os.Pipe()
    if err != nil {
        log.Fatalf("[NewParentProcess] create anonymous pipe failure!")
        return nil, nil
    }

    cmd := exec.Command("/proc/self/exe", "init")

    cmd.SysProcAttr = &syscall.SysProcAttr{
        Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS |
        syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
    }
    if tty {
        cmd.Stdin = os.Stdin
        cmd.Stdout = os.Stdout
        cmd.Stderr = os.Stderr
	}

    cmd.ExtraFiles = []*os.File{readPipe}    // extra file descriptor besides in, out and err
    return cmd, writePipe
}