package i18n

import (
	"strings"
	"time"

	"github.com/goodsign/monday"

	"github.com/golang/protobuf/ptypes"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ToUtcTimestamp create utc timestamp from time
//
func ToUtcTimestamp(t time.Time) (*timestamppb.Timestamp, error) {
	utcTime := t.UTC()
	return ptypes.TimestampProto(utcTime)
}

// FromUtcTimestamp create time form utc timestamp
//
func FromUtcTimestamp(t *timestamppb.Timestamp, zone string, offset int) time.Time {
	if t == nil {
		return time.Time{}
	}
	loc := time.FixedZone(zone, offset)
	utcTime := t.AsTime()
	return utcTime.In(loc)
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

// DateToStr convert utc date to local string
//
//		So(DateToStr(utcTime, "zh_TW"), ShouldEqual, "2021年1月2日")
//		So(DateToStr(utcTime, "zh_CN"), ShouldEqual, "2021年1月2日")
//		So(DateToStr(utcTime, "en_US"), ShouldEqual, "Jan 2, 2021")
//
func DateToStr(t time.Time, locale string) string {
	switch locale {
	case "zh_TW":
		return t.Format("2006年1月2日")
	case "zh_CN":
		return t.Format("2006年1月2日")
	}
	return monday.Format(t, "Jan 2, 2006", "en_US")
}

// TimeToStr convert utc time to local string
//
//		So(TimeToStr(utcTime, "en_US"), ShouldEqual, "11:55 PM")
//		So(TimeToStr(utcTime, "zh_TW"), ShouldEqual, "下午11:55")
//		So(TimeToStr(utcTime, "zh_CN"), ShouldEqual, "下午11:55")
//
func TimeToStr(t time.Time, locale string) string {
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

// DateTimeToStr convert utc date time to local string
//
//		So(DateTimeToStr(utcTime, "zh_TW"), ShouldEqual, "2021年1月2日 下午11:55")
//		So(DateTimeToStr(utcTime, "zh_CN"), ShouldEqual, "2021年1月2日 下午11:55")
//		So(DateTimeToStr(utcTime, "en_US"), ShouldEqual, "Jan 2, 2021 11:55 PM")
//
func DateTimeToStr(t time.Time, locale string) string {
	return DateToStr(t, locale) + " " + TimeToStr(t, locale)
}
