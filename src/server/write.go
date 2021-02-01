package server

import (
	"io"
	"net/http"
)

// commandDateline cache os env COMMAND_SLOW value
//
var commandSlow int = -1

// WriteBinary to response
//
//	WriteBinary(w, bytes)
//
func WriteBinary(w http.ResponseWriter, bytes []byte) {
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(bytes)
}

// WriteText to response
//
//	WriteText(w, "code")
//
func WriteText(w http.ResponseWriter, text string) {
	w.Header().Set("Content-Type", "text/plain")
	io.WriteString(w, text)
}

// WriteError to response
//
//	WriteError(w, errors.New("error"), 500, "error")
//
func WriteError(w http.ResponseWriter, statusCode int, errID string, err error) {
	w.WriteHeader(statusCode)
	WriteText(w, errID+"-"+err.Error())
}

// WriteStatus write status code and text response
//
//	WriteStatus(w, 500, "error")
//
func WriteStatus(w http.ResponseWriter, statusCode int, text string) {
	w.WriteHeader(statusCode)
	WriteText(w, text)
}
