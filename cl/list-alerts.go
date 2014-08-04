package cl

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func ListAlertsCL() error {

	resp, err := http.Get("http://127.0.0.1:9999/alerts")

	if err != nil {
		return err
	}

	data, err := ioutil.ReadAll(resp.Body)

	resp.Body.Close()

	if err != nil {
		return err
	}

	fmt.Printf("Alert List:\n%s\n", string(data))

	return nil
}
