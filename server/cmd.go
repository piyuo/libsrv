package server

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/piyuo/libsrv/env"
)

// deadline cache os env deadline value
//
var deadline time.Duration = -1

// setDeadline set context deadline using os.Getenv("DEADLINE"), return CancelFunc that Canceling this context releases resources associated with it, so code should call cancel as soon as the operations running in this Context complete.
//
func setDeadline(ctx context.Context) (context.Context, context.CancelFunc) {
	if deadline == -1 {
		text := os.Getenv("DEADLINE")
		ms, err := strconv.Atoi(text)
		if err != nil {
			ms = 20000
			fmt.Println("use default 20 seconds for deadline")
		}
		deadline = time.Duration(ms) * time.Millisecond
	}
	expired := time.Now().Add(deadline)
	return context.WithDeadline(ctx, expired)
}

// createCmdHandler create command handler
//
func (s *Server) createCmdHandler() http.Handler {
	withoutArchive := http.HandlerFunc(s.cmdServe)
	withArchive := ArchiveHandler(withoutArchive)
	return withArchive
}

// cmdServe serve command request, it filter empty and bad request and send correct one to dispatch
//
//	Cross origin access enabled
//
func (s *Server) cmdServe(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	//add deadline to context
	ctx, cancel := setDeadline(r.Context())
	defer cancel()

	//add request to context
	ctx = env.SetRequest(ctx, r)

	if r.Body == nil {
		writeBadRequest(ctx, w, "request has no body")
		return
	}
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeBadRequest(ctx, w, "bad request, "+err.Error())
		return
	}
	if len(bytes) == 0 {
		writeBadRequest(ctx, w, "empty request")
		return
	}

	bytes, err = s.dispatch.Route(ctx, bytes)
	if err != nil {
		handleRouteException(ctx, w, err)
		return
	}
	writeBinary(w, bytes)
}
