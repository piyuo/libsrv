package server

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/piyuo/libsrv/src/db"
	"github.com/piyuo/libsrv/src/google/gaccount"
	"github.com/piyuo/libsrv/src/google/gdb"
	"github.com/pkg/errors"

	"github.com/piyuo/libsrv/src/log"
)

// deadlineTask cache os env DEADLINE_TASK value
//
var deadlineTask time.Duration = -1

var lockHeader string

// TaskLock keep task lock records
//
type TaskLock struct {
	db.BaseObject
}

func (c *TaskLock) Factory() db.Object {
	return &TaskLock{}
}

func (c *TaskLock) Collection() string {
	return "TaskLock"
}

// setDeadlineTask set context deadline using os.Getenv("DEADLINE_TASK"), return CancelFunc that Canceling this context releases resources associated with it, so code should call cancel as soon as the operations running in this Context complete.
//
func setDeadlineTask(ctx context.Context) (context.Context, context.CancelFunc) {
	if deadlineTask == -1 {
		lockHeader = os.Getenv("NAME") + "-" + os.Getenv("BRANCH")
		text := os.Getenv("DEADLINE_TASK")
		ms, err := strconv.Atoi(text)
		if err != nil {
			ms = 840000
			log.Warn(ctx, "use default 840 seconds for DEADLINE_TASK")
		}
		deadlineTask = time.Duration(ms) * time.Millisecond
	}
	expired := time.Now().Add(deadlineTask)
	return context.WithDeadline(ctx, expired)
}

// CreateLockID from taskID
//
func CreateLockID(taskID string) string {
	return lockHeader + "-" + taskID
}

// TaskCreateFunc create task handler function
//
func TaskCreateFunc(taskHandler TaskHandler) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		//add deadline to context
		ctx, cancel := setDeadlineTask(r.Context())
		defer cancel()

		taskID, found := Query(r, "TaskID")
		if !found {
			log.Error(ctx, errors.New("TaskID not found"))
			return
		}
		lockID := CreateLockID(taskID)

		locked, err := lockTask(ctx, lockID, 15*time.Minute)
		if err != nil {
			log.Error(ctx, err)
			return
		}
		if !locked {
			log.Info(ctx, "task in progress. just wait")
			w.WriteHeader(http.StatusTooManyRequests) // return 429/503 will let google cloud slowing down execution
			return
		}
		defer unlockTask(ctx, lockID)
		retry, err := taskHandler(ctx, w, r)
		if err != nil {
			log.Error(ctx, err)
		}

		if retry {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
	return http.HandlerFunc(f)
}

// tasklockClient return tasklockClient
//
//	client,err := tasklockClient(ctx)
//	defer client.Close()
//
func newClient(ctx context.Context) (db.Client, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	cred, err := gaccount.GlobalCredential(ctx)
	if err != nil {
		return nil, err
	}
	client, err := gdb.NewClient(ctx, cred)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// createTaskLock create lock on task
//
//	err = createTaskLock(ctx, client, "lock id")
//
func createTaskLock(ctx context.Context, client db.Client, lockID string) error {
	lock := &TaskLock{}
	lock.SetID(lockID)
	return client.Set(ctx, lock)
}

// isTaskLockExists check lock exists
//
//	found, createTime, err := isTaskLockExists(ctx, client, "lock id")
//
func isTaskLockExists(ctx context.Context, client db.Client, lockID string) (bool, time.Time, error) {
	obj, err := client.Get(ctx, &TaskLock{}, lockID)
	if err != nil {
		return false, time.Time{}, errors.Wrap(err, "get task lock:"+lockID)
	}
	if obj != nil {
		return true, obj.CreateTime(), nil
	}
	return false, time.Time{}, nil
}

// lockTask lock task for 15 mins
//
//	locked, err := lockTask(ctx, "lock id")
//
func lockTask(ctx context.Context, lockID string, duration time.Duration) (bool, error) {
	client, err := newClient(ctx)
	if err != nil {
		return false, err
	}
	defer client.Close()

	found, createTime, err := isTaskLockExists(ctx, client, lockID)
	if err != nil {
		return false, errors.Wrap(err, "check lock exists:"+lockID)
	}
	if !found {
		if err := createTaskLock(ctx, client, lockID); err != nil {
			return false, errors.Wrap(err, "create task lock")
		}
		return true, nil
	}

	deadline := time.Now().UTC().Add(-duration)
	if createTime.Before(deadline) {
		// this target is too old, maybe something went wrong
		lock := &TaskLock{}
		lock.SetID(lockID)
		if err := client.Update(ctx, lock, map[string]interface{}{"CreateTime": time.Now().UTC()}); err != nil {
			return false, errors.Wrap(err, "add task lock")
		}
		return true, nil
	}
	return false, nil
}

// unlockTask unlock task
//
//	locked, err := unlockTask(ctx, "lock id")
//
func unlockTask(ctx context.Context, lockID string) error {
	client, err := newClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	lock := &TaskLock{}
	lock.SetID(lockID)
	return client.Delete(ctx, lock)
}
