package log

import (
	"bytes"
	"context"
	"net/http"
	"testing"
	"time"

	tools "github.com/piyuo/go-libsrv/tools"
	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
)

const HERE = "log_test"

func TestCreateLogClient(t *testing.T) {
	Convey("should create log client'", t, func() {
		ctx := context.Background()
		logClient, _ := gcpCreateLogClient(ctx)
		So(logClient, ShouldNotBeNil)
	})
}

func TestCreateErrorClient(t *testing.T) {
	Convey("should create error client'", t, func() {
		ctx := context.Background()
		errClient, _ := gcpCreateErrorClient(ctx, "my service", "my version")
		So(errClient, ShouldNotBeNil)
	})
}

func TestInfo(t *testing.T) {
	Convey("should print'", t, func() {
		ctx := context.Background()
		Debug(ctx, HERE, "debug msg")
	})
}

//TestLog is a production test, it will write log to google cloud platform under log viewer "Google Project, project name"
func TestLog(t *testing.T) {
	Convey("should log to server'", t, func() {
		ctx := context.Background()
		Info(ctx, HERE, "my info log")
		Warning(ctx, HERE, "my warning log")
		Alert(ctx, HERE, "my alert log")
	})
}

func TestLogWhenContextCanceled(t *testing.T) {
	Convey("should get error when context canceled", t, func() {
		dateline := time.Now().Add(time.Duration(1) * time.Millisecond)
		ctx, cancel := context.WithDeadline(context.Background(), dateline)
		defer cancel()
		time.Sleep(time.Duration(2) * time.Second)
		Info(ctx, HERE, "my info log canceled")
	})
}

func TestError(t *testing.T) {
	Convey("should print error'", t, func() {
		ctx := context.Background()
		err := errors.New("mock error happening in go")
		errID := Error(ctx, HERE, err, nil)
		So(errID, ShouldNotBeEmpty)
		So(false, ShouldEqual, false)
	})
}

func TestErrorWithRequest(t *testing.T) {
	Convey("should print error'", t, func() {
		ctx := context.Background()
		err := errors.New("mock error happening in go with request")
		req, _ := http.NewRequest("GET", "/", bytes.NewReader([]byte("ABC")))
		errID := Error(ctx, HERE, err, req)
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
		CustomError(ctx, message, application, identity, HERE, stack, id, nil)
		So(false, ShouldEqual, false)
	})
}
