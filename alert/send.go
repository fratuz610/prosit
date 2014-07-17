package alert

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

var notificationMap map[string][]string
var nMutex sync.RWMutex

func SendAlert(alertID, format string, a ...interface{}) {

	nMutex.Lock()
	defer nMutex.Unlock()

	message := fmt.Sprintf(format, a...)

	log.Printf("ALERT: *** %s ***", message)

	if alertID == "" {
		return
	}

	// lazy init
	if notificationMap == nil {
		notificationMap = make(map[string][]string)
	}

	if _, ok := notificationMap[alertID]; !ok {
		notificationMap[alertID] = make([]string, 0)
	}

	// we save to the notification map
	notificationMap[alertID] = append(notificationMap[alertID], message)
}

func deliverAlerts() {

	nMutex.Lock()
	defer nMutex.Unlock()

	// no notification, return
	if notificationMap == nil {
		return
	}

	for _, alert := range ListAlerts() {

		if _, ok := notificationMap[alert.Id]; !ok {
			return
		}

		subject, body := getEmailTemplate(notificationMap[alert.Id])

		// we send the email in another goroutine
		go sendMailgunEmail(alert, subject, body)

		// successful delivery
		delete(notificationMap, alert.Id)

	}
}

func sendMailgunEmail(alert Alert, subject, body string) {

	apiURL := fmt.Sprintf("https://api.mailgun.net/v2/%s/messages", alert.Domain)

	payload := url.Values{}

	payload.Set("from", alert.FromEmail)

	for _, to := range alert.ToEmailList {
		payload.Add("to", to)
	}

	payload.Set("subject", subject)
	payload.Set("text", body)

	req, err := http.NewRequest("POST", apiURL, strings.NewReader(payload.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth("api", alert.ApiKey)
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Printf("sendMailgunEmail: Http error: %v", err)
		return
	}
	responseBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		log.Printf("sendMailgunEmail: email not sent: response code %d / body: '%s'", resp.StatusCode, string(responseBody))
		return
	}
}
