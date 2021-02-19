package server

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/piyuo/libsrv/src/env"
	"github.com/piyuo/libsrv/src/log"
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
	if s.Map == nil {
		WriteStatus(w, http.StatusBadRequest, "this server don't have map to accept command")
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")

	//add deadline to context
	ctx, cancel := setDeadline(r.Context())
	defer cancel()

	//add request to context
	ctx = env.SetRequest(ctx, r)

	if r.Body == nil {
		WriteStatus(w, http.StatusBadRequest, "no request")
		return
	}
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(ctx, here, err)
		WriteStatus(w, http.StatusBadRequest, "failed to read request")
		return
	}
	if len(bytes) == 0 {
		WriteStatus(w, http.StatusBadRequest, "bad request")
		return
	}

	bytes, err = s.dispatch.Route(ctx, bytes)
	if err != nil {
		handleRouteException(ctx, w, err)
		return
	}
	WriteBinary(w, bytes)
}
