package server

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/pkg/errors"
)

func TestWriteResponse(t *testing.T) {
	w := httptest.NewRecorder()
	bytes := newTestAction(textLong)
	writeBinary(w, bytes)
	writeText(w, "code")
	writeError(w, 500, "error", errors.New("error"))
	writeBadRequest(context.Background(), w, "message")
}
