package cl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/mail"
	"prosit/alert"
	"prosit/web"
	"strings"
)

func CreateAlertCL() error {

	req := &web.CreateAlertReq{}

	// we request the run
	req.Id = readLine("Enter the alertID", "")

	if alert.AlertExists(req.Id) {
		return fmt.Errorf("Alert %s already exits. Delete it first", req.Id)
	}

	req.FromEmail = readLine("From Email:", "")

	if _, err := mail.ParseAddress(req.FromEmail); err != nil {
		return fmt.Errorf("Invalid from email address '%s': %v", req.FromEmail, err)
	}

	toEmailCSV := readLine("To Email list (csv)", "")

	req.ToEmailList = strings.Split(toEmailCSV, ",")

	if len(req.ToEmailList) == 0 {
		return fmt.Errorf("Invalid to email list, no emails found")
	}

	for _, toEmail := range req.ToEmailList {
		if _, err := mail.ParseAddress(toEmail); err != nil {
			return fmt.Errorf("Invalid TO email address '%s': %v", toEmail, err)
		}
	}

	req.ApiKey = readLine("MailGun API Key", "")
	req.Domain = readLine("MailGun Domain Name", "")

	data, err := json.Marshal(req)

	if err != nil {
		return err
	}

	buffer := bytes.NewBuffer(data)

	resp, err := http.Post("http://127.0.0.1:9999/alerts", "application/json", buffer)

	if err != nil {
		return err
	}

	if resp.StatusCode != 202 {
		return fmt.Errorf("Server returned status code %d", resp.StatusCode)
	}

	fmt.Printf("Alert '%s' created successfully\n\n", req.Id)
	return nil
}
