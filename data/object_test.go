package data

import (
	"reflect"
	"strconv"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type DataObjectChild struct {
	DBObject
}

//NewDataObjectChild create DataObjectChild instance
func NewDataObjectChild() interface{} {
	return new(DataObjectChild)
}

//New create new object instance
func (do *DataObjectChild) New() interface{} {
	return new(DataObjectChild)
}

//Class represent name in database
func (do *DataObjectChild) Class() string {
	return "DataObjectChild"
}

func TestClass(t *testing.T) {
	d := &DataObjectChild{}
	name := d.Class()
	Convey("check object db name correct", t, func() {
		So(name, ShouldEqual, "DataObjectChild")
	})
}

func TestNew(t *testing.T) {
	newObj := NewDataObjectChild().(*DataObjectChild)
	newObj.SetID("1")
	Convey("new object should work", t, func() {
		So(newObj.ID(), ShouldEqual, "1")
	})

	newObj2 := newObj.New().(*DataObjectChild)
	Convey("create empty new object", t, func() {
		So(newObj2.ID(), ShouldEqual, "")
	})
}

func TestID(t *testing.T) {
	d := &DataObjectChild{}
	id := d.ID()
	Convey("id should be empty", t, func() {
		So(id, ShouldBeEmpty)
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

func ThrowTimeout() error {
	return ErrOperationTimeout
}

func TestCustomError(t *testing.T) {
	err := ThrowTimeout()
	if err != ErrOperationTimeout {
		err = nil
	}
	Convey("compare custom error", t, func() {
		So(err, ShouldEqual, ErrOperationTimeout)
	})
}

var result string

func BenchmarkReflectTypeSpeed(b *testing.B) {
	d := DataObjectChild{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		name := reflect.TypeOf(d).String()
		result = name
		//obj := reflect.New(reflect.TypeOf(d))
		//result = obj.Interface().(*DataObjectChild).Class()
	}
}

func BenchmarkNativeTypeSpeed(b *testing.B) {
	d := DataObjectChild{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		name := d.Class()
		result = name
	}
}

func BenchmarkReflectNewSpeed(b *testing.B) {
	d := DataObjectChild{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		obj := reflect.New(reflect.TypeOf(d)).Interface().(*DataObjectChild)
		result = obj.Class()
	}
}

func BenchmarkNativeNewSpeed(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		obj := new(DataObjectChild)
		result = obj.Class()
	}
}
