package cl

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"prosit/process"
	"time"
)

func GetProcessLogs() error {

	processID := readLine("Process ID", "")

	justStarted := true
	lastItemTime := int64(0)

	// infinite cycle until we presst ctrl+c
	for {

		// we wait 1/2 sec
		time.Sleep(500 * time.Millisecond)

		// we check the logs
		resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:9999/processes/%s/logs", processID))

		if err != nil {
			return err
		}

		data, err := ioutil.ReadAll(resp.Body)

		resp.Body.Close()

		if err != nil {
			return err
		}

		logItemList := make([]process.LogItem, 0)

		err = json.Unmarshal(data, &logItemList)

		if err != nil {
			return err
		}

		// if the process has an empty log, continue
		if len(logItemList) == 0 {
			continue
		}

		if justStarted {

			// we check if we can show at least 10 entriess
			maxIndex := 10
			if len(logItemList) < 10 {
				maxIndex = 0
			}

			// we show the last 10 entries if available
			for i := maxIndex - 1; i >= 0; i-- {
				logItem := logItemList[i]
				fmt.Printf("%v : %s\n", logItem.Time, logItem.Message)
				lastItemTime = logItem.Time
			}

			// we reset the just started flag
			justStarted = false
		} else {

			// we show only new items since last update
			for i := len(logItemList) - 1; i >= 0; i-- {
				logItem := logItemList[i]

				if logItem.Time > lastItemTime {
					fmt.Printf("%v : %s\n", logItem.Time, logItem.Message)
					lastItemTime = logItem.Time
				}
			}
		}

	}

	return nil
}
