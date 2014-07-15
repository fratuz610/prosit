package web

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"prosit/cerr"
	"prosit/process"
	"strconv"
)

func getProcessLogs(params martini.Params, r render.Render) {

	processID, err := strconv.Atoi(params["processID"])

	if err != nil {
		outputError(cerr.NewBadRequestError("Unable to convert %s into a number: %v", params["processID"], err), r)
		return
	}

	logList, err := process.GetProcessLogs(processID)

	if err != nil {
		outputError(err, r)
		return
	}

	r.JSON(200, logList)
}
