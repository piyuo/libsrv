package command

import (
	"context"
	"os"
	"strconv"
	"testing"

	mock "github.com/piyuo/libsrv/command/mock"
	shared "github.com/piyuo/libsrv/command/shared"
	"github.com/stretchr/testify/assert"
)

func TestEncodeDecodeCommand(t *testing.T) {
	assert := assert.New(t)
	act := &mock.RespondAction{
		Text: "Hi",
	}
	dispatch := &Dispatch{
		Map: &mock.MapXXX{},
	}
	actBytes, err := dispatch.encodeCommand(act.XXX_MapID(), act)
	actID, iAct2, err2 := dispatch.decodeCommand(actBytes)
	assert.Nil(err2)
	act2 := iAct2.(*mock.RespondAction)
	assert.Nil(err)
	assert.Equal(act.XXX_MapID(), actID)
	assert.Equal(act.Text, act2.Text)
}

func TestBetterResponseName(t *testing.T) {
	assert := assert.New(t)
	errOK := OK().(*shared.PbOK)
	result := betterResponseName(errOK.XXX_MapID(), errOK)
	assert.Equal("OK", result)

	err := Error("failed").(*shared.PbError)
	result = betterResponseName(err.XXX_MapID(), err)
	assert.Equal("failed", result)

	errText := &shared.PbString{}
	result = betterResponseName(errText.XXX_MapID(), errText)
	assert.Equal("PbString", result)
}

func TestActionNoRespose(t *testing.T) {
	assert := assert.New(t)
	act := &mock.NoRespondAction{
		Text: "Hi",
	}
	dispatch := &Dispatch{
		Map: &mock.MapXXX{},
	}
	//no response action,will cause &shared.Err{}
	actBytes, err := dispatch.encodeCommand(act.XXX_MapID(), act)
	assert.Nil(err)
	resultBytes, err2 := dispatch.Route(context.Background(), actBytes)
	assert.NotNil(err2)
	assert.Nil(resultBytes)
}

func TestRoute(t *testing.T) {
	assert := assert.New(t)
	act := &mock.RespondAction{
		Text: "Hi",
	}
	dispatch := &Dispatch{
		Map: &mock.MapXXX{},
	}
	actBytes, err := dispatch.encodeCommand(act.XXX_MapID(), act)
	resultBytes, err2 := dispatch.Route(context.Background(), actBytes)
	_, resp, err3 := dispatch.decodeCommand(resultBytes)
	actualResponse := resp.(*shared.PbOK)
	assert.Nil(err)
	assert.Nil(err2)
	assert.Nil(err3)
	assert.NotNil(actualResponse)
}

func TestHandle(t *testing.T) {
	assert := assert.New(t)

	//create sample data
	act := &mock.RespondAction{
		Text: "Hi",
	}
	//create dispatch and register
	dispatch := &Dispatch{
		Map: &mock.MapXXX{},
	}
	_, respInterface, err := dispatch.runAction(context.Background(), act)
	response := respInterface.(*shared.PbOK)
	assert.Nil(err)
	assert.NotNil(response)
}

func TestTimeExecuteAction(t *testing.T) {
	assert := assert.New(t)
	os.Setenv("PIYUO_SLOW", "1")
	act := &mock.SlowAction{}
	dispatch := &Dispatch{
		Map: &mock.MapXXX{},
	}
	_, respInterface, err := dispatch.timeExecuteAction(context.Background(), act)
	assert.Nil(err)
	assert.NotNil(respInterface)
}

var benchmarkResult string

func BenchmarkStringMapSpeed(b *testing.B) {
	var list map[string]string
	list = make(map[string]string, 100)
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
	var list map[int]string
	list = make(map[int]string, 100)
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
	act := &mock.RespondAction{
		Text: "Hi",
	}
	dispatch := &Dispatch{
		Map: &mock.MapXXX{},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		actBytes, _ := dispatch.encodeCommand(act.XXX_MapID(), act)
		resultBytes, _ := dispatch.Route(context.Background(), actBytes)
		_, _, _ = dispatch.decodeCommand(resultBytes)
	}
}
