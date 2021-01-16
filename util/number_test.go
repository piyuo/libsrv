package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldConvertInterfaceToInt(t *testing.T) {
	assert := assert.New(t)
	num, err := ToInt(2)
	assert.Nil(err)
	assert.Equal(2, num)
	num, err = ToInt(float64(2))
	assert.Nil(err)
	assert.Equal(2, num)
	num, err = ToInt(float32(2))
	assert.Nil(err)
	assert.Equal(2, num)
	num, err = ToInt(int64(2))
	assert.Nil(err)
	assert.Equal(2, num)
	num, err = ToInt(int32(2))
	assert.Nil(err)
	assert.Equal(2, num)
	num, err = ToInt(uint64(2))
	assert.Nil(err)

	assert.Equal(2, num)
	num, err = ToInt(uint32(2))
	assert.Nil(err)

	assert.Equal(2, num)
	num, err = ToInt(uint(2))
	assert.Nil(err)

	assert.Equal(2, num)
	num, err = ToInt("err")
	assert.NotNil(err)
	assert.Equal(0, num)
}

func TestShouldConvertInterfaceToFloat64(t *testing.T) {
	assert := assert.New(t)
	num, err := ToFloat64(2)
	assert.Nil(err)
	assert.Equal(float64(2), num)

	num, err = ToFloat64(float64(2))
	assert.Nil(err)
	assert.Equal(float64(2), num)

	num, err = ToFloat64(float32(2))
	assert.Nil(err)
	assert.Equal(float64(2), num)

	num, err = ToFloat64(int64(2))
	assert.Nil(err)
	assert.Equal(float64(2), num)

	num, err = ToFloat64(int32(2))
	assert.Nil(err)
	assert.Equal(float64(2), num)

	num, err = ToFloat64(uint64(2))
	assert.Nil(err)
	assert.Equal(float64(2), num)

	num, err = ToFloat64(uint32(2))
	assert.Nil(err)
	assert.Equal(float64(2), num)

	num, err = ToFloat64(uint(2))
	assert.Nil(err)
	assert.Equal(float64(2), num)

	num, err = ToFloat64("err")
	assert.NotNil(err)
	assert.Equal(float64(0), num)
}

func TestShouldConvertInterfaceToUint32(t *testing.T) {
	assert := assert.New(t)
	num, err := ToUint32(2)
	assert.Nil(err)
	assert.Equal(uint32(2), num)

	num, err = ToUint32(float64(2))
	assert.Nil(err)
	assert.Equal(uint32(2), num)

	num, err = ToUint32(float32(2))
	assert.Nil(err)
	assert.Equal(uint32(2), num)

	num, err = ToUint32(int64(2))
	assert.Nil(err)
	assert.Equal(uint32(2), num)

	num, err = ToUint32(int32(2))
	assert.Nil(err)
	assert.Equal(uint32(2), num)

	num, err = ToUint32(uint64(2))
	assert.Nil(err)
	assert.Equal(uint32(2), num)

	num, err = ToUint32(uint32(2))
	assert.Nil(err)
	assert.Equal(uint32(2), num)

	num, err = ToUint32(uint(2))
	assert.Nil(err)
	assert.Equal(uint32(2), num)

	num, err = ToUint32("err")
	assert.NotNil(err)
	assert.Equal(uint32(0), num)
}
