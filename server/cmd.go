package server

import (
	"io/ioutil"
	"net/http"

	"github.com/piyuo/libsrv/env"
)

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
