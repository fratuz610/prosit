package launch

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

func Launch(req *LaunchRequest) {

	// we need to be root
	if os.Getuid() != 0 {
		panic(fmt.Errorf("This program must run as root when launching processes. Got UID %d", os.Geteuid()))
	}

	runList := strings.Split(req.FullPath, " ")

	exe, err := exec.LookPath(runList[0])

	if err != nil {
		panic(fmt.Errorf("Unable to find programm '%s' to run", runList[0]))
	}

	argList := runList[1:]

	// we create the program
	cmd := exec.Command(exe, argList...)

	if req.Folder != "" {
		cmd.Dir = req.Folder
	}

	// we redirect stdout and err to ourself
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	uid, err := getUidFromUsername(req.RunAs)

	if err != nil {
		panic(fmt.Errorf("Unable to determine the UID for user '%s': %v", req.RunAs, err))
	}

	// we set the UID
	if err = syscall.Setuid(uid); err != nil {
		panic(fmt.Errorf("Unable to set the UID: %v", err))
	}

	// we run the executable (in the main thread)
	err = cmd.Run()

	if err != nil {
		panic(err)
	}
}

func getUidFromUsername(user string) (int, error) {

	cmd := exec.Command("id", "-u", user)

	outputData, err := cmd.CombinedOutput()

	if err != nil {
		return -1, fmt.Errorf("Unable to run id command: %v", err)
	}

	// we clean up the output
	output := strings.Trim(string(outputData), "\n\t\r ")

	uid, err := strconv.Atoi(output)

	if err != nil {
		return -1, fmt.Errorf("id returned not a number: %v", err)
	}

	return uid, nil
}
