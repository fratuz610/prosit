package web

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"prosit/process"
)

func getProcessErrors(params martini.Params, r render.Render) {

	logList, err := process.GetProcessErrors(params["processID"])

	if err != nil {
		outputError(err, r)
		return
	}

	r.JSON(200, logList)
}
