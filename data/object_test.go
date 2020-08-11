package data

import (
	"context"
	"reflect"
	"strconv"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestID(t *testing.T) {
	Convey("id should be empty", t, func() {
		d := &Sample{}
		So(d.ID, ShouldBeEmpty)
		So(d.TimeCreated().IsZero(), ShouldBeTrue)
		So(d.TimeUpdated().IsZero(), ShouldBeTrue)
	})
}

func TestMap(t *testing.T) {
	Convey("should set/get map", t, func() {
		ctx := context.Background()
		dbG, dbR := createSampleDB()
		defer removeSampleDB(dbG, dbR)
		samplesG, samplesR := createSampleTable(dbG, dbR)
		defer removeSampleTable(samplesG, samplesR)

		sample := &Sample{
			Name:  "sample",
			Value: 1,
			Map: map[string]string{
				"1": "a",
				"2": "b",
			},
			Array: []string{
				"a",
				"b",
			},
			Numbers: []int{
				1,
				2,
			},
		}

		sample.Obj = &PlainObject{
			ID:   "1",
			Name: "a",
		}

		err := samplesG.Set(ctx, sample)
		So(err, ShouldBeNil)
		sampleID := sample.ID
		sample2Obj, err := samplesG.Get(ctx, sampleID)
		So(err, ShouldBeNil)
		sample2 := sample2Obj.(*Sample)
		So(sample2.Map, ShouldNotBeNil)
		So(sample2.Map["1"], ShouldEqual, "a")
		So(sample2.Map["2"], ShouldEqual, "b")
		So(sample2.Array[0], ShouldEqual, "a")
		So(sample2.Array[1], ShouldEqual, "b")
		So(sample2.Numbers[0], ShouldEqual, 1)
		So(sample2.Numbers[1], ShouldEqual, 2)
		So(sample2.Obj.ID, ShouldEqual, "1")
		So(sample2.Obj.Name, ShouldEqual, "a")
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
