package command

import (
	"bytes"
	"context"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/piyuo/libsrv/command/mock"
	"github.com/piyuo/libsrv/command/shared"
	"github.com/stretchr/testify/assert"
)

var textLong = `{
    "_id": "55d26da7c3f96f90aa005",
    "age": 20,
    "gender": "female",
    "company": "ZOGAK",
    "phone": "+1 (915) 479-2908"
   `

func BenchmarkBigArchive(b *testing.B) {
	handler := newTestServerHandler()
	actBytes := newTestAction(textLong)
	req1, _ := http.NewRequest("GET", "/", bytes.NewReader(actBytes))
	req1.Header.Set("Accept-Encoding", "gzip")
	resp1 := httptest.NewRecorder()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		handler.ServeHTTP(resp1, req1)
		_ = resp1.Result()
	}
}

func BenchmarkSmallAction(b *testing.B) {
	handler := newTestServerHandler()
	actBytes := newTestAction("Hi")
	req1, _ := http.NewRequest("GET", "/", bytes.NewReader(actBytes))
	req1.Header.Set("Accept-Encoding", "gzip")
	resp1 := httptest.NewRecorder()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		handler.ServeHTTP(resp1, req1)
		_ = resp1.Result()
	}
}

func newBigDataAction() (*mock.BigDataAction, []byte) {
	dispatch := &Dispatch{
		Map: &mock.MapXXX{},
	}
	act := &mock.BigDataAction{}
	actBytes, _ := dispatch.encodeCommand(act.XXX_MapID(), act)
	return act, actBytes
}

func TestArchive(t *testing.T) {
	assert := assert.New(t)
	handler := newTestServerHandler()

	act, actBytes := newBigDataAction()
	sampleBytes := []byte(act.GetSample())
	sampleLen := len(sampleBytes)
	req1, _ := http.NewRequest("GET", "/", bytes.NewReader(actBytes))
	req1.Header.Set("Accept-Encoding", "gzip")
	resp1 := httptest.NewRecorder()
	handler.ServeHTTP(resp1, req1)
	res1 := resp1.Result()
	returnBytes := resp1.Body.Bytes()
	returnLen := len(returnBytes)
	assert.Equal(200, res1.StatusCode)
	assert.Greater(returnLen, 10)
	assert.Greater(sampleLen, returnLen)
	assert.Equal("gzip", res1.Header.Get("Content-Encoding"))
	assert.Equal("application/octet-stream", res1.Header.Get("Content-Type"))
}

func customHTTPHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) (bool, error) {
	w.WriteHeader(http.StatusOK)
	writeText(w, "hello")
	return true, nil
}

func TestHTTPHandler(t *testing.T) {
	assert := assert.New(t)
	server := &Server{
		Map:         &mock.MapXXX{},
		HTTPHandler: customHTTPHandler,
	}
	handler := server.newHandler()

	req1, _ := http.NewRequest("GET", "/", nil)
	req1.Header.Set("Accept-Encoding", "gzip")
	resp1 := httptest.NewRecorder()
	handler.ServeHTTP(resp1, req1)
	res1 := resp1.Result()
	returnBytes := resp1.Body.Bytes()
	bodyString := string(returnBytes)
	assert.Equal(200, res1.StatusCode)
	assert.Equal("hello", bodyString)
}

func TestServe(t *testing.T) {
	assert := assert.New(t)
	handler := newTestServerHandler()
	actBytes := newTestAction("Hi")
	req1, _ := http.NewRequest("GET", "/", bytes.NewReader(actBytes))
	resp1 := httptest.NewRecorder()
	handler.ServeHTTP(resp1, req1)
	res1 := resp1.Result()

	returnBytes := resp1.Body.Bytes()
	returnLen := len(returnBytes)
	ok := okResponse()
	okLen := len(ok)
	assert.Equal(200, res1.StatusCode)
	assert.Equal(okLen, returnLen)
	assert.Equal(ok[0], returnBytes[0])
	assert.Equal("application/octet-stream", res1.Header.Get("Content-Type"))
}

func TestServe404(t *testing.T) {
	assert := assert.New(t)
	handler := newTestServerHandler()
	req1, _ := http.NewRequest("GET", "/favicon.ico", nil)
	resp1 := httptest.NewRecorder()
	handler.ServeHTTP(resp1, req1)
	res1 := resp1.Result()
	assert.Equal(400, res1.StatusCode)
}

func newTestServerHandler() http.Handler {
	server := &Server{
		Map: &mock.MapXXX{},
	}
	server.dispatch = &Dispatch{
		Map: server.Map,
	}

	return server.newHandler()
}

func newTestAction(text string) []byte {
	dispatch := &Dispatch{
		Map: &mock.MapXXX{},
	}
	act := &mock.RespondAction{
		Text: text,
	}
	actBytes, _ := dispatch.encodeCommand(act.XXX_MapID(), act)
	return actBytes
}

func okResponse() []byte {
	dispatch := &Dispatch{
		Map: &shared.MapXXX{},
	}
	ok := OK().(*shared.PbOK)
	bytes, _ := dispatch.encodeCommand(ok.XXX_MapID(), ok)
	return bytes
}

func TestContextCanceled(t *testing.T) {
	assert := assert.New(t)
	dateline := time.Now().Add(time.Duration(1) * time.Millisecond)
	ctx, cancel := context.WithDeadline(context.Background(), dateline)
	defer cancel()

	assert.Nil(ctx.Err())
	time.Sleep(time.Duration(2) * time.Millisecond)
	assert.NotNil(ctx.Err())

	err := ctx.Err()
	assert.True(errors.Is(err, context.DeadlineExceeded))
}

func newDeadlineAction() []byte {
	dispatch := &Dispatch{
		Map: &mock.MapXXX{},
	}
	act := &mock.DeadlineAction{}
	actBytes, _ := dispatch.encodeCommand(act.XXX_MapID(), act)
	return actBytes
}

func TestServeWhenContextCanceled(t *testing.T) {
	assert := assert.New(t)
	handler := newTestServerHandler()
	actBytes := newDeadlineAction()
	req, _ := http.NewRequest("GET", "/", bytes.NewReader(actBytes))
	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	res := resp.Result()
	assert.Equal(504, res.StatusCode)
}

func TestHandleRouteException(t *testing.T) {
	//r, _ := http.NewRequest("POST", "/", nil)
	w := httptest.NewRecorder()
	handleRouteException(context.Background(), w, context.DeadlineExceeded)
	handleRouteException(context.Background(), w, errors.New(""))
}

func TestServer(t *testing.T) {
	assert := assert.New(t)
	defer func() {
		if err := recover(); err != nil {
			log.Println("panic occurred:", err)
		}
	}()
	server := &Server{
		Map:         nil,
		HTTPHandler: customHTTPHandler,
	}
	assert.Panics(server.Start)
}
