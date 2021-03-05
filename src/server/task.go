package server

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/piyuo/libsrv/src/log"
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
			log.Print(ctx, "task", "use default 840 seconds for DEADLINE_TASK")
		}
		deadlineTask = time.Duration(ms) * time.Millisecond
	}
	expired := time.Now().Add(deadlineTask)
	return context.WithDeadline(ctx, expired)
}

// TaskCreateFunc create task handler function
//
func TaskCreateFunc(taskHandler TaskHandler) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		//add deadline to context
		ctx, cancel := setDeadlineTask(r.Context())
		defer cancel()

		retry, err := taskHandler(ctx, w, r)
		if err != nil {
			log.Error(ctx, here, err)
		}

		if retry {
			w.WriteHeader(http.StatusContinue)
		}
	}
	return http.HandlerFunc(f)
}
