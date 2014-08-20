package cl

import (
	"fmt"
)

var header string = `
Prosit command line interface:

`

var help string = header + `Usage

	prosit process start|add|create|new|list|logs|log|errors|stop|terminate|restart|help
	prosit alert help list|create|add|new|delete|remove|help

`

var processHelp string = header + `Usage
	
	prosit process start|add|create|new -> starts a new process
	prosit process list -> returns the list of running managed processes
	prosit process logs|log -> shows the stdout output of a process - works like tail -f
	prosit process errors -> shows the stderr output of a process - works like tail -f
	prosit process stop|terminate -> terminates a managed process
	prosit process restart -> restarts a managed process
	
`

var alertHelp string = header + `Usage
	
	prosit alert list -> lists all alerts
	prosit alert create|add|new -> creates a new alert
	prosit alert delete|remove -> deletes an alert
	
`

func ProcessHelpCL() error {
	fmt.Printf(processHelp)
	return nil
}

func AlertHelpCL() error {
	fmt.Printf(alertHelp)
	return nil
}

func HelpCL() error {
	fmt.Printf(help)
	return nil
}
