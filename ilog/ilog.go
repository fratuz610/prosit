package ilog

import (
	"bytes"
	"container/list"
	"io"
	//"log"
	"os"
	"sync"
)

var l *list.List
var mx sync.RWMutex
var writer *os.File

func init() {
	l = list.New()
}

func RedirectOutput() {

	var r *os.File
	r, writer, _ = os.Pipe()
	os.Stdout = writer
	os.Stderr = writer

	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {

		var lineBuf bytes.Buffer
		b := make([]byte, 1024)
		for {
			read, _ := r.Read(b)

			for _, val := range b[0:read] {
				if val == '\n' {
					mx.Lock()
					//log.Printf("Got string '%s'\n", lineBuf.String())
					l.PushFront(lineBuf.String())
					mx.Unlock()

					lineBuf.Truncate(0)

					// we reduce the list to 10k items
					for l.Len() > 10000 {
						l.Remove(l.Back())
					}
				} else {

					lineBuf.WriteByte(val)
				}
			}
		}

	}()

}

func GetOutput() []string {

	mx.RLock()
	defer mx.RUnlock()

	if l == nil {
		l = list.New()
	}

	ret := make([]string, 0)

	for e := l.Front(); e != nil; e = e.Next() {
		ret = append(ret, e.Value.(string))
	}

	return ret
}

func GetWriter() io.Writer {
	return writer
}
