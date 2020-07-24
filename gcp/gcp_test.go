package gcp

import (
	"context"
	"testing"
	"time"

	app "github.com/piyuo/libsrv/app"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCredential(t *testing.T) {

	Convey("should init google credential", t, func() {
		key, err := app.Key("gcloud")
		So(err, ShouldBeNil)

		cred, err := createCredential(context.Background(), key)
		So(err, ShouldBeNil)
		So(cred, ShouldNotBeNil)

		// test multi scope
		cred, err = createCredential(context.Background(), key)
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
		cred, err := RegionalCredential(context.Background(), "us")
		So(err, ShouldBeNil)
		So(cred, ShouldNotBeNil)
		cred, err = RegionalCredential(context.Background(), "jp")
		So(err, ShouldBeNil)
		So(cred, ShouldNotBeNil)
		cred, err = RegionalCredential(context.Background(), "be")
		So(err, ShouldBeNil)
		So(cred, ShouldNotBeNil)

	})
}

func TestRegionalDataCredential(t *testing.T) {
	Convey("should get data credential in current region", t, func() {
		cred, err := CurrentRegionalCredential(context.Background())
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
	})
}
