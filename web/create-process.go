package web

import (
	"fmt"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"log"
	"os"
	"prosit/alert"
	"prosit/process"
	"strings"
)

type CreateProcessReq struct {
	Id      string `json:"id"`
	Run     string `json:"run"`
	Folder  string `json:"folder"`
	User    string `json:"user"`
	AlertID string `json:"alertID"`
}

func verifyCreateProcessRequest(req *CreateProcessReq) error {

	// we do some cleanup
	req.Id = strings.Trim(req.Id, "\n\r\t ")
	req.Run = strings.Trim(req.Run, "\n\r\t ")
	req.Folder = strings.Trim(req.Folder, "\n\r\t ")
	req.User = strings.Trim(req.User, "\n\r\t ")
	req.AlertID = strings.Trim(req.AlertID, "\n\r\t ")

	if req.Id == "" {
		cnt := 0
		for {
			req.Id = generateID(req.Run, cnt)
			cnt++
			if !process.ProcessExists(req.Id) {
				break
			}
		}
	}

	log.Printf("Process ID: " + req.Id)

	if req.Run == "" {
		return fmt.Errorf("No run command passed!")
	}

	if req.Folder == "" {
		req.Folder, _ = os.Getwd()
	}

	if _, err := os.Stat(req.Folder); os.IsNotExist(err) {
		return fmt.Errorf("Non existent run folder: '%s'", req.Folder)
	}

	if req.AlertID != "" && !alert.AlertExists(req.AlertID) {
		return fmt.Errorf("Non existent alertID: '%s'", req.AlertID)
	}

	return nil
}

func createProcess(req CreateProcessReq, params martini.Params, r render.Render) {

	err := verifyCreateProcessRequest(&req)

	if err != nil {
		outputError(err, r)
		return
	}

	err = process.AddProcess(req.Id, req.Run, req.Folder, req.AlertID)

	if err != nil {
		outputError(err, r)
		return
	}

	r.JSON(202, "")
}

func generateID(run string, counter int) string {
	ret := strings.Replace(run, " ", "", -1)
	ret = strings.Replace(ret, ".", "", -1)
	ret = strings.Replace(ret, "_", "", -1)
	ret = strings.Replace(ret, ":", "", -1)
	ret = strings.Replace(ret, "\\", "", -1)
	ret = strings.Replace(ret, "/", "", -1)

	return fmt.Sprintf("%s-%03d", ret[0:8], counter)
}
