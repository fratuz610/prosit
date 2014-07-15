package process

import (
	"container/list"
	"prosit/err"
	"sync"
)

var l *list.List
var lMutex sync.RWMutex

func init() {
	l = list.New()
}

type Process struct {
	Pid     int      `json:"pid"`
	Run     string   `json:"run"`
	ArgList []string `json:"argList"`
	Folder  string   `json:"folder"`
	Error   error    `json:"error"`
	Started int64    `json:"started"`
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
		tmpProcess.Args = tmpIntProcess.cmd.Args
		tmpProcess.Folder = tmpIntProcess.cmd.Dir
		tmpProcess.Error = tmpIntProcess.err
		tmpProcess.Started = tmpIntProcess.lastStarted
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

	return nil, err.NewBadRequestError("Process %d not found", pid)
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

	return nil, err.NewBadRequestError("Process %d not found", pid)
}
