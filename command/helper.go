package command

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/piyuo/libsrv/log"
)

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
