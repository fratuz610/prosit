package cl

import (
	"fmt"
	"net/http"
)

func DeleteAlertCL() error {

	alertID := readLine("AlertID", "")

	req, err := http.NewRequest("DELETE", fmt.Sprintf("http://127.0.0.1:9999/alerts/%s", alertID), nil)

	// handle err
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("Server returned status code %d", resp.StatusCode)
	}

	fmt.Printf("Alert '%s' delete successfully\n\n", alertID)

	return nil
}
