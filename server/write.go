package server

import (
	"context"
	"io"
	"net/http"

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
func writeError(w http.ResponseWriter, statusCode int, errID string, err error) {
	w.WriteHeader(statusCode)
	writeText(w, errID+"-"+err.Error())
}

// writeStatus write status code and text response
//
//	writeStatus(w, 500, "error")
//
func writeStatus(w http.ResponseWriter, statusCode int, text string) {
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
