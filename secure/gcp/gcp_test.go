package gcp

import (
	"context"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCredential(t *testing.T) {

	Convey("should init google credential", t, func() {
		cred, err := createCredential(context.Background(), "gcloud", "https://www.googleapis.com/auth/cloud-platform")
		So(err, ShouldBeNil)
		So(cred, ShouldNotBeNil)
	})

	Convey("should keep global log google credential", t, func() {
		So(logCredGlobal, ShouldBeNil)
		cred, err := LogCredential(context.Background())
		So(err, ShouldBeNil)
		So(cred, ShouldNotBeNil)
		So(logCredGlobal, ShouldNotBeNil)
	})

	Convey("should keep global data google credential", t, func() {
		So(dataCredGlobal, ShouldBeNil)
		cred, err := DataCredential(context.Background())
		So(err, ShouldBeNil)
		So(cred, ShouldNotBeNil)
		So(dataCredGlobal, ShouldNotBeNil)
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
		_, err2 := DataCredential(ctx)
		So(err2, ShouldNotBeNil)
	})
}
