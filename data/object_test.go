package data

import (
	"reflect"
	"strconv"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestID(t *testing.T) {
	Convey("id should be empty", t, func() {
		d := &Sample{}
		So(d.ID, ShouldBeEmpty)
		So(d.CreateTime.IsZero(), ShouldBeTrue)
		So(d.UpdateTime.IsZero(), ShouldBeTrue)
	})
}

func doWork(f func(i int) string) string {
	return f(1)
}

func TestFunctionCallback(t *testing.T) {
	text := doWork(func(i int) string {
		return strconv.Itoa(i)
	})
	Convey("doWork return work", t, func() {
		So(text, ShouldEqual, "1")
	})
}

var result string

func BenchmarkReflectTypeSpeed(b *testing.B) {
	d := Sample{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		name := reflect.TypeOf(d).String()
		result = name
		//obj := reflect.New(reflect.TypeOf(d))
		//result = obj.Interface().(*Sample).Class()
	}
}
