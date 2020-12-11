package main

import (
	"os"
	"log"
)

const usage = "simple implementation by Yufeng Yang!"

func init() {
	logFileName := "YYFdocker.log"
	logFile, logErr := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if logErr != nil {
        log.Fatalf("logFile open fail: %v", logErr)
    }
    log.SetOutput(logFile)
    log.SetFlags(log.Lshortfile|log.LstdFlags)
}

func main() {
	log.Printf("==== MAIN FUNCTION START TO RUN ====\n")
	if len(os.Args) == 1 {
		log.Fatal("YYFdocker needs a command")
	}

	switch os.Args[1] {
	case "run":
		Run(true, os.Args[3])

	case "init":
		err := container.RunContainerInitProcess(cmd, nil)
		if err != nil {
			log.Fatal("Error in Init Function")
			Printf("Error in Init Function\n")
		}
		
	case "--usage":
		Printf("Usage: %s", usage)

	default:
		log.Fatal("Wrong arguments")
		Printf("Wrong arguments\nUsage: %s\n", usage)
	}

}