package server

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/piyuo/libsrv/gtask"
	"github.com/pkg/errors"

	"github.com/piyuo/libsrv/log"
)

// deadlineTask cache os env DEADLINE_TASK value
//
var deadlineTask time.Duration = -1

// setDeadlineTask set context deadline using os.Getenv("DEADLINE_TASK"), return CancelFunc that Canceling this context releases resources associated with it, so code should call cancel as soon as the operations running in this Context complete.
//
func setDeadlineTask(ctx context.Context) (context.Context, context.CancelFunc) {
	if deadlineTask == -1 {
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

// TaskEntry create task entry function
//
func TaskEntry(taskHandler TaskHandler) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		//add deadline to context
		ctx, cancel := setDeadlineTask(r.Context())
		defer cancel()

		err := TaskRun(ctx, taskHandler, r)
		if err != nil {
			log.Error(ctx, err)
			WriteError(w, http.StatusOK, err) // return OK to stop retry
			return
		}
	}
	return http.HandlerFunc(f)
}

func TaskRun(ctx context.Context, taskHandler TaskHandler, r *http.Request) error {
	_, found := Query(r, "debug")
	taskID := ""
	if !found {
		// no need to lock task when debug
		taskID, found = Query(r, "TaskID")
		if !found {
			return errors.New("TaskID not found")
		}

		err := gtask.Lock(ctx, taskID)
		if err != nil {
			return errors.Wrap(err, "lock task "+taskID)
		}
		defer gtask.Delete(ctx, taskID)
	}

	log.Info(ctx, "start task "+taskID)
	err := taskHandler(ctx, r)
	if err != nil {
		return errors.Wrap(err, "task handler fail")
	}
	log.Info(ctx, "finish task "+taskID)
	return nil
}
