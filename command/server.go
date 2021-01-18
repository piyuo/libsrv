package command

import (
	"context"
	goerrors "errors"
	fmt "fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/piyuo/libsrv/log"
	"github.com/piyuo/libsrv/session"
)

const here = "command"

// CustomHTTPHandler let you handle http request directly, return true if request is handled, return false let command pattern do the job
//
type CustomHTTPHandler func(ctx context.Context, w http.ResponseWriter, r *http.Request) (bool, error)

// Server handle http request and call dispatch
//
//      server := &Server{
//		    Map: &TestMap{},
//	    }
//
type Server struct {
	dispatch    *Dispatch
	Map         IMap
	HTTPHandler CustomHTTPHandler
}

// Start http server to listen request and serve content, defult port is 8080, you can change use export PORT="8080"
//
//	var server = &command.Server{
//  	Map: &commands.MapXXX{},
//  }
//  func main() {
//      server.Start()
//  }
//
func (s *Server) Start() {
	rand.Seed(time.Now().UTC().UnixNano())

	portText := os.Getenv("PORT")
	if portText == "" {
		portText = "8080"
	}
	port, err := strconv.Atoi(portText)
	if err != nil {
		fmt.Printf("please set http listen port like export PORT=\"8080\"")
		return
	}
	http.Handle("/", s.newHandler())
	if s.Map == nil && s.HTTPHandler == nil {
		msg := "server need Map for command pattern or custom HTTPHandler , try &Server{Map:yourMap,HTTPHandler: customHTTPHandler}"
		panic(msg)
	}
	s.dispatch = &Dispatch{
		Map: s.Map,
	}
	s.Map = nil
	ctx := context.Background()
	log.Debug(ctx, here, fmt.Sprintf("start listening from http://localhost:%d\n", port))
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

// create handler with archive ability
//
//      server := &Server{
//		    Map: &TestMap{},
//	    }
//  return server.newHandler()
//
func (s *Server) newHandler() http.Handler {
	withoutArchive := http.HandlerFunc(s.Serve)
	// support gzip compress
	withArchive := ArchiveHandler(withoutArchive)
	return withArchive
}

// Serve entry for http request, filter empty and bad request and send correct one to dispatch
//
//cross origin access enabled
//
func (s *Server) Serve(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	//add deadline to context
	ctx, cancel := session.SetDeadline(r.Context())
	defer cancel()

	//add request to context
	ctx = session.SetRequest(ctx, r)

	// handle by custom http handler ?
	if s.HTTPHandler != nil {
		result, err := s.HTTPHandler(ctx, w, r)
		if err != nil {
			handleRouteException(ctx, w, err)
			return
		}
		if result == true {
			return
		}
		// custom handler return false mean request still need go through command dispatch
	}

	if r.Body == nil {
		writeBadRequest(ctx, w, "empty request")
		return
	}
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeBadRequest(ctx, w, "bad request, "+err.Error())
		return
	}
	if len(bytes) == 0 {
		writeBadRequest(ctx, w, "bad request")
		return
	}

	bytes, err = s.dispatch.Route(ctx, bytes)
	if err != nil {
		handleRouteException(ctx, w, err)
		return
	}

	writeBinary(w, bytes)
}

//handleRouteException convert error to status code, so client command service know how to deal with it
//
//
func handleRouteException(ctx context.Context, w http.ResponseWriter, err error) {
	log.Debug(ctx, here, "[solved] "+err.Error())

	if goerrors.Is(err, context.DeadlineExceeded) {
		errID := log.Error(ctx, here, err)
		writeError(w, err, http.StatusGatewayTimeout, errID)
		return
	}

	errID := log.Error(ctx, here, err)
	writeError(w, err, http.StatusInternalServerError, errID)
}
