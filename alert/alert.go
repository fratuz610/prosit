package alert

import (
	"container/list"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"sync"
	"time"
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
var persistFile string

func initialize() {
	l = list.New()

	// first we try the current folder
	persistFile = fmt.Sprintf("%s%cprosit-alerts.json", filepath.Dir(os.Args[0]), os.PathSeparator)

	data, err := ioutil.ReadFile(persistFile)
	if err != nil {

		user, err := user.Current()

		if err != nil {
			log.Printf("Unable to determine current running user info: %v", err)
			return
		}

		// then we try the home folder
		persistFile = fmt.Sprintf("%s%cprosit-alerts.json", user.HomeDir, os.PathSeparator)

		data, err = ioutil.ReadFile(persistFile)

		if err != nil {
			log.Printf("Unable to find/read persistance file '%s'", persistFile)
			return
		}
	}

	tempAlertList := make([]Alert, 0)
	err = json.Unmarshal(data, &tempAlertList)

	if err != nil {
		log.Printf("Unable to json decode persistance file '%s': %v", persistFile, err)
		return
	}

	lMutex.Lock()
	defer lMutex.Unlock()

	for _, tempAlert := range tempAlertList {
		l.PushBack(&Alert{tempAlert.Id, "mailgun", tempAlert.FromEmail, tempAlert.ToEmailList, tempAlert.ApiKey, tempAlert.Domain})
	}

	// we setup the email delivery service in another go routine
	go func() {

		// we start the email delivery service
		c := time.Tick(1 * time.Minute)
		for _ = range c {
			deliverAlerts()
		}
	}()
}

func AddAlert(id, fromEmail string, toEmailList []string, apiKey, domain string) error {

	if l == nil {
		initialize()
	}

	lMutex.Lock()
	defer lMutex.Unlock()

	newAlert := &Alert{id, "mailgun", fromEmail, toEmailList, apiKey, domain}
	l.PushBack(newAlert)

	go persistAlerts()
	return nil
}

func ListAlerts() []Alert {

	if l == nil {
		initialize()
	}

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

	if l == nil {
		initialize()
	}

	lMutex.RLock()
	defer lMutex.RUnlock()

	// Iterate through list and print its contents.
	for e := l.Front(); e != nil; e = e.Next() {
		var tmpAlert = e.Value.(*Alert)

		if tmpAlert.Id == id {
			return true
		}
	}

	return false
}

func DeleteAlert(id string) bool {

	if l == nil {
		initialize()
	}

	lMutex.Lock()
	defer lMutex.Unlock()

	// Iterate through list and print its contents.
	for e := l.Front(); e != nil; e = e.Next() {
		var tmpAlert = e.Value.(*Alert)

		if tmpAlert.Id == id {
			l.Remove(e)
			go persistAlerts()
			return true
		}
	}

	return false
}

func persistAlerts() {

	alertList := ListAlerts()

	data, err := json.Marshal(alertList)

	if err != nil {
		log.Printf("Unable to json encode alert list: %v", err)
		return
	}

	err = ioutil.WriteFile(persistFile, data, 0644)

	if err != nil {
		log.Printf("Unable to write alert list file: %v", err)
		return
	}
}
