package server

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func mockHTTPHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	WriteStatus(w, http.StatusForbidden, "Forbidden")
	return nil
}

func mockErrorHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	return errors.New("myError")
}

func TestServerHttpDefaultReturn403(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	req1, _ := http.NewRequest("GET", "/", nil)
	resp1 := httptest.NewRecorder()
	HTTPCreateFunc(mockHTTPHandler).ServeHTTP(resp1, req1)
	res1 := resp1.Result()
	assert.Equal(http.StatusForbidden, res1.StatusCode)

	//cleanup http.Handle mapping
	http.DefaultServeMux = new(http.ServeMux)
}

func TestServerHttpHandlerReturnError(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	req1, _ := http.NewRequest("GET", "/", nil)
	resp1 := httptest.NewRecorder()
	HTTPCreateFunc(mockErrorHandler).ServeHTTP(resp1, req1)
	res1 := resp1.Result()
	assert.Equal(http.StatusInternalServerError, res1.StatusCode)

	//cleanup http.Handle mapping
	http.DefaultServeMux = new(http.ServeMux)
}

func TestServerHttpDeadline(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	assert.Nil(ctx.Err())

	backup := os.Getenv("DEADLINE_HTTP")
	os.Setenv("DEADLINE_HTTP", "20")
	defer os.Setenv("DEADLINE_HTTP", backup)
	deadlineHTTP = -1 // remove cache

	ctx, cancel := setDeadlineHTTP(ctx)
	defer cancel()

	assert.Nil(ctx.Err())
	time.Sleep(time.Duration(31) * time.Millisecond)
	assert.NotNil(ctx.Err())

	deadlineHTTP = -1 // remove cache
}

func TestServerHttpDeadlineNotSet(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	assert.Nil(ctx.Err())

	backup := os.Getenv("DEADLINE_HTTP")
	os.Setenv("DEADLINE_HTTP", "")
	defer os.Setenv("DEADLINE_HTTP", backup)
	deadlineHTTP = -1 // remove cache

	ctx, cancel := setDeadlineHTTP(ctx)
	defer cancel()

	ms := deadlineHTTP.Milliseconds()
	assert.Equal(int64(30000), ms)
	deadlineHTTP = -1 // remove cache
}
