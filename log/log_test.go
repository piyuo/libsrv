package log

import (
	"bytes"
	"context"
	"net/http"
	"os"
	"testing"
	"time"

	tools "github.com/piyuo/go-libsrv/tools"
	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
)

func TestAiFromContext(t *testing.T) {
	Convey("should get application from context'", t, func() {
		ctx := context.Background()
		backupPiyuoApp := os.Getenv("PIYUO_APP")
		os.Setenv("PIYUO_APP", "dev")
		application, identity := aiFromContext(ctx)
		So(application, ShouldEqual, "dev")
		So(identity, ShouldEqual, "")
		os.Setenv("PIYUO_APP", backupPiyuoApp)
	})
}

func TestLogHead(t *testing.T) {
	Convey("should get head from application and identity'", t, func() {
		HERE := "log_test"
		h := head("piyuo-m-us-sys", "user-store", HERE)
		So(h, ShouldEqual, "user-store@piyuo-m-us-sys/log_test: ")
	})
}

func TestDebug(t *testing.T) {
	Convey("should print info to debug console'", t, func() {
		ctx := context.Background()
		Debug(ctx, here, "debug msg")
	})
}

//TestLog is a production test, it will write log to google cloud platform under log viewer "Google Project, project name"
func TestLog(t *testing.T) {
	Convey("should log to server'", t, func() {
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
		Info(ctx, here, "my info log canceled")
	})
}

func TestBeautyStack(t *testing.T) {
	Convey("should return beauty formatted stack '", t, func() {
		err := errors.New("beauty stack")
		stack := beautyStack(err)
		So(stack, ShouldNotBeEmpty)
	})
}

func TestExtractFilename(t *testing.T) {
	Convey("should return filename '", t, func() {
		path := "/convey/doc.go:75"
		filename := extractFilename(path)
		So(filename, ShouldEqual, "doc.go:75")
	})
}

func TestIsLineUsable(t *testing.T) {
	Convey("should check line usable '", t, func() {
		line := "/smartystreets/convey/doc.go:75"
		So(isLineUsable(line), ShouldBeFalse)
	})
}

func TestIsLineDuplicate(t *testing.T) {
	Convey("should check line duplicated '", t, func() {
		list := []string{"/doc.go:75", "/doc.go:75"}
		So(isLineDuplicate(list, 0), ShouldBeFalse)
		So(isLineDuplicate(list, 1), ShouldBeTrue)
	})
}

func TestError(t *testing.T) {
	Convey("should print error'", t, func() {
		ctx := context.Background()
		err := errors.New("mock error happening in go")
		errID := Error(ctx, here, err, nil)
		So(errID, ShouldNotBeEmpty)
		So(false, ShouldEqual, false)
	})
}

func TestErrorWithRequest(t *testing.T) {
	Convey("should print error'", t, func() {
		ctx := context.Background()
		err := errors.New("mock error happening in go with request")
		req, _ := http.NewRequest("GET", "/", bytes.NewReader([]byte("ABC")))
		errID := Error(ctx, here, err, req)
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
		ErrorLog(ctx, message, application, identity, here, stack, id, nil)
		So(false, ShouldEqual, false)
	})
}

func TestErrorOpenWrite(t *testing.T) {
	Convey("should open and write error'", t, func() {
		ctx := context.Background()
		application, identity := aiFromContext(ctx)
		message := "mock error happening in flutter"
		stack := "at firstLine (a.js:3)\nat secondLine (b.js:3)"
		id := tools.UUID()
		client, close, err := ErrorOpen(ctx, application, here)
		So(err, ShouldBeNil)
		defer close()
		ErrorWrite(ctx, client, message, application, identity, here, stack, id, nil)
	})
}

func TestLogOpenWrite(t *testing.T) {
	Convey("should open and write log'", t, func() {
		ctx := context.Background()
		application, _ := aiFromContext(ctx)
		message := "mock error happening in flutter"
		logger, close, err := Open(ctx)
		So(err, ShouldBeNil)
		defer close()
		Write(ctx, logger, message, application, "001-CHIENCHIH", here, info)
	})
}
