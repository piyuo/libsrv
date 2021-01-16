package command

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestWriteResponse(t *testing.T) {
	w := httptest.NewRecorder()
	bytes := newTestAction(textLong)
	writeBinary(w, bytes)
	writeText(w, "code")
	writeError(w, errors.New("error"), 500, "error")
	writeBadRequest(context.Background(), w, "message")
}

func TestIsSlow(t *testing.T) {
	assert := assert.New(t)
	// 3 seconds execution time is not slow
	assert.Equal(0, IsSlow(5000))
	// 20 seconds execution time is really slow
	assert.Greater(IsSlow(20000000), 5000)
}
