package main

import (
    "fmt"
    "log"
    "os"
    "./cgroups/subsystems"
    "./container"
)

const usage = "yyfdocker run [-it] [-m val] [-cpushare val] [-cpuset val] [cmd]\n[simple implementation by Yufeng Yang]"

func init() {
    logFileName := "YYFdocker.log"
    logFile, logErr := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
    if logErr != nil {
        log.Fatalf("logFile open fail: %v", logErr)
    }
    log.SetOutput(logFile)
    log.SetFlags(log.Lshortfile | log.LstdFlags)
}

func main() {
    if len(os.Args) == 1 {
        log.Fatal("YYFdocker needs a command")
    }

    log.Printf("==**== YYFDOCKER START ==**== (%v)\n", os.Args)

    switch os.Args[1] {
    case "run":
        log.Printf("==== RUN START ====\n")

        tty := false
        detach := false
        allCfg := &subsystems.ResourceConfig{}
        containerName := ""
        var cmdArray []string
        var volume string

        argIdx := 2
        for argIdx < len(os.Args)  {
            arg := os.Args[argIdx]
            switch arg {
                case "-m":        allCfg.MemoryLimit = os.Args[argIdx + 1]; argIdx += 2
                case "-cpushare": allCfg.CpuShare = os.Args[argIdx + 1]; argIdx += 2
                case "-cpuset":   allCfg.CpuSet = os.Args[argIdx + 1]; argIdx += 2
                case "-v":        volume = os.Args[argIdx + 1]; argIdx += 2
                case "-it":       tty = true; argIdx += 1
                case "-d":        detach = true; argIdx += 1
                case "--name":    containerName = os.Args[argIdx + 1]; argIdx += 2
                default:          cmdArray = append(cmdArray, os.Args[argIdx:]...); argIdx = len(os.Args)
            }
        }
        imageName := cmdArray[0]
        cmdArray = cmdArray[1:]

        if tty && detach {
            log.Panicf("[RUN CMD] -d and -it parameters cannot be existed at the same time.")
        }

        Run(tty, cmdArray, allCfg, volume, containerName, imageName)

    case "init":
        log.Printf("==== INIT START ====\n")
        err := container.RunContainerInitProcess()
        if err != nil {
            log.Fatal("Error in Init Function")
            fmt.Printf("Error in Init Function\n")
        }

    case "commit":
        log.Printf("==== COMMIT START ====\n")
        CommitContainer(os.Args[2], os.Args[3])

    case "logs":
        log.Printf("==== LOGS START ====\n")
        LogContainer(os.Args[2])

    case "ps":
        log.Printf("==== PS START ====\n")
        ListContainers()

    case "stop":
        log.Printf("==== STOP START ====\n")
        StopContainer(os.Args[2])

    case "rm":
        log.Printf("==== RM START ====\n")
        RemoveContainer(os.Args[2])

    case "exec":
        log.Printf("==== EXEC START ====\n")

        if os.Getenv(ENV_EXEC_PID) != "" {
            log.Printf("[EXEC CMD] callback"); return
        }

        var cmdArray []string
        for _, cmd := range os.Args[3:] {
            cmdArray = append(cmdArray, cmd)
        }
        ExecContainer(os.Args[2], cmdArray)

    case "--usage":
        fmt.Printf("Usage: %s", usage)

    default:
        log.Fatal("Wrong arguments")
        fmt.Printf("Wrong arguments\nUsage: %s\n", usage)
    }

}
