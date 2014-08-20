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

		// we need to be root
		if os.Getuid() != 0 {
			fmt.Printf("ERROR: This program must run as root in service mode. Got UID %d\n", os.Geteuid())
			os.Exit(2)
		}

		log.Printf("Service mode detected, redirecting output\n")

		// we redirect stdout/err
		ilog.RedirectOutput()
		log.SetOutput(ilog.GetWriter())

		log.Printf("Starting as daemon process\n")
		web.StartWeb(9999)
		return
	}

	var err error

	if len(os.Args) <= 2 {
		fmt.Printf("ERROR: missing command line arguments\n")
		return
	}

	switch strings.ToLower(os.Args[1]) {
	case "process":
		switch strings.ToLower(os.Args[2]) {
		case "start", "add", "create", "new":
			err = cl.StartProcessCL()
		case "list":
			err = cl.ListProcessesCL()
		case "logs", "log":
			err = cl.GetProcessLogs()
		case "errors":
			err = cl.GetProcessErrors()
		case "stop", "terminate":
			err = cl.StopProcessCL()
		case "restart":
			err = cl.RestartProcessCL()
		default:
			err = cl.ProcessHelpCL()
		}
	case "alert":
		switch strings.ToLower(os.Args[2]) {
		case "list":
			err = cl.ListAlertsCL()
		case "delete", "remove":
			err = cl.DeleteAlertCL()
		case "create", "add", "new":
			err = cl.CreateAlertCL()
		default:
			err = cl.AlertHelpCL()
		}
	default:
		err = cl.HelpCL()
	}

	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}
}
