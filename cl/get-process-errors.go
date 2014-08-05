package cl

import (
	"fmt"
	"prosit/process"
	"time"
)

func GetProcessErrors() error {

	// we get the process ID
	processID, err := getUserProcessID()

	if err != nil {
		return err
	}

	justStarted := true
	lastItemTime := int64(0)

	// infinite cycle until we presst ctrl+c
	for {

		// we wait 1/2 sec
		time.Sleep(500 * time.Millisecond)

		errorItemList := make([]process.LogItem, 0)

		err := getJSON(fmt.Sprintf("http://127.0.0.1:9999/processes/%s/errors", processID), &errorItemList)

		if err != nil {
			return err
		}

		// if the process has an empty log, continue
		if len(errorItemList) == 0 {
			continue
		}

		if justStarted {

			// we check if we can show at least 10 entriess
			maxIndex := 10
			if len(errorItemList) < 10 {
				maxIndex = 0
			}

			// we show the last 10 entries if available
			for i := maxIndex - 1; i >= 0; i-- {
				logItem := errorItemList[i]
				fmt.Printf("%v: %s\n", formatTime(logItem.Time), logItem.Message)
				lastItemTime = logItem.Time
			}

			// we reset the just started flag
			justStarted = false
		} else {

			// we show only new items since last update
			for i := len(errorItemList) - 1; i >= 0; i-- {
				logItem := errorItemList[i]

				if logItem.Time > lastItemTime {
					fmt.Printf("%v: %s\n", formatTime(logItem.Time), logItem.Message)
					lastItemTime = logItem.Time
				}
			}
		}

	}

	return nil
}
