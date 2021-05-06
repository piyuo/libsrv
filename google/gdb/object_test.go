package gdb

import (
	"context"
	"reflect"
	"strconv"
	"testing"

	"github.com/piyuo/libsrv/identifier"
	"github.com/stretchr/testify/assert"
)

func TestObjectID(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	d := &Sample{}
	assert.Empty(d.ID())
	assert.True(d.CreateTime().IsZero())
	assert.True(d.UpdateTime().IsZero())
}

func TestObjectNilSafety(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	ctx := context.Background()
	client := sampleClient()
	sample := &Sample{}
	err := client.Set(ctx, sample)
	assert.Nil(err)
	defer client.Delete(ctx, sample)

	sample2Obj, err := client.Get(ctx, &Sample{}, sample.ID())
	assert.Nil(err)
	sample2 := sample2Obj.(*Sample)

	assert.NotNil(sample2.Array)
	assert.NotNil(sample2.Numbers)
	assert.NotNil(sample2.PObj)
}

func TestObjectTime(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	ctx := context.Background()
	client := sampleClient()
	sample := &Sample{}
	err := client.Set(ctx, sample)
	assert.Nil(err)
	defer client.Delete(ctx, sample)

	sample2Obj, err := client.Get(ctx, &Sample{}, sample.ID())
	assert.Nil(err)
	sample2 := sample2Obj.(*Sample)
	assert.False(sample2.CreateTime().IsZero())
	assert.False(sample2.UpdateTime().IsZero())
}

func TestObjectUserID(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	d := &Sample{}
	assert.Empty(d.UserID())
	d.SetUserID("user1")
	assert.Equal("user1", d.UserID())
}

func TestObjectAccountID(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	d := &Sample{}
	assert.Empty(d.AccountID())
	d.SetAccountID("account1")
	assert.Equal("account1", d.AccountID())
}

func TestObjectMap(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	ctx := context.Background()
	client := sampleClient()
	rand := identifier.RandomString(8)
	name := "test-object-map-" + rand

	sample := &Sample{
		Name:  name,
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

	sample.PObj = &PlainObject{
		ID:   "1",
		Name: "a",
	}

	err := client.Set(ctx, sample)
	assert.Nil(err)
	defer client.Delete(ctx, sample)

	sample2Obj, err := client.Get(ctx, &Sample{}, sample.ID())
	assert.Nil(err)
	sample2 := sample2Obj.(*Sample)
	assert.NotNil(sample2.Map)
	assert.Equal("a", sample2.Map["1"])
	assert.Equal("b", sample2.Map["2"])
	assert.Equal("a", sample2.Array[0])
	assert.Equal("b", sample2.Array[1])
	assert.Equal(1, sample2.Numbers[0])
	assert.Equal(2, sample2.Numbers[1])
	assert.Equal("1", sample2.PObj.ID)
	assert.Equal("a", sample2.PObj.Name)
}

func TestClientObjectWithoutFactory(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	ctx := context.Background()
	client := sampleClient()
	sampleNoFactory := &SampleNoFactory{}
	err := client.Set(ctx, sampleNoFactory)
	assert.NotNil(err)
}

func doWork(f func(i int) string) string {
	return f(1)
}

func TestFunctionCallback(t *testing.T) {
	t.Parallel()
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
