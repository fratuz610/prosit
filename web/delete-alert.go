package web

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"prosit/alert"
	"prosit/cerr"
)

func deleteAlert(params martini.Params, r render.Render) {

	alertID := params["alertID"]

	if !alert.AlertExists(alertID) {
		outputError(cerr.NewNotFoundError("No alert with id '%s'", alertID), r)
		return
	}

	alert.DeleteAlert(alertID)

	r.JSON(200, "")
}
