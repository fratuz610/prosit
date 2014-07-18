package cl

import (
	"fmt"
	"net/http"
)

func RestartProcessCL() error {

	processToRestart := readLine("Process ID", "")

	resp, err := http.Post(fmt.Sprintf("http://127.0.0.1:9999/processes/%s/restart", processToRestart), "", nil)

	if err != nil {
		return err
	}

	if resp.StatusCode != 202 {
		return fmt.Errorf("Server returned status code %d", resp.StatusCode)
	}

	fmt.Printf("Process '%s' restarted successfully\n\n", processToRestart)

	return nil
}
