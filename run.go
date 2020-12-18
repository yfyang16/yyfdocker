package main

import (
    "log"
    "./container"
    "./cgroups"
    "./cgroups/subsystems"
    "strings"
)

/* 
 * Run: yyfdocker run [args]
 *   para: 
 *     - tty: if connect the stdin, stdout, stderr between my process and the os
 *     - command: argument of "yyfdocker init"
 */
func Run(tty bool, cmdArray []string, cfg *subsystems.ResourceConfig) {
    log.Printf("** Run START (%v) **\n", cfg)
    defer log.Printf("** Run END **\n")

    parent, writePipe := container.NewParentProcess(tty)
    if parent == nil {
        log.Panicf("[Run] Maybe anonymous pipe creation failure!")
        return
    }

    if err := parent.Start(); err != nil {
        log.Panicln(err)
    }

    cgroupManager := cgroups.NewCgroupManager("yyfdocker-cgroup")
    defer cgroupManager.Destroy()
    cgroupManager.Set(cfg)
    cgroupManager.Apply(parent.Process.Pid)

    rawCommand := strings.Join(cmdArray, " ")
    writePipe.WriteString(rawCommand)
    writePipe.Close()

    parent.Wait()
}