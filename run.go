package main

import (
    "log"
    "./container"
    "./cgroups"
    "./cgroups/subsystems"
    "strings"
    "os"
)

/* 
 * Run: yyfdocker run [cmd]
 */
func Run(tty bool, cmdArray []string, cfg *subsystems.ResourceConfig, volume string) {
    log.Printf("** Run START (cfg: %v); (cmdArray: %v) **\n", cfg, cmdArray)
    defer log.Printf("** Run END **\n")

    parent, writePipe := container.NewParentProcess(tty, volume)
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

    container.DeleteWorkSpace("/root", "/root/mnt", volume)
    os.Exit(0)
}