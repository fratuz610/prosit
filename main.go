// prosit project main.go
package main

import (
	"flag"
	"log"
	"os"
	"prosit/web"
)

var daemon bool
var daemonPort int
var user string
var runCommand string
var runFolder string
var notifyID string

func init() {
	flag.BoolVar(&daemon, "daemon", true, "Run main daemon")
	flag.IntVar(&daemonPort, "daemonPort", 9999, "The local http port to start the daemon server on")

	flag.StringVar(&user, "user", "", "The user to run the process as")
	flag.StringVar(&runCommand, "run", "", "The command to run")
	flag.StringVar(&runFolder, "folder", "", "The run folder")
	flag.StringVar(&notifyID, "notify", "", "The id of the notification schema to use")
}

func main() {

	flag.Parse()

	for i, val := range os.Args {
		log.Printf("Argument: %d, value: %s", i, val)
	}

	log.Printf("daemon: %v", daemon)
	log.Printf("daemonPort: %v", daemonPort)
	log.Printf("user: %v", user)
	log.Printf("runCommand: %v", runCommand)

	if daemon {
		web.StartWeb(daemonPort)
	}
}
