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

		cred, err := createCredential(context.Background(), key, "https://www.googleapis.com/auth/cloud-platform")
		So(err, ShouldBeNil)
		So(cred, ShouldNotBeNil)
	})

	Convey("should keep global log google credential", t, func() {
		So(globalLogCredential, ShouldBeNil)
		cred, err := LogCredential(context.Background())
		So(err, ShouldBeNil)
		So(cred, ShouldNotBeNil)
		So(globalLogCredential, ShouldNotBeNil)
	})

	Convey("should keep global data google credential", t, func() {
		So(globalDataCredential, ShouldBeNil)
		cred, err := GlobalDataCredential(context.Background())
		So(err, ShouldBeNil)
		So(cred, ShouldNotBeNil)
		So(globalDataCredential, ShouldNotBeNil)
	})

}

func TestDataCredentialByRegion(t *testing.T) {
	Convey("should get data credential by region", t, func() {
		cred, err := DataCredentialByRegion(context.Background(), "us")
		So(err, ShouldBeNil)
		So(cred, ShouldNotBeNil)
		cred, err = DataCredentialByRegion(context.Background(), "jp")
		So(err, ShouldBeNil)
		So(cred, ShouldNotBeNil)
		cred, err = DataCredentialByRegion(context.Background(), "be")
		So(err, ShouldBeNil)
		So(cred, ShouldNotBeNil)

	})
}

func TestRegionalDataCredential(t *testing.T) {
	Convey("should get data credential in current region", t, func() {
		cred, err := RegionalDataCredential(context.Background())
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
		_, err := LogCredential(ctx)
		So(err, ShouldNotBeNil)
		_, err2 := GlobalDataCredential(ctx)
		So(err2, ShouldNotBeNil)
	})
}
