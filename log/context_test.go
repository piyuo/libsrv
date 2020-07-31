package log

import (
	"bytes"
	"context"
	"net/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetRequest(t *testing.T) {
	Convey("should get request", t, func() {
		req, _ := http.NewRequest("GET", "/", bytes.NewReader([]byte("ABC")))
		ctx := context.Background()
		ctx = context.WithValue(ctx, keyRequest, req)
		req2 := getRequest(ctx)
		So(req2, ShouldEqual, req)
	})
}

func TestGetID(t *testing.T) {
	Convey("should get ID", t, func() {
		ctx := context.Background()
		ctx = context.WithValue(ctx, keyToken, map[string]string{"id": "user1"})
		id := getID(ctx)
		So(id, ShouldEqual, "user1")
	})
}
