package libsrv

import (
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSystem(t *testing.T) {

	Convey("should has only one instance and pass check'", t, func() {
		s1 := Sys()
		s2 := Sys()
		So(s1.IsProduction(), ShouldEqual, s2.IsProduction())
		Sys().Check()
		So(s1.IsProduction(), ShouldEqual, false)
	})

	Convey("should able join dir and current dir'", t, func() {
		text := Sys().JoinCurrentDir("../../")
		So(strings.HasSuffix(text, "/go"), ShouldEqual, true)
	})
}

func TestGetID(t *testing.T) {
	Convey("should get id'", t, func() {
		id := Sys().ID()
		So(id, ShouldEqual, "dev")
	})
}

func TestGetLogHead(t *testing.T) {
	Convey("should get log head'", t, func() {
		s := &system{}
		So(s.getLogHead("dev", ""), ShouldEqual, "<dev>")
		So(s.getLogHead("dev", "user-store"), ShouldEqual, "<dev> user-store")
		So(s.getLogHead("PIYUO-US-M-SYS", ""), ShouldEqual, "[PIYUO-US-M-SYS]")
		So(s.getLogHead("PIYUO-US-M-SYS", "user-store"), ShouldEqual, "[PIYUO-US-M-SYS] user-store")
		So(s.getLogHead("piyuo-us-m-web-index", ""), ShouldEqual, "(piyuo-us-m-web-index)")
		So(s.getLogHead("piyuo-us-m-web-index", "user-store"), ShouldEqual, "(piyuo-us-m-web-index) user-store")
	})
}

func TestCheckProduction(t *testing.T) {
	s := &system{}
	Convey("should set production correctly", t, func() {
		So(s.checkProduction("dev"), ShouldEqual, false)
		So(s.checkProduction("piyuo-tw-a-app"), ShouldEqual, false)
		So(s.checkProduction("piyuo-tw-b-app"), ShouldEqual, false)
		So(s.checkProduction("piyuo-tw-m-app"), ShouldEqual, true)
		So(s.checkProduction("PIYUO-TW-M-SYS"), ShouldEqual, true)
	})
}

func TestCredential(t *testing.T) {
	s := &system{}
	Convey("should get attributes from credential", t, func() {
		filename, scope := s.getAttributesFromCredential(LOG)
		So(filename, ShouldEqual, "log.key")
		So(scope, ShouldNotBeEmpty)
	})

	Convey("should init google credential", t, func() {
		cred, _ := s.initGoogleCloudCredential(LOG)
		So(cred, ShouldNotBeNil)
	})

	Convey("should keep google credential", t, func() {
		So(s.googleCred, ShouldBeNil)
		cred, _ := s.GetGoogleCloudCredential(LOG)
		So(cred, ShouldNotBeNil)
		So(s.googleCred, ShouldNotBeNil)
	})
}

func TestInfo(t *testing.T) {
	Convey("should print'", t, func() {
		Sys().Info("hello log")
		So(true, ShouldEqual, true)
	})
}

/*

//TestLog is a production test, it will write log to google cloud platform under log viewer "Google Project, project name"
func TestLog(t *testing.T) {
	Convey("should print'", t, func() {
		Sys().Notice("my notice log")
		Sys().Warning("my warning log")
		Sys().Alert("my alert log")
	})
}

// TestError is a production test, it will write error to google cloud platform under Error Reporting
func TestError(t *testing.T) {
	Convey("should print error'", t, func() {
		err := errors.New("my error1")
		Sys().Error(err)
		So(false, ShouldEqual, false)
	})
}

func TestErrorBy5(t *testing.T) {
	Convey("should print error with cause by id'", t, func() {
		err := errors.New("my error by user5")
		Sys().ErrorBy(err, "user5")
		So(false, ShouldEqual, false)
	})
}

func TestErrorFrom(t *testing.T) {
	Convey("should print error from'", t, func() {
		stack := "at firstLine (a.js:3)\nat secondLine (b.js:3)"
		location := "piyuo-tw-m-web-index"
		id := "user1"
		language := "flutter"
		Sys().ErrorFrom("error manually", stack, location, id, language)
		So(false, ShouldEqual, false)
	})
}
*/
