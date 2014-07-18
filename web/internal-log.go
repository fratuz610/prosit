package web

import (
	"github.com/martini-contrib/render"
	"prosit/ilog"
)

func getInternalLog(r render.Render) {
	r.JSON(200, ilog.GetOutput())
}
