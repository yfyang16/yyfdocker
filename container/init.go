package container

import (
    "log"
    "os"
    "os/exec"
    "syscall"
    "io/ioutil"
    "fmt"
    "strings"
)

/*
 * RunContainerInitProcess:
 *   This function will be executed in a yyfdocker inside.
 *   mount "/proc" to yydocker /proc
 */
func RunContainerInitProcess() error {
    log.Printf("** RunContainerInitProcess START **\n")

    // wait and read user cmds from the read pipe
    cmdArray := ReadUserCmd() 
    if cmdArray == nil {
        log.Panicf("[ReadUserCmd] Could not read pipe!")
        return fmt.Errorf("[ReadUserCmd] Could not read pipe!")
    }

    // func Mount(source string, target string, fstype string, flags uintptr, data string)
    // We should explicitly state that this mount namespace is independent.
    err := syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, "")
    if err != nil {
        log.Printf(err.Error())
    }

    // mount "/proc"
    // flag meaning:
    //   - MS_NOEXEC: no other program can be executed in this file system
    //   - MS_NOSUID: "set-user-ID" and "set-group-ID" are not allowed
    //   - MS_NODEV : default flag
    defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
    syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")

    // find the absolute path of a cmd in $PATH of OS
    // e.g. /bin/bash is the absolute path of bash cmd
    absolutePath, err := exec.LookPath(cmdArray[0])
    if err != nil {
        log.Panicf("[ReadUserCmd] Could not find the corresponding path of %s!", cmdArray[0])
        return err
    }

    log.Printf("** RunContainerInitProcess END **\n")

    // syscall.Exec will execute "command" and replace the init process with
    // "command" process. (So the first process (pid == 1) will be "command")
    err = syscall.Exec(absolutePath, cmdArray[0:], os.Environ())
    if err != nil {
        log.Printf(err.Error())
    }

    return err
}


/** Init process receive an extra file which is a read pipe recording the user cmds */
func ReadUserCmd() []string {
    readPipe := os.NewFile(uintptr(3), "readPipe")
    rawCommand, err := ioutil.ReadAll(readPipe)
    if err != nil {
        log.Panicf("[ReadUserCmd] Could not read pipe!")
        return nil
    }
    return strings.Split(string(rawCommand), " ")

}
