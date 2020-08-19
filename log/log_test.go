package log

import (
	"context"
	"os"
	"testing"
	"time"

	identifier "github.com/piyuo/libsrv/identifier"
	"github.com/piyuo/libsrv/session"
	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
)

var here = "log_test"

func TestShouldPrintToConsole(t *testing.T) {
	Convey("should not print to console if it's stable branch", t, func() {
		os.Setenv("BRANCH", "master")
		So(shouldPrintToConsole(), ShouldBeTrue)
		os.Setenv("BRANCH", "stable")
		So(shouldPrintToConsole(), ShouldBeFalse)
	})
}

func TestGetHeader(t *testing.T) {
	Convey("should get header", t, func() {
		appName = "test"
		ctx := context.Background()
		ctx = session.SetUserID(ctx, "user1")
		header, id := getHeader(ctx, here)
		So(header, ShouldEqual, "user1@test/log_test: ")
		So(id, ShouldEqual, "user1")
	})
}

func TestDebug(t *testing.T) {
	Convey("should print info to debug console", t, func() {
		ctx := context.Background()
		Debug(ctx, here, "debug...")
	})
}

//TestLog is a production test, it will write log to google cloud platform under log viewer "Google Project, project name"
func TestLog(t *testing.T) {
	Convey("should log to server", t, func() {
		ctx := context.Background()
		Info(ctx, here, "my info log")
		Warning(ctx, here, "my warning log")
		Alert(ctx, here, "my alert log")
	})
}

func TestLogWhenContextCanceled(t *testing.T) {
	Convey("should get error when context canceled", t, func() {
		dateline := time.Now().Add(time.Duration(1) * time.Millisecond)
		ctx, cancel := context.WithDeadline(context.Background(), dateline)
		defer cancel()
		time.Sleep(time.Duration(2) * time.Millisecond)
		logger, err := NewLogger(ctx)
		So(err, ShouldNotBeNil)
		So(logger, ShouldBeNil)

		Log(ctx, DEBUG, here, "")
		WriteError(ctx, here, "", "", "")
		errID := Error(ctx, here, nil)
		So(errID, ShouldBeEmpty)
		Info(ctx, here, "my info log canceled")

		errorer, err := NewErrorer(ctx)
		So(err, ShouldNotBeNil)
		So(errorer, ShouldBeNil)
		logger, err = NewLogger(ctx)
		So(err, ShouldNotBeNil)
		So(logger, ShouldBeNil)
	})
}

func TestBeautyStack(t *testing.T) {
	Convey("should return beauty formatted stack", t, func() {
		err := errors.New("beauty stack")
		stack := beautyStack(err)
		So(stack, ShouldNotBeEmpty)
	})
}

func TestExtractFilename(t *testing.T) {
	Convey("should return filename", t, func() {
		path := "/convey/doc.go:75"
		filename := extractFilename(path)
		So(filename, ShouldEqual, "doc.go:75")
		path = "doc.go:75"
		filename = extractFilename(path)
		So(filename, ShouldEqual, "doc.go:75")
	})
}

func TestIsLineUsable(t *testing.T) {
	Convey("should check line usable", t, func() {
		line := "/smartystreets/convey/doc.go:75"
		So(isLineUsable(line), ShouldBeFalse)
	})
}

func TestIsLineDuplicate(t *testing.T) {
	Convey("should check line duplicated", t, func() {
		list := []string{"/doc.go:75", "/doc.go:75"}
		So(isLineDuplicate(list, 0), ShouldBeFalse)
		So(isLineDuplicate(list, 1), ShouldBeTrue)
	})
}

func TestError(t *testing.T) {
	Convey("should print error", t, func() {
		ctx := context.Background()
		err := errors.New("mock error happening in go")
		errID := Error(ctx, here, err)
		So(errID, ShouldNotBeEmpty)

		errID = Error(ctx, here, nil)
		So(errID, ShouldBeEmpty)

	})
}

func TestErrorWithRequest(t *testing.T) {
	Convey("should error", t, func() {
		ctx := context.Background()
		err := errors.New("mock error happening in go with request")
		errID := Error(ctx, here, err)
		So(errID, ShouldNotBeEmpty)
	})
}

func TestCustomError(t *testing.T) {
	Convey("should write error", t, func() {
		ctx := context.Background()
		message := "mock error happening in flutter"
		stack := "at firstLine (a.js:3)\nat secondLine (b.js:3)"
		id := identifier.UUID()
		WriteError(ctx, here, message, stack, id)
		So(false, ShouldEqual, false)
	})
}
