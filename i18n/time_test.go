package i18n

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUtcTimestamp(t *testing.T) {
	Convey("should create utc timestamp", t, func() {
		t, err := UtcTimestamp(time.Now())
		So(err, ShouldBeNil)
		So(t, ShouldNotBeEmpty)
	})
}

func TestUtcToLocal(t *testing.T) {
	Convey("should create local time", t, func() {
		utcTime := time.Date(2021, time.January, 2, 23, 55, 0, 0, time.UTC)
		locTime := UtcToLocal(utcTime, "PDT", -25200)
		utcTime2 := locTime.In(time.UTC)
		So(utcTime2.Year(), ShouldEqual, 2021)
		So(utcTime2.Month(), ShouldEqual, time.January)
		So(utcTime2.Day(), ShouldEqual, 2)
		So(utcTime2.Hour(), ShouldEqual, 23)
		So(utcTime2.Minute(), ShouldEqual, 55)
		utcTime3 := LocalToUtc(locTime)
		So(utcTime3.Year(), ShouldEqual, 2021)
		So(utcTime3.Month(), ShouldEqual, time.January)
		So(utcTime3.Day(), ShouldEqual, 2)
		So(utcTime3.Hour(), ShouldEqual, 23)
		So(utcTime3.Minute(), ShouldEqual, 55)
	})
}

func TestDateToLocalStr(t *testing.T) {
	Convey("should create local time string", t, func() {
		utcTime := time.Date(2021, time.January, 2, 23, 55, 0, 0, time.UTC)
		locStr := DateToLocalStr(utcTime, "PDT", -25200)
		So(locStr, ShouldEqual, "2021-01-02")
	})
}

func TestDateTimeToLocalStr(t *testing.T) {
	Convey("should create local time string", t, func() {
		utcTime := time.Date(2021, time.January, 2, 23, 55, 0, 0, time.UTC)
		locStr := DateTimeToLocalStr(utcTime, "PDT", -25200)
		So(locStr, ShouldEqual, "2021-01-02 16:55") //PDT Time
	})
}

func TestTimeToLocalStr(t *testing.T) {
	Convey("should create local time string", t, func() {
		utcTime := time.Date(2021, time.January, 2, 23, 55, 0, 0, time.UTC)
		locStr := TimeToLocalStr(utcTime, "PDT", -25200)
		So(locStr, ShouldEqual, "16:55") //PDT Time
	})
}
