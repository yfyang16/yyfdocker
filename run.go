package main

import (
    "log"
    "./container"
    "./cgroups"
    "./cgroups/subsystems"
    "strings"
    "strconv"
    "time"
    "math/rand"
    "fmt"
)

/* 
 * Run: yyfdocker run [cmd]
 */
func Run(tty bool, cmdArray []string, cfg *subsystems.ResourceConfig, volume string, containerName string, imageName string) {
    log.Printf("** Run START (cfg: %v); (cmdArray: %v); (tty: %v) **\n", cfg, cmdArray, tty)
    defer log.Printf("** Run END **\n")

    containerId := randString(10)
    if containerName == "" {containerName = containerId}

    parent, writePipe := container.NewParentProcess(tty, volume, containerName, imageName)
    if parent == nil {
        log.Panicf("[Run] Maybe anonymous pipe creation failure!")
        return
    }

    if err := parent.Start(); err != nil {
        log.Panicln(err)
    }

    if err := RecordContainerInfo(parent.Process.Pid, cmdArray, containerName, containerId, volume); err != nil {
        log.Panicf("[Run] RecordContainerInfo throws error: %v", err)
    }

    cgroupManager := cgroups.NewCgroupManager("yyfdocker-cgroup")
    cgroupManager.Set(cfg)
    cgroupManager.Apply(parent.Process.Pid)

    rawCommand := strings.Join(cmdArray, " ")
    fmt.Printf("[Run] cmd: %v\n", rawCommand)
    writePipe.WriteString(rawCommand)
    writePipe.Close()

    if tty {
        parent.Wait()
        DeleteContainerInfo(containerName)
        container.DeleteWorkSpace(volume, containerName)
        cgroupManager.Destroy()
    }
}

func RecordContainerInfo(pid int, cmdArray []string, containerName string, id string, volume string) error {

    containerInfo := &container.ContainerInfo{
        Id:          id, 
        Pid:         strconv.Itoa(pid), 
        Command:     strings.Join(cmdArray, " "), 
        CreatedTime: time.Now().Format("15:01:01"),
        Status:      container.RUNNING,
        Name:        containerName,
        Volume:      volume,
    }

    return WriteContainerInfo(containerInfo)
}

func randString(n int) string {
    letterBytes := "1234567890"
    rand.Seed(time.Now().UnixNano())
    b := make([]byte, n)
    for i := range b {
        b[i] = letterBytes[rand.Intn(10)]
    }
    return string(b)
}