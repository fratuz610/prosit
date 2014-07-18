package launch

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func Launch(user, run string) {

	runList := strings.Split(run, " ")

	exe, err := exec.LookPath(runList[0])

	if err != nil {
		panic(fmt.Errorf("Unable to find programm '%s' to run", runList[0]))
	}

	argList := runList[1:]

	// we create the program
	cmd := exec.Command(exe, argList...)

	// we redirect stdout and err to ourself
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// we run the executable (in the main thread)
	err = cmd.Run()

	if err != nil {
		panic(err)
	}
}
