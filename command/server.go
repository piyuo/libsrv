package command

import (
	"context"
	goerrors "errors"
	fmt "fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	app "github.com/piyuo/libsrv/app"
	log "github.com/piyuo/libsrv/log"
)

const here = "command"

//CustomHTTPHandler let you handle http request directly, return true if request is handled, return false let command pattern do the job
type CustomHTTPHandler func(w http.ResponseWriter, r *http.Request) (bool, error)

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
//      var server = &command.Server{
//      Map: &commands.MapXXX{},
//     }
//     func main() {
//      server.Start()
//     }
//
func (s *Server) Start() {
	rand.Seed(time.Now().UTC().UnixNano())
	app.Check()
	if s.Map == nil {
		msg := "server need Map for command pattern, try &Server{Map:yourMap}"
		panic(msg)
	}
	portText := os.Getenv("PORT")
	if portText == "" {
		portText = "8080"
	}
	port, err := strconv.Atoi(portText)
	if err != nil {
		fmt.Printf("please set http listen port like export PORT=\"8080\"")
		return
	}
	fmt.Printf("start listening from http://localhost:%d\n", port)
	http.Handle("/", s.newHandler())
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
	ctx, cancel, token, err := contextWithTokenAndDeadline(r)
	if err != nil {
		handleEnvNotReady(ctx, w)
		return
	}
	defer cancel()

	if s.HTTPHandler != nil {
		result, err := s.HTTPHandler(w, r)
		if err != nil {
			handleRouteException(ctx, w, r, err)
			return
		}
		if result == true {
			return
		}
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

	s.dispatch = &Dispatch{
		Map: s.Map,
	}
	bytes, err = s.dispatch.Route(ctx, bytes)
	if err != nil {
		handleRouteException(ctx, w, r, err)
		return
	}

	//check to see if token need revive
	if token != nil && token.Revive() {
		token.ToCookie(w)
	}
	writeBinary(w, bytes)
}

// handleEnvNotReady write error to response, the server environment variable is not set
//
//	logEnvMissing(context.Background(), w)
//
func handleEnvNotReady(ctx context.Context, w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotImplemented)
	msg := "failed to set context deadline, may be missing environment variable, use app.Check() to make sure all var are set"
	writeText(w, msg)
	log.Debug(ctx, here, msg)
}

//handleRouteException convert error to status code, so client command service know how to deal with it
//
//	context,cancel, token, err := contextWithTokenAndDeadline(req)
//
func handleRouteException(ctx context.Context, w http.ResponseWriter, r *http.Request, err error) {
	log.Debug(ctx, here, "[solved] "+err.Error())
	if goerrors.Is(err, context.DeadlineExceeded) {
		errID := log.Error(ctx, here, err, r)
		writeError(w, err, http.StatusGatewayTimeout, errID)
		return
	} else if goerrors.Is(err, ErrAccessTokenRequired) {
		writeError(w, err, http.StatusNetworkAuthenticationRequired, err.Error())
		return
	} else if goerrors.Is(err, ErrAccessTokenExpired) {
		writeError(w, err, http.StatusPreconditionFailed, err.Error())
		return
	} else if goerrors.Is(err, ErrPaymentTokenRequired) {
		writeError(w, err, http.StatusPaymentRequired, err.Error())
		return
	}
	errID := log.Error(ctx, here, err, r)
	writeError(w, err, http.StatusInternalServerError, errID)
}

//contextWithTokenAndDeadline add token to context if token exist in cookies and deadline
//
//	context,cancel, token, err := contextWithTokenAndDeadline(req)
//
func contextWithTokenAndDeadline(r *http.Request) (context.Context, context.CancelFunc, app.Token, error) {
	dateline := app.ContextDateline()
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

// writeBinary to response
//
//	writeBinary(w, bytes)
//
func writeBinary(w http.ResponseWriter, bytes []byte) {
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(bytes)
}

// writeText to response
//
//	writeText(w, "code")
//
func writeText(w http.ResponseWriter, text string) {
	w.Header().Set("Content-Type", "text/plain")
	io.WriteString(w, text)
}

// writeError to response
//
//	writeError(w, errors.New("error"), 500, "error")
//
func writeError(w http.ResponseWriter, err error, statusCode int, text string) {
	w.WriteHeader(statusCode)
	writeText(w, text)
}

// writeBadRequest to response
//
//	writeBadRequest(context.Background(), w, "message")
//
func writeBadRequest(ctx context.Context, w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusBadRequest)
	writeText(w, msg)
	log.Debug(ctx, here, msg)
}
