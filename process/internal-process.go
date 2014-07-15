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
	pid         int
	cmd         *exec.Cmd
	err         error
	terminate   bool
	runCount    int
	avgDuration int
	lastStarted int64
	stdout      *Consumer
	stderr      *Consumer
	alertID     string
}

func newInternalProcess(run, folder, alertID string) (*internalProcess, error) {

	var err error

	runList := strings.Split(run, " ")

	runList[0], err = exec.LookPath(runList[0])

	if err != nil {
		return nil, fmt.Errorf("Unable to find programm '%s' to run", runList[0])
	}

	run = strings.Join(runList, " ")

	log.Printf("Final run string '%s'\n", run)

	ret := &internalProcess{}
	ret.cmd = exec.Command(runList[0], runList[1:]...)
	ret.cmd.Dir = folder
	ret.avgDuration = -1
	ret.alertID = alertID

	return ret, nil
}

func (p *internalProcess) start() {

	go func(p *internalProcess) {

		for {
			var startProcess = time.Now()

			// we clean the flag
			p.err = nil

			log.Printf("Starting internalProcess '%s'...\n", p.cmd.Path)

			p.lastStarted = time.Now().Unix()

			// we set the out and err streams to something readable and autorotating
			p.stdout = &Consumer{}
			p.stderr = &Consumer{}

			p.cmd.Stdout = p.stdout
			p.cmd.Stderr = p.stderr

			p.err = p.cmd.Start()

			if p.err != nil {

				log.Printf("Unable to start process '%s' because: %v\n", p.cmd.Path, p.err)

				if p.alertID != "" {
					// we have an alert to send
					alert.SendAlert(p.alertID, fmt.Sprintf("Unable to start process %s: %v\n", p.cmd.Path, p.err))
				}

				break
			}

			log.Printf("Process %s' started\n", p.cmd.Path)

			// we save the PID
			p.pid = p.cmd.Process.Pid

			p.err = p.cmd.Wait()

			log.Printf("Process '%s' exited with error %v\n", p.err)

			if p.alertID != "" {
				// we have an alert to send
				alert.SendAlert(p.alertID, fmt.Sprintf("Process '%s' exited with error (if any) %v\n", p.err))
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

			if p.avgDuration < 10 {
				log.Printf("Average life %d too low (< 10 seconds), something is wrong, closing down\n", p.avgDuration)
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
