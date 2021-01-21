package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

// apiDeadline cache os env apiDeadline value
//
var apiDeadline time.Duration = -1

// setAPIDeadline set context deadline using os.Getenv("API_DEADLINE"), return CancelFunc that Canceling this context releases resources associated with it, so code should call cancel as soon as the operations running in this Context complete.
//
func setAPIDeadline(ctx context.Context) (context.Context, context.CancelFunc) {
	if apiDeadline == -1 {
		text := os.Getenv("API_DEADLINE")
		ms, err := strconv.Atoi(text)
		if err != nil {
			ms = 20000
			fmt.Println("use default 20 seconds for apiDeadline")
		}
		apiDeadline = time.Duration(ms) * time.Millisecond
	}
	expired := time.Now().Add(apiDeadline)
	return context.WithDeadline(ctx, expired)
}

// createAPIHandler create general handler
//
func (s *Server) createAPIHandler() http.Handler {
	return http.HandlerFunc(s.apiServe)
}

// apiServe serve api request
//
func (s *Server) apiServe(w http.ResponseWriter, r *http.Request) {

	//add deadline to context
	ctx, cancel := setAPIDeadline(r.Context())
	defer cancel()

	err := s.APIHandler(ctx, w, r)
	if err != nil {
		handleRouteException(ctx, w, err)
		return
	}
}
