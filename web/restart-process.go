package web

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"prosit/cerr"
	"prosit/process"
)

func restartProcess(params martini.Params, r render.Render) {

	if !process.ProcessExists(params["processID"]) {
		outputError(cerr.NewNotFoundError("Process %s not found", params["processID"]), r)
		return
	}

	err := process.RestartProcess(params["processID"])

	if err != nil {
		outputError(err, r)
		return
	}

	r.JSON(202, "")
}
