package web

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"prosit/err"
	"prosit/process"
	"strconv"
)

func getProcessLogs(params martini.Params, r render.Render) {

	processID, cerr := strconv.Atoi(params["processID"])

	if cerr != nil {
		outputError(err.NewBadRequestError("Unable to convert %s into a number: %v", params["processID"], cerr), r)
		return
	}

	logList, cerr := process.GetProcessLogs(processID)

	if cerr != nil {
		outputError(cerr, r)
		return
	}

	r.JSON(200, logList)
}
