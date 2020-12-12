package container

import (
	"log"
	"os"
	"syscall"
)

/*
 * RunContainerInitProcess:
 *	 This function will be executed in a yyfdocker inside.
 *	 mount "/proc" to yydocker /proc
 *
 *   para:
 *     - command: the expected first process in yyfdocker. e.g. /bin/bash
 *     - args: argument of "command"
 */
func RunContainerInitProcess(command string, args []string) error {
	log.Printf("** RunContainerInitProcess START **\n")

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
	argv := []string{command}

	log.Printf("** RunContainerInitProcess END **\n")

	// syscall.Exec will execute "command" and replace the init process with
	// "command" process. (So the first process (pid == 1) will be "command")
	err = syscall.Exec(command, argv, os.Environ())
	if err != nil {
		log.Printf(err.Error())
	}

	log.Printf("** RunContainerInitProcess ERROR END **\n")
	return err
}
