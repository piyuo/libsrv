package gcp

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCredential(t *testing.T) {

	Convey("should init google credential", t, func() {
		cred, _ := createCredential(context.Background(), "log-gcp", "https://www.googleapis.com/auth/cloud-platform")
		So(cred, ShouldNotBeNil)
	})

	Convey("should keep google credential", t, func() {
		So(logCred, ShouldBeNil)
		cred, _ := LogCredential(context.Background())
		So(cred, ShouldNotBeNil)
		So(logCred, ShouldNotBeNil)
	})
}
