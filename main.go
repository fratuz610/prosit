// prosit project main.go
package main

import (
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"prosit/cl"
	"prosit/ilog"
	"prosit/launch"
	"prosit/web"
	"strings"
)

func main() {

	lrs := os.Getenv("_PROSIT_LAUNCH_REQ")

	if lrs != "" {
		// we have a launch request

		lr := &launch.LaunchRequest{}
		d := gob.NewDecoder(base64.NewDecoder(base64.StdEncoding, strings.NewReader(lrs)))
		err := d.Decode(lr)

		if err != nil {
			log.Fatalf("Failed to decode LaunchRequest in child: %v", err)
		}

		launch.Launch(lr)

		os.Exit(0)
	}

	if len(os.Args) == 1 {

		// we redirect stdout/err
		ilog.RedirectOutput()
		log.SetOutput(ilog.GetWriter())

		log.Printf("Starting as daemon process\n")
		web.StartWeb(9999)
		return
	}

	var err error

	switch strings.ToLower(os.Args[1]) {
	case "start-process", "add-process":
		err = cl.StartProcessCL()
	case "list-processes":
		err = cl.ListProcessesCL()
	case "logs":
		err = cl.GetProcessLogs()
	case "stop-process":
		err = cl.StopProcessCL()
	case "restart-process":
		err = cl.RestartProcessCL()
	case "list-alerts":
		err = cl.ListAlertsCL()
	case "delete-alert":
		err = cl.DeleteAlertCL()
	case "create-alert":
		err = cl.CreateAlertCL()
	}

	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}
}
