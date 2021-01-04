package i18n

import (
	"strings"
	"time"

	"github.com/goodsign/monday"

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
//		So(DateToLocalStr(utcTime, "zh_TW"), ShouldEqual, "2021年1月2日")
//		So(DateToLocalStr(utcTime, "zh_CN"), ShouldEqual, "2021年1月2日")
//		So(DateToLocalStr(utcTime, "en_US"), ShouldEqual, "Jan 2, 2021")
//
func DateToLocalStr(t time.Time, locale string) string {
	switch locale {
	case "zh_TW":
		return t.Format("2006年1月2日")
	case "zh_CN":
		return t.Format("2006年1月2日")
	}
	return monday.Format(t, "Jan 2, 2006", "en_US")
}

// TimeToLocalStr convert utc time to local string
//
//		So(TimeToLocalStr(utcTime, "en_US"), ShouldEqual, "11:55 PM")
//		So(TimeToLocalStr(utcTime, "zh_TW"), ShouldEqual, "下午11:55")
//		So(TimeToLocalStr(utcTime, "zh_CN"), ShouldEqual, "下午11:55")
//
func TimeToLocalStr(t time.Time, locale string) string {
	switch locale {
	case "zh_TW":
		result := monday.Format(t, "PM3:04", "en_US")
		result = strings.Replace(result, "PM", "下午", 1)
		result = strings.Replace(result, "AM", "上午", 1)
		return result
	case "zh_CN":
		result := monday.Format(t, "PM3:04", "en_US")
		result = strings.Replace(result, "PM", "下午", 1)
		result = strings.Replace(result, "AM", "上午", 1)
		return result
	}
	return monday.Format(t, "3:04 PM", "en_US")
}

// DateTimeToLocalStr convert utc date time to local string
//
//		So(DateTimeToLocalStr(utcTime, "zh_TW"), ShouldEqual, "2021年1月2日 下午11:55")
//		So(DateTimeToLocalStr(utcTime, "zh_CN"), ShouldEqual, "2021年1月2日 下午11:55")
//		So(DateTimeToLocalStr(utcTime, "en_US"), ShouldEqual, "Jan 2, 2021 11:55 PM")
//
func DateTimeToLocalStr(t time.Time, locale string) string {
	return DateToLocalStr(t, locale) + " " + TimeToLocalStr(t, locale)
}
