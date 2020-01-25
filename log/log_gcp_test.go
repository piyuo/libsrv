package log

import (
	"bytes"
	"context"
	"net/http"
	"testing"

	tools "github.com/piyuo/go-libsrv/tools"
	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
)

func TestInfo(t *testing.T) {
	Convey("should print'", t, func() {
		ctx := context.Background()
		Info(ctx, "hello log")
	})
}

//TestLog is a production test, it will write log to google cloud platform under log viewer "Google Project, project name"
func TestLog(t *testing.T) {
	Convey("should log to server'", t, func() {
		ctx := context.Background()
		Notice(ctx, "my notice log")
		Warning(ctx, "my warning log")
		Alert(ctx, "my alert log")
	})
}

func TestError(t *testing.T) {
	Convey("should print error'", t, func() {
		ctx := context.Background()
		err := errors.New("mock error happening in go")
		errID := Error(ctx, err, nil)
		So(errID, ShouldNotBeEmpty)
		So(false, ShouldEqual, false)
	})
}

func TestErrorWithRequest(t *testing.T) {
	Convey("should print error'", t, func() {
		ctx := context.Background()
		err := errors.New("mock error happening in go with request")
		req, _ := http.NewRequest("GET", "/", bytes.NewReader([]byte("ABC")))
		errID := Error(ctx, err, req)
		So(errID, ShouldNotBeEmpty)
		So(false, ShouldEqual, false)
	})
}

func TestCustomError(t *testing.T) {
	Convey("should print error from'", t, func() {
		ctx := context.Background()
		application, identity := aiFromContext(ctx)
		message := "mock error happening in flutter"
		stack := "at firstLine (a.js:3)\nat secondLine (b.js:3)"
		id := tools.UUID()
		CustomError(ctx, message, application, identity, stack, id, true, nil)
		So(false, ShouldEqual, false)
	})
}
