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

	// Map is command map for command handler
	//
	Map command.IMap

	// APIHandler is API handler
	//
	APIHandler HTTPHandler
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
	name := os.Getenv("NAME")
	if name == "" {
		panic("missing environment variable: NAME, try set NAME=\"serviceName\"")
	}

	region := os.Getenv("REGION")
	if region == "" {
		panic("missing environment variable: REGION, try set REGION=\"US\"")
	}

	branch := os.Getenv("BRANCH")
	if branch == "" {
		panic("missing environment variable: BRANCH, try set BRANCH=\"master\"")
	}

	if s.Map == nil && s.APIHandler == nil {
		msg := " Map or APIHandler is missing, try &Server{Map:yourMap, ApiHandler: yourHandler}"
		panic(msg)
	}
	port := s.prepare()
	log.Debug(context.Background(), here, fmt.Sprintf("start listening from http://localhost%v\n", port))

	if err := http.ListenAndServe(port, nil); err != nil {
		fmt.Printf("Server failed: %s\n", err)
	}
}

// prepare server variable and return listening port like :8080
//
func (s *Server) prepare() string {

	rand.Seed(time.Now().UTC().UnixNano())

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.Handle("/", s.createCmdHandler())
	http.Handle("/api", s.createAPIHandler())

	// set default map
	s.dispatch = &command.Dispatch{
		Map: s.Map,
	}
	s.Map = nil

	// set default task handler
	if s.APIHandler == nil {
		s.APIHandler = defaultHandler
	}

	return ":" + port
}

func defaultHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	writeStatus(w, http.StatusForbidden, "Forbidden")
	return nil
}

// handleRouteException convert error to status code, so client command service know how to deal with it
//
func handleRouteException(ctx context.Context, w http.ResponseWriter, err error) {
	log.Debug(ctx, here, "[logged] "+err.Error())

	if goerrors.Is(err, context.DeadlineExceeded) {
		writeStatus(w, http.StatusGatewayTimeout, "Deadline Exceeded")
		return
	}

	errID := log.Error(ctx, here, err)
	writeError(w, http.StatusInternalServerError, errID, err)
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
