package web

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"prosit/cerr"
	"prosit/process"
	"strconv"
)

func getProcessErrors(params martini.Params, r render.Render) {

	processID, err := strconv.Atoi(params["processID"])

	if err != nil {
		outputError(cerr.NewBadRequestError("Unable to convert %s into a number: %v", params["processID"], err), r)
		return
	}

	logList, err := process.GetProcessErrors(processID)

	if err != nil {
		outputError(err, r)
		return
	}

	r.JSON(200, logList)
}
