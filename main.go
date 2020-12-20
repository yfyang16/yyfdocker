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
    logFile, logErr := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY, 0666)
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
        allCfg := &subsystems.ResourceConfig{}
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
                default:          cmdArray = append(cmdArray, os.Args[argIdx:]...); argIdx = len(os.Args)
            }
        }

        Run(tty, cmdArray, allCfg, volume)

    case "init":
        log.Printf("==== INIT START ====\n")
        err := container.RunContainerInitProcess()
        if err != nil {
            log.Fatal("Error in Init Function")
            fmt.Printf("Error in Init Function\n")
        }

    case "commit":
        log.Printf("==== COMMIT START ====\n")
        CommitContainer(os.Args[2])

    case "--usage":
        fmt.Printf("Usage: %s", usage)

    default:
        log.Fatal("Wrong arguments")
        fmt.Printf("Wrong arguments\nUsage: %s\n", usage)
    }

}
