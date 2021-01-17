package data

import (
	"context"
	"reflect"
	"strconv"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
)

func TestID(t *testing.T) {
	assert := assert.New(t)
	d := &Sample{}
	assert.Empty(d.ID)
	assert.True(d.TimeCreated().IsZero())
	assert.True(d.TimeUpdated().IsZero())
}

func TestBy(t *testing.T) {
	assert := assert.New(t)
	d := &Sample{}
	assert.Empty(d.GetBy())
	d.SetBy("user1")
	assert.Equal("user1", d.GetBy())
}

func TestMap(t *testing.T) {
	assert := assert.New(t)
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
	assert.Nil(err)
	sampleID := sample.ID
	sample2Obj, err := samplesG.Get(ctx, sampleID)
	assert.Nil(err)
	sample2 := sample2Obj.(*Sample)
	assert.NotNil(sample2.Map)
	assert.Equal("a", sample2.Map["1"])
	assert.Equal("b", sample2.Map["2"])
	assert.Equal("a", sample2.Array[0])
	assert.Equal("b", sample2.Array[1])
	assert.Equal(1, sample2.Numbers[0])
	assert.Equal(2, sample2.Numbers[1])
	assert.Equal("1", sample2.Obj.ID)
	assert.Equal("a", sample2.Obj.Name)
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
