package command

import (
	"strconv"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestEncodeDecodeCommand(t *testing.T) {
	act := &TestAction{
		Text: "Hi",
	}
	dispatch := &Dispatch{
		Map: &TestMap{},
	}
	actBytes, err := dispatch.encodeCommand(act.XXX_MapID(), act)
	actID, iAct2, err2 := dispatch.decodeCommand(actBytes)
	act2 := iAct2.(*TestAction)
	Convey("test decode command is right", t, func() {
		So(err, ShouldBeNil)
		So(err2, ShouldBeNil)
		So(actID, ShouldEqual, act.XXX_MapID())
		So(act2.Text, ShouldEqual, act.Text)
	})
}

func TestActionNoRespond(t *testing.T) {
	act := &TestActionNotRespond{
		Text: "Hi",
	}
	dispatch := &Dispatch{
		Map: &TestMap{},
	}
	actBytes, err := dispatch.encodeCommand(act.XXX_MapID(), act)
	resultBytes, err2 := dispatch.Route(actBytes)

	Convey("test dispatch route", t, func() {
		So(err, ShouldBeNil)
		So(err2, ShouldBeNil)
		So(resultBytes, ShouldBeNil)
	})
}

func TestRoute(t *testing.T) {
	act := &TestAction{
		Text: "Hi",
	}
	dispatch := &Dispatch{
		Map: &TestMap{},
	}
	actBytes, err := dispatch.encodeCommand(act.XXX_MapID(), act)
	resultBytes, err2 := dispatch.Route(actBytes)
	respID, resp, err3 := dispatch.decodeCommand(resultBytes)
	actualResponse := resp.(*TestResponse)
	stringResp := &TestResponse{}
	Convey("test dispatch route", t, func() {
		So(err, ShouldBeNil)
		So(err2, ShouldBeNil)
		So(err3, ShouldBeNil)
		So(respID, ShouldEqual, stringResp.XXX_MapID())
		So(actualResponse.Text, ShouldEqual, "Hi")
	})
}

func TestHandle(t *testing.T) {

	//create sample data
	act := &TestAction{
		Text: "Hi",
	}
	//create dispatch and register
	dispatch := &Dispatch{
		Map: &TestMap{},
	}

	//test dispatch route
	respID, respInterface, err := dispatch.handle(act)
	response := respInterface.(*TestResponse)
	stringResponse := &TestResponse{}
	Convey("test despatch handle", t, func() {
		So(err, ShouldBeNil)
		So(respID, ShouldEqual, stringResponse.XXX_MapID())
		So(response.Text, ShouldEqual, "Hi")
	})
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
	act := &TestAction{
		Text: "Hi",
	}
	dispatch := &Dispatch{
		Map: &TestMap{},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		actBytes, _ := dispatch.encodeCommand(act.XXX_MapID(), act)
		resultBytes, _ := dispatch.Route(actBytes)
		_, _, _ = dispatch.decodeCommand(resultBytes)
	}
}
