package gcp

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCredential(t *testing.T) {

	Convey("should init google credential", t, func() {
		cred, err := createCredential(context.Background(), "log-gcp", "https://www.googleapis.com/auth/cloud-platform")
		So(err, ShouldBeNil)
		So(cred, ShouldNotBeNil)
	})

	Convey("should keep global log google credential", t, func() {
		So(globalCredentialLog, ShouldBeNil)
		cred, err := LogCredential(context.Background())
		So(err, ShouldBeNil)
		So(cred, ShouldNotBeNil)
		So(globalCredentialLog, ShouldNotBeNil)
	})

	Convey("should keep global data google credential", t, func() {
		So(globalCredentialData, ShouldBeNil)
		cred, err := DataCredential(context.Background())
		So(err, ShouldBeNil)
		So(cred, ShouldNotBeNil)
		So(globalCredentialData, ShouldNotBeNil)
	})

}
