package command

import (
	fmt "fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
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
	if s.Map == nil {
		panic("Server need Map for command pattern, try &Server{Map:yourMap}")
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
	//todo:skip check url for now
	//if len(r.URL.Path) != 1 {
	//	w.WriteHeader(http.StatusNotFound)
	//	fmt.Println(r.URL.Path + " not found")
	//	return
	//}

	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.writeText(w, err.Error())
		return
	}
	if len(bytes) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		s.writeText(w, "command not exist")
		return
	}

	s.dispatch = &Dispatch{
		Map: s.Map,
	}
	bytes, err = s.dispatch.Route(bytes)
	if err == ErrCommandParsing {
		w.WriteHeader(http.StatusBadRequest)
		s.writeText(w, err.Error())
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		s.writeText(w, err.Error())
		s.handleException(err)
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

func (s *Server) handleException(err error) {
	log.Fatalf("dispatch error : %v", err)

}
