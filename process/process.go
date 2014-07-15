package process

import (
	"container/list"
	"prosit/cerr"
	"sync"
)

var l *list.List
var lMutex sync.RWMutex

func init() {
	l = list.New()
}

type Process struct {
	Pid       int      `json:"pid"`
	Run       string   `json:"run"`
	ArgList   []string `json:"argList"`
	Folder    string   `json:"folder"`
	Error     error    `json:"error"`
	Started   int64    `json:"started"`
	IsRunning bool     `json:"isRunning"`
	AlertID   string   `json:"alertID"`
}

func AddProcess(run, folder, alertID string) error {

	lMutex.Lock()
	defer lMutex.Unlock()

	// we create an internal process
	intProc, err := newInternalProcess(run, folder, alertID)

	if err != nil {
		return err
	}

	// we add to the list
	l.PushBack(intProc)

	// we start the Process
	intProc.start()

	return nil
}

func ListProcesses() []Process {

	lMutex.RLock()
	defer lMutex.RUnlock()

	ret := make([]Process, 0)

	// Iterate through list and print its contents.
	for e := l.Front(); e != nil; e = e.Next() {
		var tmpIntProcess = e.Value.(*internalProcess)

		var tmpProcess = &Process{}
		tmpProcess.Pid = tmpIntProcess.pid
		tmpProcess.Run = tmpIntProcess.cmd.Path
		tmpProcess.ArgList = tmpIntProcess.cmd.Args[1:]
		tmpProcess.Folder = tmpIntProcess.cmd.Dir
		tmpProcess.Error = tmpIntProcess.err
		tmpProcess.Started = tmpIntProcess.lastStarted
		tmpProcess.IsRunning = tmpIntProcess.isRunning
		tmpProcess.AlertID = tmpIntProcess.alertID

		ret = append(ret, *tmpProcess)
	}

	return ret
}

func GetProcessLogs(pid int) ([]string, error) {

	lMutex.RLock()
	defer lMutex.RUnlock()

	// Iterate through list and print its contents.
	for e := l.Front(); e != nil; e = e.Next() {
		var tmpIntProcess = e.Value.(*internalProcess)

		if tmpIntProcess.pid == pid {
			return tmpIntProcess.stdout.LogList(), nil
		}
	}

	return nil, cerr.NewBadRequestError("Process %d not found", pid)
}

func GetProcessErrors(pid int) ([]string, error) {

	lMutex.RLock()
	defer lMutex.RUnlock()

	// Iterate through list and print its contents.
	for e := l.Front(); e != nil; e = e.Next() {
		var tmpIntProcess = e.Value.(*internalProcess)

		if tmpIntProcess.pid == pid {
			return tmpIntProcess.stderr.LogList(), nil
		}
	}

	return nil, cerr.NewBadRequestError("Process %d not found", pid)
}
