package log

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

const HERE = "log_gcp_test"

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
