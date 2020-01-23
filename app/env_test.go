package libsrv

import (
	"context"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestBasicFunction(t *testing.T) {
	Convey("should able join dir and current dir'", t, func() {
		text := EnvJoinCurrentDir("../../../")
		So(strings.HasSuffix(text, "/go"), ShouldBeTrue)
	})
	Convey("should env PIYUO_ENV'", t, func() {
		id := EnvPiyuoApp()
		So(id, ShouldEqual, "piyuo-dev")
		So(EnvProduction(), ShouldEqual, false)
	})
	Convey("should get key path'", t, func() {
		path, err := EnvKeyPath("log")
		So(err, ShouldBeNil)
		So(strings.HasSuffix(path, "/log.key"), ShouldBeTrue)
	})
}

func TestCredential(t *testing.T) {
	Convey("should get attributes from credential", t, func() {
		keyname, scope := getAttributesFromCredential(LOG)
		So(keyname, ShouldEqual, "log")
		So(scope, ShouldNotBeEmpty)
	})

	Convey("should init google credential", t, func() {
		cred, _ := createGoogleCloudCredential(context.Background(), LOG)
		So(cred, ShouldNotBeNil)
	})

	Convey("should keep google credential", t, func() {
		So(logCred, ShouldBeNil)
		cred, _ := EnvGoogleCredential(context.Background(), LOG)
		So(cred, ShouldNotBeNil)
		So(logCred, ShouldNotBeNil)
	})
}

func TestUUID(t *testing.T) {
	Convey("should generate uuid", t, func() {
		id := UUID()
		So(id, ShouldNotBeEmpty)
	})

}
