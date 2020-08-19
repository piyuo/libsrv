package log

import (
	"context"
	"testing"

	identifier "github.com/piyuo/libsrv/identifier"
	"github.com/piyuo/libsrv/session"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGCPErrorer(t *testing.T) {
	Convey("should write error", t, func() {
		appName = "error-gcp_test"
		ctx := context.Background()
		ctx = session.SetUserID(ctx, "user1")
		errorer, err := NewGCPErrorer(ctx)
		So(err, ShouldBeNil)
		So(errorer, ShouldNotBeNil)
		defer errorer.Close()
		stack := "at firstLine (a.js:3)\nat secondLine (b.js:3)"
		id := identifier.UUID()
		errorer.Write(ctx, "TestGCPLogger", "write error", stack, id)
	})
}

func TestGCPEmptyStack(t *testing.T) {
	Convey("should write empty stack", t, func() {
		appName = "error-gcp_test"
		ctx := context.Background()
		errorer, err := NewGCPErrorer(ctx)
		So(err, ShouldBeNil)
		So(errorer, ShouldNotBeNil)
		defer errorer.Close()
		id := identifier.UUID()
		errorer.Write(ctx, "TestGCPLogger", "write error", "", id)
	})
}
