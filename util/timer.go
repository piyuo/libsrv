package util

import (
	"time"
)

// Timer implement timer function
type Timer interface {
	// Start timer
	//
	//	timer := NewTimer()
	//	timer.Start()
	//
	Start()

	// Stop timer
	//
	//	timer := NewTimer()
	//	timer.Start()
	//
	Stop() int64

	// TimeSpan get current time span and not stop timer
	//
	//	timer := NewTimer()
	//	timer.Start()
	//	timeSpan := timer.Current()
	//
	TimeSpan() int64
}

//NewTimer create a timer
//
//	timer := NewTimer()
//
func NewTimer() Timer {
	return &timer{}
}

type timer struct {
	current time.Time
}

// Start timer
//
//	timer := NewTimer()
//	timer.Start()
//
func (t *timer) Start() {
	t.current = time.Now()
}

// Stop timer
//
//	timer := NewTimer()
//	timer.Start()
//
func (t *timer) Stop() int64 {
	return t.TimeSpan()
}

// TimeSpan get current time span and not stop timer
//
//	timer := NewTimer()
//	timer.Start()
//	ms := timer.TimeSpan()
//
func (t *timer) TimeSpan() int64 {
	duration := time.Since(t.current)
	ns := duration.Nanoseconds()
	ms := ns / 1000000
	return ms
}
