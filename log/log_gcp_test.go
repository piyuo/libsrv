package log

import (
	"context"
	"testing"
	"time"

	app "github.com/piyuo/libsrv/app"
	util "github.com/piyuo/libsrv/util"
	. "github.com/smartystreets/goconvey/convey"
)

const here = "log_gcp_test"

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
		errClient, _ := gcpCreateErrorClient(ctx, app.PiyuoID())
		So(errClient, ShouldNotBeNil)
	})
}

func TestGcpErrorOpenWrite(t *testing.T) {
	Convey("should open and write error'", t, func() {
		ctx := context.Background()
		application := "go"
		identity := "test"
		message := "mock error happening in flutter"
		stack := "at firstLine (a.js:3)\nat secondLine (b.js:3)"
		id := util.UUID()
		client, close, err := gcpErrorOpen(ctx, app.PiyuoID())
		So(err, ShouldBeNil)
		defer close()
		gcpErrorWrite(client, message, application, identity, here, stack, id, nil)
	})
}

func TestGcpLogOpenWrite(t *testing.T) {
	Convey("should open and write log'", t, func() {
		ctx := context.Background()
		application, identity := aiFromContext(ctx)
		message := "mock error happening in flutter"
		logger, close, err := gcpLogOpen(ctx)
		So(err, ShouldBeNil)
		defer close()
		gcpLogWrite(logger, time.Now(), message, application, identity, here, LevelInfo)
	})
}
