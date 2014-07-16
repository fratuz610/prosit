package web

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"prosit/process"
)

func getProcessLogs(params martini.Params, r render.Render) {

	logList, err := process.GetProcessLogs(params["processID"])

	if err != nil {
		outputError(err, r)
		return
	}

	r.JSON(200, logList)
}
