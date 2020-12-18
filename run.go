package main

import (
    "os"
    "log"
    "./container"
    "./cgroups"
    "./cgroups/subsystems"
)

/* 
 * Run: yyfdocker run [args]
 *   para: 
 *     - tty: if connect the stdin, stdout, stderr between my process and the os
 *     - command: argument of "yyfdocker init"
 */
func Run(tty bool, cmdArray []string, cfg *subsystems.ResourceConfig) {
    log.Printf("** Run START **\n")
    defer log.Printf("** Run END **\n")

    parent, writePipe := container.NewParentProcess(tty)
    if parent == nil {
        log.Panicf("[Run] Maybe anonymous pipe creation failure!")
        return
    }

    if err := parent.Start(); err != nil {
        log.Panicln(err)
    }

    rawCommand := strings.Join(cmdArray, " ")
    writePipe.WriteString(rawCommand)
    writePipe.Close()

    parent.Wait()
    os.Exit(0)
}