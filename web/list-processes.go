package web

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"prosit/process"
)

func listProcesses(params martini.Params, r render.Render) {

	processList := process.ListProcesses()

	r.JSON(200, processList)
}
