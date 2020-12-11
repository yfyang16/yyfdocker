package main

import (
    "os"
    "log"
    "./container"
)

/* 
 * Run: yyfdocker run [args]
 *   para: 
 *     - tty: if connect the stdin, stdout, stderr between my process and the os
 *     - command: argument of "yyfdocker init"
 */
func Run(tty bool, command string) {
    log.Printf("** Run START **\n")

    parent := container.NewParentProcess(tty, command)
    if err := parent.Start(); err != nil {
        log.Error(err)
    }
    parent.Wait()

    log.Printf("** Run END **\n")
    os.Exit(-1)
}