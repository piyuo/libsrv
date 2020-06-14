package data

import (
	"reflect"
	"strconv"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type Sample struct {
	StoredObject
}

//NewSample create Sample instance
func NewSample() interface{} {
	return new(Sample)
}

//New create new object instance
func (do *Sample) New() interface{} {
	return new(Sample)
}

//Class represent name in database
func (do *Sample) ModelName() string {
	return "Sample"
}

func TestModelName(t *testing.T) {
	Convey("check object db name correct", t, func() {
		d := &Sample{}
		name := d.ModelName()
		So(name, ShouldEqual, "Sample")
	})
}

func TestNew(t *testing.T) {
	newObj := NewSample().(*Sample)
	newObj.SetID("1")
	Convey("new object should work", t, func() {
		So(newObj.ID(), ShouldEqual, "1")
	})

	newObj2 := newObj.New().(*Sample)
	Convey("create empty new object", t, func() {
		So(newObj2.ID(), ShouldEqual, "")
	})
}

func TestID(t *testing.T) {
	Convey("id should be empty", t, func() {
		d := &Sample{}
		id := d.ID()
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
	d := Sample{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		name := reflect.TypeOf(d).String()
		result = name
		//obj := reflect.New(reflect.TypeOf(d))
		//result = obj.Interface().(*Sample).Class()
	}
}

func BenchmarkNativeTypeSpeed(b *testing.B) {
	d := Sample{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		name := d.ModelName()
		result = name
	}
}

func BenchmarkReflectNewSpeed(b *testing.B) {
	d := Sample{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		obj := reflect.New(reflect.TypeOf(d)).Interface().(*Sample)
		result = obj.ModelName()
	}
}

func BenchmarkNativeNewSpeed(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		obj := new(Sample)
		result = obj.ModelName()
	}
}
