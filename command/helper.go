package command

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	log "github.com/piyuo/libsrv/log"
	util "github.com/piyuo/libsrv/util"
)

// keyContext define key used in ctx
//
type keyContext int

const (
	// REQUEST is context key name for request
	//
	keyRequest keyContext = iota

	// TOKEN is context key name for token
	//
	keyToken
)

// commandDateline cache os env COMMAND_DEADLINE value
//
var commandDateline time.Duration = -1

// commandDateline cache os env COMMAND_SLOW value
//
var commandSlow int = -1

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

// GetIP return ip from current request, return empty if anything wrong
//
//	ip := GetIP(ctx)
//
func GetIP(ctx context.Context) string {
	value := ctx.Value(keyRequest)
	if value == nil {
		return ""
	}
	req := value.(*http.Request)
	return util.GetIP(req)
}

// GetLocale return locale from current request, return en-us if anything else
//
//	lang := GetLocale(ctx)
//
func GetLocale(ctx context.Context) string {
	value := ctx.Value(keyRequest)
	if value == nil {
		return "en-us"
	}
	req := value.(*http.Request)
	return util.GetLocale(req)
}

// getDeadline get context deadline,dateline should not greater than 10 min.
//
//	deadline := getDeadline()
//
func getDeadline() time.Time {
	if commandDateline == -1 {
		text := os.Getenv("DEADLINE")
		var err error
		ms, err := strconv.Atoi(text)
		if err != nil {
			ms = 20000
			fmt.Print("use default deadline 20 seconds")
		}
		commandDateline = time.Duration(ms)
	}
	return time.Now().Add(commandDateline * time.Millisecond)
}

//IsSlow check execution time is greater than slow definition,if so return slow limit, other return 0
//
//	So(IsSlow(5), ShouldBeFalse)
func IsSlow(executionTime int) int {
	if commandSlow == -1 {
		text := os.Getenv("SLOW")
		var err error
		commandSlow, err = strconv.Atoi(text)
		if err != nil {
			commandSlow = 12000
			fmt.Print("use default slow detection 12 seconds")
		}
	}

	if executionTime > commandSlow {
		return commandSlow
	}
	return 0
}
