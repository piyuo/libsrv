package command

import (
	"context"
	fmt "fmt"
	"io"
	"io/ioutil"
	"net/http"

	app "github.com/piyuo/go-libsrv/app"
	log "github.com/piyuo/go-libsrv/log"
	"github.com/pkg/errors"
)

const here = "command"

// Server handle http request and call dispatch
//
//      server := &Server{
//		    Map: &TestMap{},
//	    }
type Server struct {
	dispatch *Dispatch
	Map      IMap
}

// Start http server to listen request and serve content
//
//      var server = &command.Server{
//      Map: &commands.MapXXX{},
//     }
//     func main() {
//      server.Start(80)
//     }
func (s *Server) Start(port int) {
	app.Check()

	if s.Map == nil {
		msg := "server need Map for command pattern, try &Server{Map:yourMap}"
		panic(msg)
	}
	http.Handle("/", s.newHandler())
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

// create handler with archive ability
//
//      server := &Server{
//		    Map: &TestMap{},
//	    }
//  return server.newHandler()
func (s *Server) newHandler() http.Handler {
	withoutArchive := http.HandlerFunc(s.Serve)
	// support local server gzip compress
	// withArchive := ArchiveHandler(withoutArchive)
	return withoutArchive
}

// Serve entry for http request, filter empty and bad request and send correct one to dispatch
//
//cross origin access enabled
//
func (s *Server) Serve(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	ctx, cancel, token, err := contextWithTokenAndDeadline(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		msg := "failed to set context deadline, may need app.Check()"
		s.writeText(w, msg)
		log.Debug(ctx, here, msg)
		return
	}
	defer cancel()

	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		msg := "bad request. request is empty"
		s.writeText(w, msg)
		log.Debug(ctx, here, msg)
		return
	}

	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.writeText(w, err.Error())
		log.Debug(ctx, here, "bad request. "+err.Error())
		return
	}
	if len(bytes) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		msg := "bad request, need command in request"
		s.writeText(w, msg)
		log.Debug(ctx, here, msg)
		return
	}

	s.dispatch = &Dispatch{
		Map: s.Map,
	}
	bytes, err = s.dispatch.Route(ctx, bytes)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		//if anything wrong just log error and send error id to client
		errID := log.Error(ctx, here, err, r)
		s.writeText(w, errID)
		return
	}

	//check to see if token need revive
	if token != nil && token.Revive() {
		token.ToCookie(w)
	}
	s.writeBinary(w, bytes)
}

func (s *Server) writeText(w http.ResponseWriter, text string) {
	w.Header().Set("Content-Type", "text/plain")
	io.WriteString(w, text)
}

func (s *Server) writeBinary(w http.ResponseWriter, bytes []byte) {
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(bytes)
}

//contextWithTokenAndDeadline add token to context if token exist in cookies and deadline
//
//	context,cancel, token, err := contextWithTokenAndDeadline(req)
func contextWithTokenAndDeadline(r *http.Request) (context.Context, context.CancelFunc, app.Token, error) {
	dateline, err := app.ContextDateline()
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "failed to get context deadline")
	}

	ctx, cancel := context.WithDeadline(r.Context(), dateline)
	if len(r.Cookies()) == 0 {
		return ctx, cancel, nil, nil
	}
	token, err := app.TokenFromCookie(r)
	if err != nil { // it is fine with no token, just return the context
		return ctx, cancel, nil, nil
	}
	//return new context with token
	return token.ToContext(ctx), cancel, token, nil
}
