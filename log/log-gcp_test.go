package log

import (
	"context"
	"os"
	"testing"

	tools "github.com/piyuo/go-libsrv/tools"
	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGetLogHead(t *testing.T) {
	Convey("should generate log head'", t, func() {
		ctx := context.Background()
		backupPiyuoApp := os.Getenv("PIYUO_APP")
		os.Setenv("PIYUO_APP", "dev")
		head, _, _ := generateLogHead(ctx, false)
		So(head, ShouldEqual, "[dev]: ")
		head, _, _ = generateLogHead(ctx, true)
		So(head, ShouldEqual, "<dev>: ")

		os.Setenv("PIYUO_APP", "piyuo-m-us-sys")
		head, _, _ = generateLogHead(ctx, false)
		So(head, ShouldEqual, "[piyuo-m-us-sys]: ")
		os.Setenv("PIYUO_APP", "piyuo-m-web-index")
		head, _, _ = generateLogHead(ctx, true)
		So(head, ShouldEqual, "<piyuo-m-web-index>: ")
		os.Setenv("PIYUO_APP", backupPiyuoApp)
	})
}

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

// TestError is a production test, it will write error to google cloud platform under Error Reporting
func TestErr(t *testing.T) {
	Convey("should print error'", t, func() {
		ctx := context.Background()
		err := errors.New("mock error happening in go")
		errID := Error(ctx, err)
		So(errID, ShouldNotBeEmpty)
		So(false, ShouldEqual, false)
	})
}

func TestError(t *testing.T) {
	Convey("should print error from'", t, func() {
		ctx := context.Background()
		message := "mock error happening in flutter"
		stack := "at firstLine (a.js:3)\nat secondLine (b.js:3)"
		id := tools.UUID()
		CustomError(ctx, message, stack, id, true)
		So(false, ShouldEqual, false)
	})
}
