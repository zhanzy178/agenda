package auth

import (
	"log"
	"os/exec"
	"strconv"
	"strings"
)

func CurrentBashPid() int {
	// This function will return current bash main pid.
	// It will parse 'ps' result to get last bash process pid
	cmd := exec.Command("ps", "--sort", "start_time", "-o", "pid")
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	outStr := strings.Split(string(out), "\n")
	pidStr := strings.Trim(outStr[1], " ")
	pid, err := strconv.Atoi(pidStr)

	if err != nil {
		log.Fatal(err)
	}
	return pid
}
