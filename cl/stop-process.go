package cl

import (
	"fmt"
	"net/http"
)

func StopProcessCL() error {

	processToStop := readLine("Process ID", "")

	resp, err := http.Post(fmt.Sprintf("http://127.0.0.1:9999/processes/%s/stop", processToStop), "", nil)

	if err != nil {
		return err
	}

	if resp.StatusCode != 202 {
		return fmt.Errorf("Server returned status code %d", resp.StatusCode)
	}

	fmt.Printf("Process '%s' stopped successfully\n\n", processToStop)

	return nil
}
