package server

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/piyuo/libsrv/identifier"
	"github.com/stretchr/testify/assert"
)

func TestGzipEnabled(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Add("Accept-Encoding", "gzip")
	resp := httptest.NewRecorder()
	HTTPEntry(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		str := identifier.RandomNumber(256)
		w.Write([]byte(str))
		return nil
	}).ServeHTTP(resp, req)
	result := resp.Result()
	assert.Equal(http.StatusOK, result.StatusCode)
	assert.Equal("gzip", result.Header.Get("Content-Encoding"))
	l := len(resp.Body.Bytes())
	assert.True(l < 256)
}

func TestGzipNotEnable(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	req, _ := http.NewRequest("GET", "/", nil)
	resp := httptest.NewRecorder()
	HTTPEntry(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		str := identifier.RandomNumber(128)
		w.Write([]byte(str))
		return nil
	}).ServeHTTP(resp, req)
	result := resp.Result()
	assert.Equal(http.StatusOK, result.StatusCode)
	assert.Empty(result.Header.Get("Content-Encoding"))
	l := len(resp.Body.Bytes())
	assert.True(l == 128)
}

func TestGzipSmallChunk(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Add("Accept-Encoding", "gzip")
	resp := httptest.NewRecorder()
	HTTPEntry(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		str := identifier.RandomNumber(128)
		w.Write([]byte(str))
		return nil
	}).ServeHTTP(resp, req)
	result := resp.Result()
	assert.Equal(http.StatusOK, result.StatusCode)
	//	assert.Equal("gzip", result.Header.Get("Content-Encoding"))
	l := len(resp.Body.Bytes())
	assert.True(l == 128)
}

func TestDefaultReturn403(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	req1, _ := http.NewRequest("GET", "/", nil)
	resp1 := httptest.NewRecorder()

	HTTPEntry(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		WriteStatus(w, http.StatusForbidden, "Forbidden")
		return nil
	}).ServeHTTP(resp1, req1)
	res1 := resp1.Result()
	assert.Equal(http.StatusForbidden, res1.StatusCode)

	//cleanup http.Handle mapping
	http.DefaultServeMux = new(http.ServeMux)
}

func TestHandlerReturnError(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	req1, _ := http.NewRequest("GET", "/", nil)
	resp1 := httptest.NewRecorder()
	HTTPEntry(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		return errors.New("myError")
	}).ServeHTTP(resp1, req1)
	res1 := resp1.Result()
	assert.Equal(http.StatusInternalServerError, res1.StatusCode)

	//cleanup http.Handle mapping
	http.DefaultServeMux = new(http.ServeMux)
}

func TestDeadline(t *testing.T) {
	t.Parallel()
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

func TestDeadlineNotSet(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	ctx := context.Background()
	assert.Nil(ctx.Err())

	backup := os.Getenv("DEADLINE_HTTP")
	os.Setenv("DEADLINE_HTTP", "")
	defer os.Setenv("DEADLINE_HTTP", backup)
	deadlineHTTP = -1 // remove cache

	_, cancel := setDeadlineHTTP(ctx)
	defer cancel()

	ms := deadlineHTTP.Milliseconds()
	assert.NotZero(ms)
	deadlineHTTP = -1 // remove cache
}
