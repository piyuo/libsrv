package server

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/piyuo/libsrv/src/log"
)

// deadlineTask cache os env DEADLINE_TASK value
//
var deadlineTask time.Duration = -1

var lockHeader string

// setDeadlineTask set context deadline using os.Getenv("DEADLINE_TASK"), return CancelFunc that Canceling this context releases resources associated with it, so code should call cancel as soon as the operations running in this Context complete.
//
func setDeadlineTask(ctx context.Context) (context.Context, context.CancelFunc) {
	if deadlineTask == -1 {
		lockHeader = os.Getenv("NAME") + "-" + os.Getenv("BRANCH")
		text := os.Getenv("DEADLINE_TASK")
		ms, err := strconv.Atoi(text)
		if err != nil {
			ms = 840000
			log.Print(ctx, "task", "use default 840 seconds for DEADLINE_TASK")
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
			log.Error(ctx, here, errors.New("TaskID not found"))
			return
		}
		lockID := CreateLockID(taskID)

		db, err := New(ctx)
		if !found {
			log.Error(ctx, here, errors.New("failed to create task lock database connection"))
			return
		}
		defer db.Close()

		locked, err := db.LockTask(ctx, lockID, 15*time.Minute)
		if err != nil {
			log.Error(ctx, here, errors.New("failed to lock task"))
			return
		}
		if !locked {
			log.Print(ctx, here, "task already in progress. just wait")
			w.WriteHeader(http.StatusTooManyRequests) // return 429/503 will let google cloud slowing down execution
			return
		}
		defer db.DeleteTaskLock(ctx, lockID)

		retry, err := taskHandler(ctx, w, r)
		if err != nil {
			log.Error(ctx, here, err)
		}

		if retry {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
	return http.HandlerFunc(f)
}
