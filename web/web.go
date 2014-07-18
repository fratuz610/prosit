package web

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"net/http"
	"os"
	"strconv"
)

func StartWeb(port int) {
	// we define the HTTP port
	os.Setenv("PORT", strconv.Itoa(port))

	m := martini.Classic()

	// JSON middleware
	m.Use(render.Renderer(render.Options{IndentJSON: true}))

	// add json content type to each POST and PUT requests
	m.Use(addJSONContentTypeMW)

	m.Get("/processes", listProcesses)
	m.Post("/processes", binding.Bind(CreateProcessReq{}), createProcess)

	m.Post("/processes/:processID/stop", stopProcess)
	m.Post("/processes/:processID/restart", restartProcess)

	m.Get("/processes/:processID/logs", getProcessLogs)
	m.Get("/processes/:processID/errors", getProcessErrors)

	m.Get("/alerts", listAlerts)
	m.Post("/alerts", binding.Bind(CreateAlertReq{}), createAlert)
	m.Delete("/alerts/:alertID", deleteAlert)

	m.Get("/internal-log", getInternalLog)

	//	m.Put("/cc/bill2bill/recycling-cassette", binding.Bind(ccbill.SetRecyclingCassetteTypeRequest{}), setRecyclingCassette)

	m.Run()
}

func addJSONContentTypeMW(res http.ResponseWriter, req *http.Request) {

	if req.Method == "POST" || req.Method == "PUT" {
		req.Header.Set("content-type", "application/json")
	}
}
