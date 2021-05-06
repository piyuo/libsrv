package gtask

import (
	"time"

	"github.com/piyuo/libsrv/db"
	"github.com/pkg/errors"
)

// Task keep task lock records, id will be task id
//
type Task struct {
	db.Entity

	// Name is task name
	//
	Name string `firestore:"Name,omitempty"`

	// Duration is task lock duration in seconds. during lock no another task can be lock again
	//
	Duration int `firestore:"Duration,omitempty"`

	// LockTime is when this task is locked
	//
	LockTime time.Time `firestore:"LockTime,omitempty"`

	// MaxRetry is max retry count for this task
	//
	MaxRetry int `firestore:"MaxRetry,omitempty"`

	// Retry is total execute count
	//
	Retry int `firestore:"Retry,omitempty"`
}

// Factory create a empty object, return object must be nil safe, no nil in any field
//

func (c *Task) Factory() db.Object {
	return &Task{}
}

// Collection return the name in database
//
func (c *Task) Collection() string {
	return "Task"
}

// Lock return true if task can be locked and increment Execute count
//
func (c *Task) Lock() error {
	if c.Retry == c.MaxRetry {
		return errors.New("max retry exceeded")
	}

	deadline := time.Now().UTC().Add(-time.Duration(c.Duration) * time.Second)
	if c.LockTime.After(deadline) {
		return errors.New("task is running, the lock time is " + c.LockTime.Format("2006-01-02 15:04:05") + ", deadline is " + deadline.Format("2006-01-02 15:04:05"))
	}

	c.LockTime = time.Now().UTC()
	c.Retry++
	return nil
}

// newTask create new task
//
func newTask(id string, name string, duration, maxRetry int) *Task {
	task := &Task{
		Name:     name,
		MaxRetry: maxRetry,
		Duration: duration,
		Retry:    0,
	}
	task.SetID(id)
	return task
}
