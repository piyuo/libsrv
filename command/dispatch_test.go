package command

import (
	"context"
	"strconv"
	"testing"

	mock "github.com/piyuo/libsrv/command/mock"
	"github.com/piyuo/libsrv/command/types"
	"github.com/stretchr/testify/assert"
)

func TestEncodeDecodeCommand(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	act := &mock.CmdResponse{
		Value: "Hi",
	}
	dispatch := &Dispatch{
		Map: &mock.MapXXX{},
	}
	actBytes, err := dispatch.EncodeCommand(act.XXX_MapID(), act)
	actID, iAct2, err2 := dispatch.DecodeCommand(actBytes)
	assert.Nil(err2)
	act2 := iAct2.(*mock.CmdResponse)
	assert.Nil(err)
	assert.Equal(act.XXX_MapID(), actID)
	assert.Equal(act.Value, act2.Value)
}

func TestBetterResponseName(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	errOK := OK
	result := betterResponseName(errOK.XXX_MapID(), OK)
	assert.Equal("OK", result)

	err := Error("failed").(*types.Error)
	result = betterResponseName(err.XXX_MapID(), err)
	assert.Equal("failed", result)

	errText := &types.String{}
	result = betterResponseName(errText.XXX_MapID(), errText)
	assert.Equal("String", result)
}

func TestActionNoRespose(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	act := &mock.CmdNoRespond{
		Text: "Hi",
	}
	dispatch := &Dispatch{
		Map: &mock.MapXXX{},
	}
	//no response action,will cause &types.r{}
	actBytes, err := dispatch.EncodeCommand(act.XXX_MapID(), act)
	assert.Nil(err)
	resultBytes, err2 := dispatch.Route(context.Background(), actBytes)
	assert.NotNil(err2)
	assert.Nil(resultBytes)
}

func TestRoute(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	act := &mock.CmdRespond{
		Text: "Hi",
	}
	dispatch := &Dispatch{
		Map: &mock.MapXXX{},
	}
	actBytes, err := dispatch.EncodeCommand(act.XXX_MapID(), act)
	assert.Nil(err)

	resultBytes, err2 := dispatch.Route(context.Background(), actBytes)
	assert.Nil(err2)

	_, resp, err3 := dispatch.DecodeCommand(resultBytes)
	assert.Nil(err3)
	actualResponse := resp.(*mock.CmdResponse)
	assert.NotNil(actualResponse)
}

func TestHandle(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	//create sample data
	act := &mock.CmdRespond{
		Text: "Hi",
	}
	//create dispatch and register
	dispatch := &Dispatch{
		Map: &mock.MapXXX{},
	}
	_, respInterface, err := dispatch.runAction(context.Background(), act)
	response := respInterface.(*mock.CmdResponse)
	assert.Nil(err)
	assert.NotNil(response)
}

var benchmarkResult string

func BenchmarkStringMapSpeed(b *testing.B) {
	list := make(map[string]string, 100)
	for x := 0; x < 100; x++ {
		list[strconv.Itoa(x)] = strconv.Itoa(x)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for x := 0; x < 100; x++ {
			benchmarkResult = list[strconv.Itoa(x)]
		}
	}
}

func BenchmarkIntMapSpeed(b *testing.B) {
	list := make(map[int]string, 100)
	for x := 0; x < 100; x++ {
		list[x] = strconv.Itoa(x)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for x := 0; x < 100; x++ {
			benchmarkResult = list[x]
		}
	}
}

var tmp []byte

func BenchmarkAppend(b *testing.B) {
	bytes1 := []byte("my first slice")
	bytes2 := []byte("second slice")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for x := 0; x < 100; x++ {
			tmp = append(bytes1[:], bytes2[:]...)
		}
	}
}

func BenchmarkCopyPreAllocate(b *testing.B) {
	bytes1 := []byte("my first slice")
	bytes2 := []byte("second slice")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for x := 0; x < 100; x++ {
			totalLen := len(bytes1) + len(bytes2)
			tmp := make([]byte, totalLen)
			var i int
			i += copy(tmp, bytes1)
			i += copy(tmp[i:], bytes2)
		}
	}
}

func BenchmarkDispatch(b *testing.B) {
	act := &mock.CmdRespond{
		Text: "Hi",
	}
	dispatch := &Dispatch{
		Map: &mock.MapXXX{},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		actBytes, _ := dispatch.EncodeCommand(act.XXX_MapID(), act)
		resultBytes, _ := dispatch.Route(context.Background(), actBytes)
		_, _, _ = dispatch.DecodeCommand(resultBytes)
	}
}
