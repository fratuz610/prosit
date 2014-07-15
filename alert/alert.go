package alert

import (
	"container/list"
	"fmt"
	"log"
	"net/mail"
	"sync"
)

type alert struct {
	id          string
	alertType   string
	fromEmail   string
	toEmailList []string
	apiKey      string
	domain      string
}

var l *list.List
var lMutex sync.RWMutex

func init() {
	l = list.New()
}

func AddAlert(id, fromEmail string, toEmailList []string, apiKey, domain string) error {

	if AlertExists(id) {
		return fmt.Errorf("Alert with id '%s' already exists")
	}

	if _, err := mail.ParseAddress(fromEmail); err != nil {
		return fmt.Errorf("Invalid from email address '%s': %v", fromEmail, err)
	}

	for _, toEmail := range toEmailList {
		if _, err := mail.ParseAddress(toEmail); err != nil {
			return fmt.Errorf("Invalid to email address '%s': %v", toEmail, err)
		}
	}

	lMutex.Lock()
	defer lMutex.Unlock()

	newAlert := &alert{id, "mailgun", fromEmail, toEmailList, apiKey, domain}
	l.PushBack(newAlert)

	return nil
}

func AlertExists(id string) bool {

	lMutex.RLock()
	defer lMutex.RUnlock()

	// Iterate through list and print its contents.
	for e := l.Front(); e != nil; e = e.Next() {
		var tmpAlert = e.Value.(alert)

		if tmpAlert.id == id {
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
		var tmpAlert = e.Value.(alert)

		if tmpAlert.id == id {
			l.Remove(e)
			return true
		}
	}

	return false
}

func SendAlert(alertID, message string) {
	log.Printf("Got alert for '%s' and message '%s'", alertID, message)
}
