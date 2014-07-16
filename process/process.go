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
	Id        string   `json:"id"`
	Pid       int      `json:"pid"`
	Run       string   `json:"run"`
	ArgList   []string `json:"argList"`
	Folder    string   `json:"folder"`
	Error     error    `json:"error"`
	Started   int64    `json:"started"`
	IsRunning bool     `json:"isRunning"`
	AlertID   string   `json:"alertID"`
}

func AddProcess(id, run, folder, alertID string) error {

	lMutex.Lock()
	defer lMutex.Unlock()

	// we create an internal process
	intProc, err := newInternalProcess(id, run, folder, alertID)

	if err != nil {
		return err
	}

	// we add to the list
	l.PushBack(intProc)

	// we start the Process
	intProc.start()

	return nil
}

func ProcessExists(id string) bool {

	lMutex.RLock()
	defer lMutex.RUnlock()

	// Iterate through list and print its contents.
	for e := l.Front(); e != nil; e = e.Next() {
		var tmpProc = e.Value.(*internalProcess)

		if tmpProc.id == id {
			return true
		}
	}

	return false
}

func ListProcesses() []Process {

	lMutex.RLock()
	defer lMutex.RUnlock()

	ret := make([]Process, 0)

	// Iterate through list and print its contents.
	for e := l.Front(); e != nil; e = e.Next() {
		var tmpIntProcess = e.Value.(*internalProcess)

		var tmpProcess = &Process{}
		tmpProcess.Id = tmpIntProcess.id
		tmpProcess.Pid = tmpIntProcess.pid
		tmpProcess.Run = tmpIntProcess.cmd.Path
		tmpProcess.ArgList = tmpIntProcess.cmd.Args[1:]
		tmpProcess.Folder = tmpIntProcess.folder
		tmpProcess.Error = tmpIntProcess.err
		tmpProcess.Started = tmpIntProcess.lastStarted
		tmpProcess.IsRunning = tmpIntProcess.isRunning
		tmpProcess.AlertID = tmpIntProcess.alertID

		ret = append(ret, *tmpProcess)
	}

	return ret
}

func GetProcessLogs(id string) ([]string, error) {

	lMutex.RLock()
	defer lMutex.RUnlock()

	// Iterate through list and print its contents.
	for e := l.Front(); e != nil; e = e.Next() {
		var tmpIntProcess = e.Value.(*internalProcess)

		if tmpIntProcess.id == id {
			return tmpIntProcess.stdout.LogList(), nil
		}
	}

	return nil, cerr.NewBadRequestError("Process '%s' not found", id)
}

func GetProcessErrors(id string) ([]string, error) {

	lMutex.RLock()
	defer lMutex.RUnlock()

	// Iterate through list and print its contents.
	for e := l.Front(); e != nil; e = e.Next() {
		var tmpIntProcess = e.Value.(*internalProcess)

		if tmpIntProcess.id == id {
			return tmpIntProcess.stderr.LogList(), nil
		}
	}

	return nil, cerr.NewBadRequestError("Process '%s' not found", id)
}
