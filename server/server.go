package server

import (
	"context"
	goerrors "errors"
	fmt "fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/piyuo/libsrv/command"
	"github.com/piyuo/libsrv/log"
)

const here = "server"

// HTTPHandler let you handle http request directly, return true if request is handled successfully, return false will result  404 bad request
//
type HTTPHandler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

// Server handle http request and call dispatch
//
//      server := &Server{
//		    Map: &TestMap{},
//	    }
//
type Server struct {
	dispatch *command.Dispatch

	// Map is command map
	//
	Map command.IMap

	// TaskHandler is for long running task
	//
	TaskHandler HTTPHandler
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
	if s.Map == nil && s.TaskHandler == nil {
		msg := " Map or TaskHandler is missing, try &Server{Map:yourMap,TaskHandler: yourHandler}"
		panic(msg)
	}
	port, _ := s.prepare()
	log.Debug(context.Background(), here, fmt.Sprintf("start listening from http://localhost%v\n", port))
	http.ListenAndServe(port, nil)
}

// prepare server variable and return listening port like :8080
//
func (s *Server) prepare() (string, http.Handler) {

	rand.Seed(time.Now().UTC().UnixNano())

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	var rootHandler http.Handler
	if s.Map != nil {
		rootHandler = s.cmdHandler()
		http.Handle("/", rootHandler)
		http.Handle("/task", s.taskHandler())
	} else {
		rootHandler = s.taskHandler()
		http.Handle("/", rootHandler)
	}

	// set default map
	s.dispatch = &command.Dispatch{
		Map: s.Map,
	}
	s.Map = nil

	// set default task handler
	if s.TaskHandler == nil {
		s.TaskHandler = defaultTaskHandler
	}
	return ":" + port, rootHandler
}

func defaultTaskHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	return nil
}

// handleRouteException convert error to status code, so client command service know how to deal with it
//
func handleRouteException(ctx context.Context, w http.ResponseWriter, err error) {
	log.Debug(ctx, here, "[logged] "+err.Error())

	if goerrors.Is(err, context.DeadlineExceeded) {
		errID := log.Error(ctx, here, err)
		writeError(w, err, http.StatusGatewayTimeout, errID)
		return
	}

	errID := log.Error(ctx, here, err)
	writeError(w, err, http.StatusInternalServerError, errID)
}

// Query return value from query string
//
//	value, ok := Query(r, "type")
//	assert.True(ok)
//	assert.Equal("maintenance", value)
//
func Query(r *http.Request, param string) (string, bool) {
	params, ok := r.URL.Query()[param]
	if !ok || len(params[0]) < 1 {
		return "", false
	}
	return params[0], true
}
