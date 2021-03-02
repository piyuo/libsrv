package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pkg/errors"
)

func TestWriteResponse(t *testing.T) {
	w := httptest.NewRecorder()
	bytes := newTestAction(textLong)
	WriteBinary(w, bytes)
	WriteText(w, "code")
	WriteError(w, 500, errors.New("error"))
	WriteStatus(w, http.StatusBadRequest, "bad request")
}
