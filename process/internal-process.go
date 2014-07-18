package process

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"os/exec"
	"prosit/alert"
	"prosit/cerr"
	"prosit/launch"
	"strings"
	"time"
)

type internalProcess struct {
	id            string
	pid           int
	path          string
	fullPath      string
	argList       []string
	folder        string
	cmd           *exec.Cmd
	err           error
	runAs         string
	terminate     bool
	runCount      int
	avgDuration   int
	lastStarted   int64
	stdout        *Consumer
	stderr        *Consumer
	alertID       string
	isRunning     bool
	isInterrupted bool
}

func newInternalProcess(id, run, folder, alertID, runAs string) (*internalProcess, error) {

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
	ret.fullPath = strings.Join(runList, " ")
	ret.folder = folder
	ret.avgDuration = -1
	ret.alertID = alertID
	ret.isRunning = false
	ret.runAs = runAs

	return ret, nil
}

func (p *internalProcess) start() {

	go func(p *internalProcess) {

		for {
			var startProcess = time.Now()

			// we clean the flag
			p.err = nil

			log.Printf("Process '%s': starting...\n", p.fullPath)

			p.lastStarted = time.Now().Unix()

			lr := &launch.LaunchRequest{RunAs: p.runAs, FullPath: p.fullPath, Folder: p.folder}

			var buf bytes.Buffer
			b64enc := base64.NewEncoder(base64.StdEncoding, &buf)
			p.err = gob.NewEncoder(b64enc).Encode(lr)
			b64enc.Close()

			if p.err != nil {
				// we have an alert to send
				alert.SendAlert(p.alertID, "Process '%s': unable to start because we can base64 encode ENV data: %v\n", p.fullPath, p.err)
				break
			}

			// we launch ourself with a specific ENV variable
			p.cmd = exec.Command(os.Args[0], "launch")

			p.cmd.Env = append(os.Environ(), "_PROSIT_LAUNCH_REQ="+buf.String())

			// we set the out and err streams to something readable and autorotating
			p.stdout = &Consumer{}
			p.stderr = &Consumer{}

			p.cmd.Stdout = p.stdout
			p.cmd.Stderr = p.stderr

			p.err = p.cmd.Start()

			if p.err != nil {

				// we have an alert to send
				alert.SendAlert(p.alertID, "Process '%s': unable to start because: %v\n", p.fullPath, p.err)

				break
			}

			alert.SendAlert(p.alertID, "Process '%s': started at %v", p.fullPath, time.Now())

			// the process is running now
			p.isRunning = true

			// we save the PID
			p.pid = p.cmd.Process.Pid

			p.err = p.cmd.Wait()

			// the process is NOT running now
			p.isRunning = false

			alert.SendAlert(p.alertID, "Process '%s': stopped at %v", p.fullPath, time.Now())

			if p.terminate {
				log.Printf("Process '%s': terminate flag on, breaking out of the loop\n", p.fullPath)
				break
			}

			if p.isInterrupted {
				alert.SendAlert(p.alertID, "Process '%s': manual stop detected", p.fullPath)
				p.isInterrupted = false
				continue
			}

			// we have an alert to send
			if p.err != nil {
				alert.SendAlert(p.alertID, "Process '%s': exited with error %v", p.fullPath, p.err)
			} else {
				alert.SendAlert(p.alertID, "Process '%s': exited with NO error", p.fullPath)
			}

			runDuration := time.Since(startProcess)

			if p.avgDuration == -1 {
				p.avgDuration = int(runDuration.Seconds())
			} else {
				p.avgDuration = (p.avgDuration*p.runCount + int(runDuration.Seconds())) / (p.runCount + 1)
			}

			p.runCount++

			log.Printf("Process '%s' has an average duration of %d seconds\n", p.fullPath, p.avgDuration)

			if p.avgDuration < 10 {
				alert.SendAlert(p.alertID, "Process '%s' keeps on crashing (average life %d seconds). Shutting down.", p.fullPath, p.avgDuration)
				break
			}
		}

	}(p)
}

func (p *internalProcess) stop() error {
	if p.cmd.Process == nil {
		return cerr.NewBadRequestError("Process '%s': Unable to stop a process that hasn't started yet", p.fullPath)
	}

	if !p.isRunning {
		return cerr.NewBadRequestError("Process '%s': Unable to stop a process that isn't running", p.fullPath)
	}

	// we set the terminate flag
	p.terminate = true

	p.isInterrupted = true

	// we send an ctrl+C signal
	return p.cmd.Process.Signal(os.Interrupt)

}

func (p *internalProcess) restart() error {
	if p.cmd.Process == nil {
		return cerr.NewBadRequestError("Process '%s': Unable to stop a process that hasn't started yet", p.fullPath)
	}

	if !p.isRunning {
		return cerr.NewBadRequestError("Process '%s': Unable to stop a process that isn't running", p.fullPath)
	}

	// we set the terminate flag to false
	p.terminate = false

	p.isInterrupted = true

	// we send an ctrl+C signal
	return p.cmd.Process.Signal(os.Interrupt)
}
