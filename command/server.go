package command

import (
	"context"
	fmt "fmt"
	"io"
	"io/ioutil"
	"net/http"

	app "github.com/piyuo/go-libsrv/app"
	log "github.com/piyuo/go-libsrv/log"
)

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

//contextWithToken add token to context if token exist in cookies
func (s *Server) contextWithToken(r *http.Request) (context.Context, app.Token) {
	ctx := r.Context()
	if len(r.Cookies()) == 0 {
		return ctx, nil
	}
	token, err := app.TokenFromCookie(r)
	if err != nil { // it is fine with no token, just return the context
		return ctx, nil
	}
	//return new context with token
	return token.ToContext(ctx), token
}

// Serve entry for http request, filter empty and bad request and send correct one to dispatch
//
// enable cross origin access
func (s *Server) Serve(w http.ResponseWriter, r *http.Request) {
	ctx, token := s.contextWithToken(r)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		msg := "bad request. request is empty"
		s.writeText(w, msg)
		log.Info(ctx, msg)
		return
	}

	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.writeText(w, err.Error())
		log.Info(ctx, "bad request. "+err.Error())
		return
	}
	if len(bytes) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		msg := "bad request, need command in request"
		s.writeText(w, msg)
		log.Info(ctx, msg)
		return
	}

	s.dispatch = &Dispatch{
		Map: s.Map,
	}
	bytes, err = s.dispatch.Route(ctx, bytes)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		//if anything wrong just log error and send error id to client
		errID := log.Error(ctx, err)
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
