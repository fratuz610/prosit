package process

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"prosit/alert"
	"strings"
	"time"
)

type internalProcess struct {
	id          string
	pid         int
	path        string
	argList     []string
	folder      string
	cmd         *exec.Cmd
	err         error
	terminate   bool
	runCount    int
	avgDuration int
	lastStarted int64
	stdout      *Consumer
	stderr      *Consumer
	alertID     string
	isRunning   bool
}

func newInternalProcess(id, run, folder, alertID string) (*internalProcess, error) {

	var err error

	runList := strings.Split(run, " ")

	runList[0], err = exec.LookPath(runList[0])

	if err != nil {
		return nil, fmt.Errorf("Unable to find programm '%s' to run", runList[0])
	}

	ret := &internalProcess{}
	ret.id = id
	ret.path = runList[0]
	ret.argList = runList[1:]
	ret.folder = folder
	ret.avgDuration = -1
	ret.alertID = alertID
	ret.isRunning = false

	return ret, nil
}

func (p *internalProcess) start() {

	go func(p *internalProcess) {

		for {
			var startProcess = time.Now()

			// we clean the flag
			p.err = nil

			log.Printf("Starting process '%s %s'...\n", p.path, strings.Join(p.argList, " "))

			p.lastStarted = time.Now().Unix()

			// we create the cmd structure every run
			p.cmd = exec.Command(p.path, p.argList...)

			// we set the startup folder
			p.cmd.Dir = p.folder

			// we set the out and err streams to something readable and autorotating
			p.stdout = &Consumer{}
			p.stderr = &Consumer{}

			p.cmd.Stdout = p.stdout
			p.cmd.Stderr = p.stderr

			p.err = p.cmd.Start()

			if p.err != nil {

				// we have an alert to send
				alert.SendAlert(p.alertID, "Unable to start process %s: %v\n", p.cmd.Path, p.err)

				break
			}

			alert.SendAlert(p.alertID, "Process %s started correctly at %v", p.cmd.Path, time.Now())

			log.Printf("Process %s' started\n", p.cmd.Path)

			// the process is running now
			p.isRunning = true

			// we save the PID
			p.pid = p.cmd.Process.Pid

			p.err = p.cmd.Wait()

			// the process is NOT running now
			p.isRunning = false

			// we have an alert to send
			if p.err != nil {
				alert.SendAlert(p.alertID, "Process '%s' exited with error %v", p.cmd.Path, p.err)
			} else {
				alert.SendAlert(p.alertID, "Process '%s' exited with NO error", p.cmd.Path)
			}

			if p.terminate {
				log.Printf("Terminate flag on, closing down\n")
				break
			}

			runDuration := time.Since(startProcess)

			if p.avgDuration == -1 {
				p.avgDuration = int(runDuration.Seconds())
			} else {
				p.avgDuration = (p.avgDuration*p.runCount + int(runDuration.Seconds())) / (p.runCount + 1)
			}

			p.runCount++

			log.Printf("Process '%s' has an average duration of %d seconds\n", p.cmd.Path, p.avgDuration)

			if p.avgDuration < 10 {
				alert.SendAlert(p.alertID, "Process '%s' keeps on crashing (average life %d seconds). Shutting down.", p.cmd.Path, p.avgDuration)
				break
			}
		}

	}(p)
}

func (p *internalProcess) stop() error {
	if p.cmd.Process == nil {
		return fmt.Errorf("Unable to stop a process that hasn't started %s: ", p.cmd.Path)
	}

	// we set the terminate flag
	p.terminate = true

	// we send an ctrl+C signal
	return p.cmd.Process.Signal(os.Interrupt)

}
