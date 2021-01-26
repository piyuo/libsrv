package i18n

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUtc(t *testing.T) {
	assert := assert.New(t)
	now := time.Now()
	zone, offset := now.Zone()

	utcTime := now.UTC()
	utcAgain := utcTime.UTC()
	locTime := utcTime.In(time.FixedZone(zone, offset))
	locAgain := utcAgain.In(time.FixedZone(zone, offset))
	locThird := locAgain.In(time.FixedZone(zone, offset))

	nowStr := DateTimeToStr(now, "en_US")
	locStr := DateTimeToStr(locTime, "en_US")
	againStr := DateTimeToStr(locAgain, "en_US")
	thirdStr := DateTimeToStr(locThird, "en_US")

	//there is no probelm to do UTC <-> local again and again
	assert.Equal(nowStr, locStr)
	assert.Equal(nowStr, againStr)
	assert.Equal(nowStr, thirdStr)
}

func TestUtcTimestamp(t *testing.T) {
	assert := assert.New(t)
	now := time.Now()
	zone, offset := now.Zone()

	timestamp, err := ToUtcTimestamp(now)
	assert.Nil(err)
	assert.NotNil(timestamp)

	locTime := FromUtcTimestamp(timestamp, zone, offset)
	nowStr := DateTimeToStr(now, "en_US")
	locStr := DateTimeToStr(locTime, "en_US")
	assert.Equal(nowStr, locStr)

	//nil timestamp will result empty time
	locTime = FromUtcTimestamp(nil, zone, offset)
	assert.True(locTime.IsZero())
}

func TestUtcToLocal(t *testing.T) {
	assert := assert.New(t)
	utcTime := time.Date(2021, time.January, 2, 23, 55, 0, 0, time.UTC)
	locTime := UtcToLocal(utcTime, "PDT", -25200)
	utcTime2 := locTime.In(time.UTC)
	assert.Equal(2021, utcTime2.Year())
	assert.Equal(time.January, utcTime2.Month())
	assert.Equal(2, utcTime2.Day())
	assert.Equal(23, utcTime2.Hour())
	assert.Equal(55, utcTime2.Minute())
	utcTime3 := LocalToUtc(locTime)
	assert.Equal(2021, utcTime3.Year())
	assert.Equal(time.January, utcTime3.Month())
	assert.Equal(2, utcTime3.Day())
	assert.Equal(23, utcTime3.Hour())
	assert.Equal(55, utcTime3.Minute())
}

func TestDateToLocalStr(t *testing.T) {
	assert := assert.New(t)
	utcTime := time.Date(2021, time.January, 2, 23, 55, 0, 0, time.UTC)
	assert.Equal("2021年1月2日", DateToStr(utcTime, "zh_TW"))
	assert.Equal("2021年1月2日", DateToStr(utcTime, "zh_CN"))
	assert.Equal("Jan 2, 2021", DateToStr(utcTime, "en_US"))
}

func TestTimeToLocalStr(t *testing.T) {
	assert := assert.New(t)
	utcTime := time.Date(2021, time.January, 2, 23, 55, 0, 0, time.UTC)
	assert.Equal("11:55 PM", TimeToStr(utcTime, "en_US"))
	assert.Equal("下午11:55", TimeToStr(utcTime, "zh_TW"))
	assert.Equal("下午11:55", TimeToStr(utcTime, "zh_CN"))
}

func TestDateTimeToLocalStr(t *testing.T) {
	assert := assert.New(t)
	utcTime := time.Date(2021, time.January, 2, 23, 55, 0, 0, time.UTC)
	assert.Equal("2021年1月2日 下午11:55", DateTimeToStr(utcTime, "zh_TW"))
	assert.Equal("2021年1月2日 下午11:55", DateTimeToStr(utcTime, "zh_CN"))
	assert.Equal("Jan 2, 2021 11:55 PM", DateTimeToStr(utcTime, "en_US"))
}
