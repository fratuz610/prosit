package web

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"net/mail"
	"prosit/alert"
	"prosit/cerr"
)

type CreateAlertReq struct {
	Id          string   `json:"id"`
	FromEmail   string   `json:"fromEmail"`
	ToEmailList []string `json:"toEmailList"`
	ApiKey      string   `json:"apiKey"`
	Domain      string   `json:"domain"`
}

func verifyCreateAlertRequest(req *CreateAlertReq) error {

	if alert.AlertExists(req.Id) {
		return cerr.NewBadRequestError("Alert with id '%s' already exists", req.Id)
	}

	if _, err := mail.ParseAddress(req.FromEmail); err != nil {
		return cerr.NewBadRequestError("Invalid from email address '%s': %v", req.FromEmail, err)
	}

	for _, toEmail := range req.ToEmailList {
		if _, err := mail.ParseAddress(toEmail); err != nil {
			return cerr.NewBadRequestError("Invalid to email address '%s': %v", toEmail, err)
		}
	}

	return nil
}

func createAlert(req CreateAlertReq, params martini.Params, r render.Render) {

	err := verifyCreateAlertRequest(&req)

	if err != nil {
		outputError(err, r)
		return
	}

	err = alert.AddAlert(req.Id, req.FromEmail, req.ToEmailList, req.ApiKey, req.Domain)

	if err != nil {
		outputError(err, r)
		return
	}

	r.JSON(202, "")
}
