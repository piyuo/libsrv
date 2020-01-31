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

const HERE = "log_test"

func TestCreateLogClient(t *testing.T) {
	Convey("should create log client'", t, func() {
		ctx := context.Background()
		logClient, _ := createLogClient(ctx)
		So(logClient, ShouldNotBeNil)
	})
}

func TestCreateErrorClient(t *testing.T) {
	Convey("should create error client'", t, func() {
		ctx := context.Background()
		errClient, _ := createErrorClient(ctx, "my service", "my version")
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
		Critical(ctx, HERE, "my critical log")
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
