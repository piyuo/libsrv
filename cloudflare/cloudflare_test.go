package data

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewCloudflare(t *testing.T) {
	Convey("should new cloudflare", t, func() {
		cflare, err := NewCloudflare(context.Background())
		So(err, ShouldBeNil)
		So(cflare, ShouldNotBeNil)
	})
}

func TestAddSubDomain(t *testing.T) {
	Convey("should add sub domain", t, func() {
		cflare, err := NewCloudflare(context.Background())
		So(err, ShouldBeNil)
		So(cflare, ShouldNotBeNil)
	})
}
