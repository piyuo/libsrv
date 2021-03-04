package server

import (
	"context"
	goerrors "errors"
	fmt "fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/piyuo/libsrv/src/command"
	"github.com/piyuo/libsrv/src/log"
)

const here = "server"

// HTTPHandler let you handle http request directly, return true if request is handled successfully, return false will result  404 bad request
//
type HTTPHandler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

// Server handle http request and call dispatch
//
//	server := &Server{
//		CommandHandlers: map[string]command.IMap{"/": &mock.MapXXX{}},
//		HTTPHandlers:    map[string]HTTPHandler{"/api": httpHandler},
//	}
//
type Server struct {
	//	dispatch *command.Dispatch

	// CommandHandlers is command map for handle command request
	//
	CommandHandlers map[string]command.IMap

	// HTTPHandlers is http handler map to handle http request
	//
	HTTPHandlers map[string]HTTPHandler
}

// Start http server to listen request and serve content, defult port is 8080, you can change use export PORT="8080"
//
//	server := &Server{
//		CommandHandlers: map[string]command.IMap{"/": &mock.MapXXX{}},
//		HTTPHandlers:    map[string]HTTPHandler{"/api": httpHandler},
//	}
//  func main() {
//      server.Start()
//  }
//
func (s *Server) Start() {
	ctx := context.Background()
	if s.CommandHandlers == nil && s.HTTPHandlers == nil {
		msg := "CommandHandlers or HTTPHandlers is missing, try add &Server{CommandHandlers:yourCommandHandler, HTTPHandlers: yourHttpHandler}"
		panic(msg)
	}

	if err := http.ListenAndServe(s.ready(ctx), nil); err != nil {
		log.Error(ctx, here, err)
	}
}

// ready server variable and return listening port like :8080
//
func (s *Server) ready(ctx context.Context) string {
	rand.Seed(time.Now().UTC().UnixNano())

	if s.CommandHandlers != nil {
		for pattern, cmdMap := range s.CommandHandlers {
			http.Handle(pattern, CMDCreateFunc(cmdMap))
			break // only allow one command handler for now
		}
		//realease command handlers, don't need anymore
		s.CommandHandlers = nil
	}

	if s.HTTPHandlers != nil {
		for pattern, httpHandler := range s.HTTPHandlers {
			http.Handle(pattern, HTTPCreateFunc(httpHandler))
		}
		//realease http handlers, don't need anymore
		s.HTTPHandlers = nil
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Print(ctx, here, fmt.Sprintf("start listening from http://localhost%v\n", port))
	return ":" + port
}

// handleRouteException convert error to status code, so client command service know how to deal with it
//
func handleRouteException(ctx context.Context, w http.ResponseWriter, err error) {
	log.Print(ctx, here, "[logged] "+err.Error())

	if goerrors.Is(err, context.DeadlineExceeded) {
		WriteStatus(w, http.StatusGatewayTimeout, "Deadline Exceeded")
		return
	}

	log.Error(ctx, here, err)
	WriteError(w, http.StatusInternalServerError, err)
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
