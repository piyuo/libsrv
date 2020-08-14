package gcp

import (
	"context"
	"testing"
	"time"

	"github.com/piyuo/libsrv/key"
	"github.com/piyuo/libsrv/region"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCredential(t *testing.T) {

	Convey("should create google credential", t, func() {
		bytes, err := key.BytesWithoutCache("gcloud.json")
		So(err, ShouldBeNil)

		cred, err := createCredential(context.Background(), bytes)
		So(err, ShouldBeNil)
		So(cred, ShouldNotBeNil)
	})

	Convey("should keep global credential", t, func() {
		So(globalCredential, ShouldBeNil)
		cred, err := GlobalCredential(context.Background())
		So(err, ShouldBeNil)
		So(cred, ShouldNotBeNil)
		So(globalCredential, ShouldNotBeNil)
	})
}

func TestDataCredentialByRegion(t *testing.T) {
	Convey("should get data credential by region", t, func() {
		region.Current = "us"
		cred, err := RegionalCredential(context.Background())
		So(err, ShouldBeNil)
		So(cred, ShouldNotBeNil)

		region.Current = "jp"
		cred, err = RegionalCredential(context.Background())
		So(err, ShouldBeNil)
		So(cred, ShouldNotBeNil)

		region.Current = "be"
		cred, err = RegionalCredential(context.Background())
		So(err, ShouldBeNil)
		So(cred, ShouldNotBeNil)

	})
}

func TestCredentialWhenContextCanceled(t *testing.T) {
	Convey("should get error when context canceled", t, func() {
		dateline := time.Now().Add(time.Duration(1) * time.Millisecond)
		ctx, cancel := context.WithDeadline(context.Background(), dateline)
		defer cancel()
		time.Sleep(time.Duration(2) * time.Millisecond)
		_, err := GlobalCredential(ctx)
		So(err, ShouldNotBeNil)
		_, err = RegionalCredential(ctx)
		So(err, ShouldNotBeNil)
	})
}
