package main

import (
	"fmt"
	"log"
	"log/syslog"
	"os"
	"os/exec"
	"strings"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("ERROR incorrect usage: onerror cmd [args...]")
		os.Exit(2)
	}

	// command to run and its args if any
	cmd := os.Args[1]
	cmdArgs := os.Args[2:]

	// Configure logger to write to the syslog NOTICE
	sysLog, err := syslog.New(syslog.LOG_NOTICE|syslog.LOG_CRON, cmd)
	if err != nil {
		log.Println("onerror: Error connecting to syslog: ", err)
		log.Println("onerror: Continuing executing, sending output to stderr")
	}

	stdOutStdErr, err := exec.Command(cmd, cmdArgs...).CombinedOutput()
	if err == nil {
		// Command was successful, output to syslog/stderr and exit
		if sysLog != nil {
			sysLog.Notice(fmt.Sprintf("onerror: SUCCESS runing %s %+v\n", cmd, cmdArgs))
			sysLog.Notice("onerror: Command output:\n")
			for _, line := range strings.Split(string(stdOutStdErr), "\n") {
				sysLog.Notice(fmt.Sprintf("onerror: %s\n", string(line)))
			}
		} else {
			log.Printf("onerror: SUCCESS running %s %+v\n", cmd, cmdArgs)
			log.Printf("onerror: Command output:\n%s\n", stdOutStdErr)
		}

		os.Exit(0)
	}

	// Command failed, print to all output to syslog if available, and stderr

	if sysLog != nil {
		sysLog.Err(fmt.Sprintf("onerror: ERROR running %s %+v: %s\n", cmd, cmdArgs, err.Error()))
		sysLog.Err("onerror: Command output:\n")
		for _, line := range strings.Split(string(stdOutStdErr), "\n") {
			sysLog.Err(fmt.Sprintf("onerror: %s\n", line))
		}
	}

	log.Printf("onerror: ERROR running %s %+v %s\n", cmd, cmdArgs, err.Error())
	log.Printf("onerror: Command output:\n%s\n", stdOutStdErr)

	// TODO: return exit value of cmd
	os.Exit(1)
}
