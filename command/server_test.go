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

	mock "github.com/piyuo/libsrv/command/mock"
	shared "github.com/piyuo/libsrv/command/shared"

	. "github.com/smartystreets/goconvey/convey"
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
	Convey("test any file request", t, func() {
		So(res1.StatusCode, ShouldEqual, 200)
		So(returnLen, ShouldBeGreaterThan, 10)
		So(sampleLen, ShouldBeGreaterThan, returnLen)
		So(res1.Header.Get("Content-Encoding"), ShouldEqual, "gzip")
		So(res1.Header.Get("Content-Type"), ShouldEqual, "application/octet-stream")
	})
}

func customHTTPHandler(w http.ResponseWriter, r *http.Request) (bool, error) {
	w.WriteHeader(http.StatusOK)
	writeText(w, "hello")
	return true, nil
}

func TestHTTPHandler(t *testing.T) {
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
	Convey("should use custom http handler", t, func() {
		So(res1.StatusCode, ShouldEqual, 200)
		So(bodyString, ShouldEqual, "hello")
	})
}

func TestServe(t *testing.T) {
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
	Convey("test any file request", t, func() {
		So(res1.StatusCode, ShouldEqual, 200)
		So(returnLen, ShouldEqual, okLen)
		So(returnBytes[0], ShouldEqual, ok[0])
		So(res1.Header.Get("Content-Type"), ShouldEqual, "application/octet-stream")
	})
}

func TestServe404(t *testing.T) {
	handler := newTestServerHandler()
	req1, _ := http.NewRequest("GET", "/favicon.ico", nil)
	resp1 := httptest.NewRecorder()
	handler.ServeHTTP(resp1, req1)
	res1 := resp1.Result()
	Convey("test any file request", t, func() {
		So(res1.StatusCode, ShouldEqual, 400)
	})
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
	Convey("should get error when context canceled", t, func() {
		dateline := time.Now().Add(time.Duration(1) * time.Millisecond)
		ctx, cancel := context.WithDeadline(context.Background(), dateline)
		defer cancel()

		So(ctx.Err(), ShouldBeNil)
		time.Sleep(time.Duration(2) * time.Millisecond)
		So(ctx.Err(), ShouldNotBeNil)

		err := ctx.Err()
		So(errors.Is(err, context.DeadlineExceeded), ShouldBeTrue)
	})
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
	Convey("should get error when context canceled", t, func() {
		handler := newTestServerHandler()
		actBytes := newDeadlineAction()
		req, _ := http.NewRequest("GET", "/", bytes.NewReader(actBytes))
		resp := httptest.NewRecorder()
		handler.ServeHTTP(resp, req)
		res := resp.Result()
		So(res.StatusCode, ShouldEqual, 504)
	})
}

func TestHandleRouteException(t *testing.T) {
	Convey("should handle route exception", t, func() {
		//r, _ := http.NewRequest("POST", "/", nil)
		w := httptest.NewRecorder()
		handleRouteException(context.Background(), w, context.DeadlineExceeded)
		handleRouteException(context.Background(), w, errors.New(""))
	})
}

func TestServer(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("panic occurred:", err)
		}
	}()
	Convey("should start server", t, func() {
		server := &Server{
			Map:         nil,
			HTTPHandler: customHTTPHandler,
		}
		So(server.Start, ShouldPanic)
	})
}
