package command

import (
	"testing"

	"github.com/piyuo/libsrv/src/command/shared"
	"github.com/stretchr/testify/assert"
)

func TestShared(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	//should create text response
	text := String("hi").(*shared.PbString)
	assert.Equal("hi", text.Value)

	//should create number response
	num := Int(201).(*shared.PbInt)
	assert.Equal(int32(201), num.Value)

	//should create bool response
	b := Bool(true).(*shared.PbBool)
	assert.True(b.Value)
	b = Bool(false).(*shared.PbBool)
	assert.False(b.Value)

	//should create error response
	err := Error("errCode").(*shared.PbError)
	assert.Equal("errCode", err.Code)
	//should be OK
	ok := OK()
	assert.True(IsOK(ok))

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

func TestPbString(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	assert.False(IsString(nil, ""))
	assert.False(IsString(String("hi"), ""))
	assert.True(IsString(String("hi"), "hi"))
}

func TestPbInt(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	assert.False(IsInt(nil, 1))
	assert.False(IsInt(Int(12), 42))
	assert.True(IsInt(Int(42), 42))
}

func TestPbBool(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	assert.False(IsBool(nil, false))
	assert.False(IsBool(Bool(false), true))
	assert.True(IsBool(Bool(true), true))
}
