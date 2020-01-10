package libsrv

import (
	"strings"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSystem(t *testing.T) {

	Convey("should has only one instance and pass check'", t, func() {
		s1 := CurrentSystem()
		s2 := CurrentSystem()
		So(s1.IsProduction(), ShouldEqual, s2.IsProduction())
		CurrentSystem().Check()
		So(s1.IsProduction(), ShouldEqual, false)
	})

	Convey("should able join dir and current dir'", t, func() {
		text := CurrentSystem().JoinCurrentDir("../../")
		So(strings.HasSuffix(text, "/go"), ShouldEqual, true)
	})
}

func TestTimer(t *testing.T) {
	Convey("should use timer'", t, func() {
		CurrentSystem().TimerStart()
		time.Sleep(1 * time.Millisecond)
		ms := CurrentSystem().TimerStop()
		So(ms >= 1, ShouldBeTrue)
	})
}

func TestGetID(t *testing.T) {
	Convey("should get id'", t, func() {
		id := CurrentSystem().ID()
		So(id, ShouldEqual, "dev")
	})
}

func TestCheckProduction(t *testing.T) {
	s := &system{}
	Convey("should set production correctly", t, func() {
		So(s.checkProduction("dev"), ShouldEqual, false)
		So(s.checkProduction("piyuo-tw-a-app"), ShouldEqual, false)
		So(s.checkProduction("piyuo-tw-b-app"), ShouldEqual, false)
		So(s.checkProduction("piyuo-tw-m-app"), ShouldEqual, true)
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
		CurrentSystem().Info("hello log")
		So(true, ShouldEqual, true)
	})
}

/*
//TestLog is a production test, it will write log to google cloud platform under log viewer "Google Project, project name"
func TestLog(t *testing.T) {
	Convey("should print'", t, func() {
		CurrentSystem().Notice("my notice log")
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
func TestErrorManually(t *testing.T) {
	Convey("should print error manually'", t, func() {
		stack := "at firstLine (a.js:3)\nat secondLine (b.js:3)"
		CurrentSystem().ErrorManually("error manually", stack)
		So(false, ShouldEqual, false)
	})
}
*/
