package process

import (
	"bytes"
	"container/list"
	"fmt"
	"prosit/alert"
	"sync"
)

type Consumer struct {
	l       *list.List
	buf     bytes.Buffer
	mx      sync.RWMutex
	alertID string
}

func (c *Consumer) Write(p []byte) (n int, err error) {

	c.mx.Lock()
	defer c.mx.Unlock()

	if c.l == nil {
		c.l = list.New()
	}

	for _, val := range p {
		if val == '\n' {
			// we save the buffer to the list
			c.l.PushFront(c.buf.String())

			if c.alertID != "" {
				// we have an alert to send
				alert.SendAlert(c.alertID, fmt.Sprintf("ERROR: %s", c.buf.String()))
			}

			// we truncate the buffer
			c.buf.Truncate(0)

		} else {

			// we simply add to the buffer
			c.buf.WriteByte(val)
		}
	}

	// we reduce the list to 10k items
	for c.l.Len() > 10000 {
		c.l.Remove(c.l.Back())
	}

	return len(p), nil
}

func (c *Consumer) SetAlertID(alertID string) {
	c.alertID = alertID
}

func (c *Consumer) LogList() []string {

	c.mx.RLock()
	defer c.mx.RUnlock()

	if c.l == nil {
		c.l = list.New()
	}

	ret := make([]string, 0)

	for e := c.l.Front(); e != nil; e = e.Next() {
		ret = append(ret, e.Value.(string))
	}

	return ret

}
