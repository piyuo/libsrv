package i18n

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUtcTimestamp(t *testing.T) {
	assert := assert.New(t)
	tt, err := UtcTimestamp(time.Now())
	assert.Nil(err)
	assert.NotEmpty(tt)
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
	assert.Equal("2021年1月2日", DateToLocalStr(utcTime, "zh_TW"))
	assert.Equal("2021年1月2日", DateToLocalStr(utcTime, "zh_CN"))
	assert.Equal("Jan 2, 2021", DateToLocalStr(utcTime, "en_US"))
}

func TestTimeToLocalStr(t *testing.T) {
	assert := assert.New(t)
	utcTime := time.Date(2021, time.January, 2, 23, 55, 0, 0, time.UTC)
	assert.Equal("11:55 PM", TimeToLocalStr(utcTime, "en_US"))
	assert.Equal("下午11:55", TimeToLocalStr(utcTime, "zh_TW"))
	assert.Equal("下午11:55", TimeToLocalStr(utcTime, "zh_CN"))
}

func TestDateTimeToLocalStr(t *testing.T) {
	assert := assert.New(t)
	utcTime := time.Date(2021, time.January, 2, 23, 55, 0, 0, time.UTC)
	assert.Equal("2021年1月2日 下午11:55", DateTimeToLocalStr(utcTime, "zh_TW"))
	assert.Equal("2021年1月2日 下午11:55", DateTimeToLocalStr(utcTime, "zh_CN"))
	assert.Equal("Jan 2, 2021 11:55 PM", DateTimeToLocalStr(utcTime, "en_US"))
}
