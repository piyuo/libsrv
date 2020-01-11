package command

import (
	fmt "fmt"
	"io"
	"io/ioutil"
	"net/http"

	libsrv "github.com/piyuo/go-libsrv"
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
	libsrv.Sys().Check()
	if s.Map == nil {
		msg := "server need Map for command pattern, try &Server{Map:yourMap}"
		libsrv.Sys().Emergency(msg)
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
	withoutArchive := http.HandlerFunc(s.Main)
	// support local server gzip compress
	// withArchive := ArchiveHandler(withoutArchive)
	return withoutArchive
}

// Main entry for http request, filter empty and bad request and send correct one to dispatch
//
// enable cross origin access
func (s *Server) Main(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		msg := "bad request. request is empty"
		s.writeText(w, msg)
		libsrv.Sys().Info(msg)
		return
	}

	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.writeText(w, err.Error())
		libsrv.Sys().Info("bad request. " + err.Error())
		return
	}
	if len(bytes) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		msg := "bad request, must include command in request"
		s.writeText(w, msg)
		libsrv.Sys().Info(msg)
		return
	}

	s.dispatch = &Dispatch{
		Map: s.Map,
	}
	bytes, err = s.dispatch.Route(bytes)
	if err == ErrCommandParsing {
		w.WriteHeader(http.StatusBadRequest)
		msg := "bad request, failed to parsing command. " + err.Error()
		s.writeText(w, msg)
		libsrv.Sys().Warning(msg)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		s.writeText(w, "internal server error, we already log this error and will be fixed ASAP.")
		libsrv.Sys().Error(err)
		return
	}

	if bytes != nil {
		s.writeBinary(w, bytes)
	}
}

func (s *Server) writeText(w http.ResponseWriter, text string) {
	w.Header().Set("Content-Type", "text/plain")
	io.WriteString(w, text)
}

func (s *Server) writeBinary(w http.ResponseWriter, bytes []byte) {
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(bytes)
}
