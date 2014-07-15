package alert

import (
	"container/list"
	"log"
	"net/mail"
	"prosit/cerr"
	"sync"
)

type Alert struct {
	Id          string   `json:"id"`
	AlertType   string   `json:"alertType"`
	FromEmail   string   `json:"fromEmail"`
	ToEmailList []string `json:"toEmailList"`
	ApiKey      string   `json:"apiKey"`
	Domain      string   `json:"domain"`
}

var l *list.List
var lMutex sync.RWMutex

func init() {
	l = list.New()
}

func AddAlert(id, fromEmail string, toEmailList []string, apiKey, domain string) error {

	if AlertExists(id) {
		return cerr.NewBadRequestError("Alert with id '%s' already exists", id)
	}

	if _, err := mail.ParseAddress(fromEmail); err != nil {
		return cerr.NewBadRequestError("Invalid from email address '%s': %v", fromEmail, err)
	}

	for _, toEmail := range toEmailList {
		if _, err := mail.ParseAddress(toEmail); err != nil {
			return cerr.NewBadRequestError("Invalid to email address '%s': %v", toEmail, err)
		}
	}

	lMutex.Lock()
	defer lMutex.Unlock()

	newAlert := &Alert{id, "mailgun", fromEmail, toEmailList, apiKey, domain}
	l.PushBack(newAlert)

	return nil
}

func ListAlerts() []Alert {

	lMutex.RLock()
	defer lMutex.RUnlock()

	ret := make([]Alert, 0)

	// Iterate through list and print its contents.
	for e := l.Front(); e != nil; e = e.Next() {
		var tmpAlert = e.Value.(*Alert)
		ret = append(ret, *tmpAlert)
	}

	return ret
}

func AlertExists(id string) bool {

	lMutex.RLock()
	defer lMutex.RUnlock()

	// Iterate through list and print its contents.
	for e := l.Front(); e != nil; e = e.Next() {
		var tmpAlert = e.Value.(Alert)

		if tmpAlert.Id == id {
			return true
		}
	}

	return false
}

func DeleteAlert(id string) bool {
	lMutex.Lock()
	defer lMutex.Unlock()

	// Iterate through list and print its contents.
	for e := l.Front(); e != nil; e = e.Next() {
		var tmpAlert = e.Value.(Alert)

		if tmpAlert.Id == id {
			l.Remove(e)
			return true
		}
	}

	return false
}

func SendAlert(alertID, message string) {
	log.Printf("Got alert for '%s' and message '%s'", alertID, message)
}
