package data

import (
	"context"
	"reflect"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestID(t *testing.T) {
	assert := assert.New(t)
	d := &Sample{}
	assert.Empty(d.ID)
	assert.True(d.GetCreateTime().IsZero())
	assert.True(d.GetUpdateTime().IsZero())
}

func TestUserID(t *testing.T) {
	assert := assert.New(t)
	d := &Sample{}
	assert.Empty(d.GetUserID())
	d.SetUserID("user1")
	assert.Equal("user1", d.GetUserID())
}

func TestAccountID(t *testing.T) {
	assert := assert.New(t)
	d := &Sample{}
	assert.Empty(d.GetAccountID())
	d.SetAccountID("account1")
	assert.Equal("account1", d.GetAccountID())
}

func TestMap(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	g, err := NewSampleGlobalDB(ctx)
	assert.Nil(err)
	defer g.Close()
	table := g.SampleTable()

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

	err = table.Set(ctx, sample)
	defer table.DeleteObject(ctx, sample)
	assert.Nil(err)
	sampleID := sample.ID
	sample2Obj, err := table.Get(ctx, sampleID)
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
	assert := assert.New(t)
	text := doWork(func(i int) string {
		return strconv.Itoa(i)
	})
	assert.Equal("1", text)
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
