package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

// taskDeadline cache os env taskDeadline value
//
var taskDeadline time.Duration = -1

// setTaskDeadline set context deadline using os.Getenv("taskDeadline"), return CancelFunc that Canceling this context releases resources associated with it, so code should call cancel as soon as the operations running in this Context complete.
//
func setTaskDeadline(ctx context.Context) (context.Context, context.CancelFunc) {
	if taskDeadline == -1 {
		text := os.Getenv("taskDeadline")
		ms, err := strconv.Atoi(text)
		if err != nil {
			ms = 20000
			fmt.Print("use default 20 seconds for taskDeadline")
		}
		taskDeadline = time.Duration(ms) * time.Millisecond
	}
	expired := time.Now().Add(taskDeadline)
	return context.WithDeadline(ctx, expired)
}

// taskHandler create handler for task
//
func (s *Server) taskHandler() http.Handler {
	return http.HandlerFunc(s.taskServe)
}

// taskServe serve task request
//
func (s *Server) taskServe(w http.ResponseWriter, r *http.Request) {

	//add deadline to context
	ctx, cancel := setTaskDeadline(r.Context())
	defer cancel()

	err := s.TaskHandler(ctx, w, r)
	if err != nil {
		handleRouteException(ctx, w, err)
		return
	}
}
