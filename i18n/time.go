package i18n

import (
	"time"

	"github.com/golang/protobuf/ptypes"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// UtcTimestamp create utc timestamp
//
func UtcTimestamp(t time.Time) (*timestamppb.Timestamp, error) {
	utcTime := t.UTC()
	return ptypes.TimestampProto(utcTime)
}

// UtcToLocal convert utc time to local time
//
//     var now = time.Now();
//     var local = UtcToLocal(now,"PDT", -25200);
//
func UtcToLocal(t time.Time, zone string, offset int) time.Time {
	loc := time.FixedZone(zone, offset)
	return t.In(loc)
}

// LocalToUtc convert local time to utc time
//
//     var now = time.Now();
//     var utc = UtcToLocal(now);
//
func LocalToUtc(t time.Time) time.Time {
	return t.In(time.UTC)
}

// DateToLocalStr convert utc date to local string
//
//	locStr := DateToLocalStr(utcTime, "PDT", -25200)
//		So(locStr, ShouldEqual, "2021-01-02")
//	})
//
func DateToLocalStr(t time.Time, zone string, offset int) string {
	local := UtcToLocal(t, zone, offset)
	layout := "2006-01-02"
	return local.Format(layout)
}

// DateTimeToLocalStr convert utc date time to local string
//
//		locStr := DateTimeToLocalStr(utcTime, "PDT", -25200)
//		So(locStr, ShouldEqual, "2021-01-02 16:55") //PDT Time
//
func DateTimeToLocalStr(t time.Time, zone string, offset int) string {
	local := UtcToLocal(t, zone, offset)
	layout := "2006-01-02 15:04"
	return local.Format(layout)
}

// TimeToLocalStr convert utc time to local string
//
//	locStr := TimeToLocalStr(utcTime, "PDT", -25200)
//	So(locStr, ShouldEqual, "16:55") //PDT Time
//
func TimeToLocalStr(t time.Time, zone string, offset int) string {
	local := UtcToLocal(t, zone, offset)
	layout := "15:04"
	return local.Format(layout)
}
