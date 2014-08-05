package cl

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"prosit/process"
)

func getProcesses() ([]process.Process, error) {

	processList := make([]process.Process, 0)

	err := getJSON("http://127.0.0.1:9999/processes", &processList)

	if err != nil {
		return nil, err
	}

	return processList, nil
}

func getJSON(url string, target interface{}) error {

	resp, err := http.Get(url)

	if err != nil {
		return err
	}

	data, err := ioutil.ReadAll(resp.Body)

	resp.Body.Close()

	if err != nil {
		return err
	}

	err = json.Unmarshal(data, target)

	if err != nil {
		return err
	}

	return nil
}
