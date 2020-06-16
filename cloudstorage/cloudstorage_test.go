package data

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewCloudstorage(t *testing.T) {
	Convey("should new cloudstorage", t, func() {
		cflare, err := NewCloudstorage(context.Background())
		So(err, ShouldBeNil)
		So(cflare, ShouldNotBeNil)
	})
}

func TestAddBucket(t *testing.T) {
	Convey("should new cloudstorage", t, func() {
		cflare, err := NewCloudstorage(context.Background())
		So(err, ShouldBeNil)
		So(cflare, ShouldNotBeNil)
	})
}
