package libsrv

import (
	"time"
)

//Timer interface
type Timer interface {
	Start()
	Stop() int64
}

type timer struct {
	current time.Time
}

//NewTimer create a timer
func NewTimer() Timer {
	return &timer{}
}

func (t *timer) Start() {
	t.current = time.Now()
}

func (t *timer) Stop() int64 {
	duration := time.Now().Sub(t.current)
	ns := duration.Nanoseconds()
	ms := ns / 1000000
	return ms
}
