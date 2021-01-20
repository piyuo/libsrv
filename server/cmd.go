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

// cmdDeadline cache os env cmdDeadline value
//
var cmdDeadline time.Duration = -1

// setCmdDeadline set context deadline using os.Getenv("cmdDeadline"), return CancelFunc that Canceling this context releases resources associated with it, so code should call cancel as soon as the operations running in this Context complete.
//
func setCmdDeadline(ctx context.Context) (context.Context, context.CancelFunc) {
	if cmdDeadline == -1 {
		text := os.Getenv("cmdDeadline")
		ms, err := strconv.Atoi(text)
		if err != nil {
			ms = 20000
			fmt.Print("use default 20 seconds for cmdDeadline")
		}
		cmdDeadline = time.Duration(ms) * time.Millisecond
	}
	expired := time.Now().Add(cmdDeadline)
	return context.WithDeadline(ctx, expired)
}

// cmdHandler create handler for command
//
func (s *Server) cmdHandler() http.Handler {
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
	ctx, cancel := setCmdDeadline(r.Context())
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
