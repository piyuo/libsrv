package command

import (
	"testing"

	"github.com/piyuo/libsrv/command/simple"
	"github.com/stretchr/testify/assert"
)

func TestHelper(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	//should create text response
	text := String("hi").(*simple.String)
	assert.Equal("hi", text.Value)

	//should create number response
	num := Int(201).(*simple.Number)
	assert.Equal(int32(201), num.Value)

	//should create bool response
	b := Bool(true).(*simple.Bool)
	assert.True(b.Value)
	b = Bool(false).(*simple.Bool)
	assert.False(b.Value)

	//should create error response
	err := Error("errCode").(*simple.Error)
	assert.Equal("errCode", err.Code)
	//should be OK
	assert.True(IsOK(OK))

	//should not be OK
	assert.False(IsOK(1))

	//should be INVALID error
	err3 := Error("INVALID")
	assert.True(IsError(err3, "INVALID"))

	//should not be INVALID error
	assert.False(IsError(nil, "INVALID"))
	err2 := 3
	assert.False(IsError(err2, "INVALID"))
}

func TestGetErrorCode(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	err := Error("errCode").(*simple.Error)
	assert.Equal("errCode", GetErrorCode(err))
	assert.Equal("", GetErrorCode("notError"))
	assert.Equal("", GetErrorCode(nil))
}

func TestString(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	assert.False(IsString(nil, ""))
	assert.False(IsString(String("hi"), ""))
	assert.True(IsString(String("hi"), "hi"))
}

func TestInt(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	assert.False(IsInt(nil, 1))
	assert.False(IsInt(Int(12), 42))
	assert.True(IsInt(Int(42), 42))
}

func TestBool(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	assert.False(IsBool(nil, false))
	assert.False(IsBool(Bool(false), true))
	assert.True(IsBool(Bool(true), true))
}

func TestBlockShort(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	assert.True(IsBlockShort(BlockShort))
	assert.False(IsBlockShort(Error("otherError")))
}

func TestBlockLong(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	assert.True(IsBlockLong(BlockLong))
	assert.False(IsBlockLong(Error("otherError")))
}
