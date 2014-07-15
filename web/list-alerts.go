package web

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"prosit/alert"
)

func listAlerts(params martini.Params, r render.Render) {

	alertList := alert.ListAlerts()

	r.JSON(200, alertList)
}
