package web

import (
	"fmt"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"os"
	"prosit/alert"
	"prosit/process"
	"strings"
)

type CreateProcessReq struct {
	Run     string `json:"run"`
	Folder  string `json:"folder"`
	User    string `json:"user"`
	AlertID string `json:"alertID"`
}

func verifyRequest(req *CreateProcessReq) error {

	// we do some cleanup
	req.Run = strings.Trim(req.Run, "\n\r\t ")
	req.Folder = strings.Trim(req.Folder, "\n\r\t ")
	req.User = strings.Trim(req.User, "\n\r\t ")
	req.AlertID = strings.Trim(req.AlertID, "\n\r\t ")

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

	err := verifyRequest(&req)

	if err != nil {
		outputError(err, r)
		return
	}

	err = process.AddProcess(req.Run, req.Folder, req.AlertID)

	if err != nil {
		outputError(err, r)
		return
	}

	r.JSON(202, "")
}
