package cl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"prosit/web"
	"strings"
)

func StartProcessCL() error {

	// we request the run
	run := readLine("Enter the command to run", "")

	// we verify the path
	runList := strings.Split(run, " ")

	var err error
	runList[0], err = exec.LookPath(runList[0])

	if err != nil {
		return fmt.Errorf("Unable to find programm '%s' to run", runList[0])
	}

	currentUser, _ := user.Current()

	// we request the folder
	runAs := readLine("Run as", currentUser.Username)

	runAsUser, err := user.Lookup(runAs)
	if err != nil {
		return fmt.Errorf("No user with id '%s'", runAs)
	}

	// we request the folder
	folder := readLine("Running folder", runAsUser.HomeDir)

	if _, err := os.Stat(folder); os.IsNotExist(err) {
		return fmt.Errorf("Folder '%s' does not exist", folder)
	}

	// we request the alertID
	alertID := readLine("AlertID", "")

	createProcessReq := &web.CreateProcessReq{}
	createProcessReq.Run = run
	createProcessReq.RunAs = runAs
	createProcessReq.Folder = folder
	createProcessReq.AlertID = alertID

	data, err := json.Marshal(createProcessReq)

	if err != nil {
		return err
	}

	buffer := bytes.NewBuffer(data)

	resp, err := http.Post("http://127.0.0.1:9999/processes", "application/json", buffer)

	if err != nil {
		return err
	}

	if resp.StatusCode != 202 {
		return fmt.Errorf("Server returned status code %d", resp.StatusCode)
	}

	fmt.Printf("Process '%s' added successfully\n\n", createProcessReq.Run)
	return nil
}
