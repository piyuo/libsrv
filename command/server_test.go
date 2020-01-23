package command

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/piyuo/go-libsrv/shared"
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

/* archive no longer support
func TestArchive(t *testing.T) {
	handler := newTestServerHandler()
	actBytes := newTestAction(textLong)
	actBytesLen := len(actBytes)
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
		So(actBytesLen, ShouldBeGreaterThan, returnLen)
		So(res1.Header.Get("Content-Encoding"), ShouldEqual, "gzip")
		So(res1.Header.Get("Content-Type"), ShouldEqual, "application/octet-stream")
	})
}
*/
func TestServe(t *testing.T) {
	handler := newTestServerHandler()
	actBytes := newTestAction("Hi")
	req1, _ := http.NewRequest("GET", "/", bytes.NewReader(actBytes))
	resp1 := httptest.NewRecorder()
	handler.ServeHTTP(resp1, req1)
	res1 := resp1.Result()

	returnBytes := resp1.Body.Bytes()
	returnLen := len(returnBytes)
	sample := newTestResponse()
	sampleLen := len(sample)
	Convey("test any file request", t, func() {
		So(res1.StatusCode, ShouldEqual, 200)
		So(returnLen, ShouldEqual, sampleLen)
		So(returnBytes[0], ShouldEqual, sample[0])
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
		Map: &TestMap{},
	}
	return server.newHandler()
}

func newTestAction(text string) []byte {
	dispatch := &Dispatch{
		Map: &TestMap{},
	}
	act := &TestAction{
		Text: text,
	}
	actBytes, _ := dispatch.encodeCommand(act.XXX_MapID(), act)
	return actBytes
}

func newTestResponse() []byte {
	dispatch := &Dispatch{
		Map: &TestMap{},
	}
	r := &shared.Err{
		Code: 1,
	}

	rBytes, _ := dispatch.encodeCommand(r.XXX_MapID(), r)
	return rBytes
}
