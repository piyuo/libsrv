package libsrv

import (
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSystem(t *testing.T) {
	Convey("should pass Check'", t, func() {
		//if any thing wrong, check will result panic
		CurrentSystem().Check()
	})

	Convey("should has only one instance'", t, func() {
		s1 := CurrentSystem()
		s2 := CurrentSystem()
		So(s1.IsDebug(), ShouldEqual, s2.IsDebug())
		s1.SetDebug(false)
		So(s1.IsDebug(), ShouldEqual, false)
		So(s2.IsDebug(), ShouldEqual, false)
	})

	Convey("should able join dir and current dir'", t, func() {
		text := CurrentSystem().JoinCurrentDir("../../")
		So(strings.HasSuffix(text, "/go"), ShouldEqual, true)
	})
}

func TestGetID(t *testing.T) {
	Convey("should get id'", t, func() {
		id := CurrentSystem().ID()
		So(id, ShouldEqual, "dev")
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
		cred, _ := s.getGoogleCloudCredential(LOG)
		So(cred, ShouldNotBeNil)
		So(s.googleCred, ShouldNotBeNil)
	})
}

func TestInfo(t *testing.T) {
	Convey("should print'", t, func() {
		CurrentSystem().Info("hello log")
		So(true, ShouldEqual, true)
	})
}

/*
//TestLog is a production test, it will write log to google cloud platform under log viewer "Google Project, project name"
func TestLog(t *testing.T) {
	Convey("should print'", t, func() {
		CurrentSystem().Log("my notice log")
		CurrentSystem().Warning("my warning log")
		CurrentSystem().Alert("my alert log")
		So(false, ShouldEqual, false)
	})
}

// TestError is a production test, it will write error to google cloud platform under Error Reporting
func TestError(t *testing.T) {
	Convey("should print error'", t, func() {
		err := errors.New("my error")
		CurrentSystem().Error(err)
		So(false, ShouldEqual, false)
	})
}
*/
