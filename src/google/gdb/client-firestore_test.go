package gdb

import (
	"context"
	"strconv"
	"testing"

	"github.com/piyuo/libsrv/src/google/gaccount"
	"github.com/piyuo/libsrv/src/identifier"
	"github.com/piyuo/libsrv/src/util"
	"github.com/stretchr/testify/assert"
)

func TestGdbObjectWithoutFactory(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	ctx := context.Background()
	client := sampleClient()
	sampleNoFactory := &SampleNoFactory{}
	err := client.Set(ctx, sampleNoFactory)
	assert.NotNil(err)
}

func TestGdbClientClose(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	ctx := context.Background()
	cred, err := gaccount.GlobalCredential(ctx)
	assert.Nil(err)
	client, err := NewClient(ctx, cred)
	assert.Nil(err)
	client.Close()
}

func TestGdbClientCRUD(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	ctx := context.Background()
	client := sampleClient()
	sample := &Sample{
		Name:  "sample",
		Value: 1,
	}
	assert.Empty(sample.ID())

	// return nil if object not exists
	o, err := client.Get(ctx, &Sample{}, "no id")
	assert.Nil(err)
	assert.Nil(o)

	// not found
	exist, err := client.Exists(ctx, &Sample{}, "no id")
	assert.Nil(err)
	assert.False(exist)

	// set object with auto id
	err = client.Set(ctx, sample)
	assert.Nil(err)
	assert.NotEmpty(sample.ID())

	// found
	exist, err = client.Exists(ctx, &Sample{}, sample.ID())
	assert.Nil(err)
	assert.True(exist)

	// get saved object
	sample2, err := client.Get(ctx, &Sample{}, sample.ID())
	assert.Nil(err)
	assert.NotNil(sample2)
	assert.Equal(sample2.(*Sample).Name, sample.Name)
	sampleCreateTime := sample2.CreateTime()
	assert.False(sampleCreateTime.IsZero())
	assert.False(sample2.UpdateTime().IsZero())

	// set sample again
	sample.Name = "modified"
	err = client.Set(ctx, sample)
	assert.Nil(err)

	m, err := client.Get(ctx, &Sample{}, sample.ID())
	sampleM := m.(*Sample)
	assert.Nil(err)
	assert.NotNil(sampleM)
	assert.Equal("modified", sampleM.Name)

	// set nil object
	err = client.Set(ctx, nil)
	assert.NotNil(err)

	// delete object
	err = client.Delete(ctx, sample2)
	assert.Nil(err)

	// manual id
	sample = &Sample{
		Name:  "sample",
		Value: 1,
	}
	sample.SetID("gdb-client-test")
	err = client.Set(ctx, sample)
	assert.Nil(err)
	assert.Equal("gdb-client-test", sample.ID())

	sample3, err := client.Get(ctx, &Sample{}, "gdb-client-test")
	defer client.Delete(ctx, sample3)
	assert.Nil(err)
	assert.NotNil(sample3)
	assert.Equal(sample3.(*Sample).Name, sample.Name)

	// delete not exists object
	err = client.Delete(ctx, &Sample{})
	assert.NotNil(err)
}

func TestGdbSelectUpdateIncrementDelete(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	ctx := context.Background()
	client := sampleClient()

	sample := &Sample{
		Name:  "sample",
		Value: 6,
	}

	err := client.Set(ctx, sample)
	assert.Nil(err)

	// not exists
	value, err := client.Select(ctx, &Sample{}, "no id", "Value")
	assert.Nil(err)
	assert.Nil(value)

	// found
	value, err = client.Select(ctx, &Sample{}, sample.ID(), "Value")
	assert.Nil(err)
	assert.Equal(int64(6), value)

	// update
	err = client.Update(ctx, sample, map[string]interface{}{
		"Name":  "sample2",
		"Value": 2,
	})
	assert.Nil(err)

	name, err := client.Select(ctx, &Sample{}, sample.ID(), "Name")
	assert.Nil(err)
	assert.Equal("sample2", name)

	value, err = client.Select(ctx, &Sample{}, sample.ID(), "Value")
	assert.Nil(err)
	assert.Equal(int64(2), value)

	err = client.Increment(ctx, &Sample{}, "Value", 3)
	assert.NotNil(err)

	err = client.Increment(ctx, sample, "Value", 3)
	assert.Nil(err)

	value, err = client.Select(ctx, &Sample{}, sample.ID(), "Value")
	assert.Nil(err)
	assert.Equal(int64(5), value)

	err = client.Delete(ctx, sample)
	assert.Nil(err)
}

func TestGdbListQueryFindCount(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	ctx := context.Background()
	client := sampleClient()

	name1 := "testGdb-" + identifier.RandomString(6)
	name2 := "testGdb-" + identifier.RandomString(6)
	sample1 := &Sample{
		Name:  name1,
		Value: 1001,
	}
	sample2 := &Sample{
		Name:  name2,
		Value: 1002,
	}
	err := client.Set(ctx, sample1)
	assert.Nil(err)
	err = client.Set(ctx, sample2)
	assert.Nil(err)
	defer client.Delete(ctx, sample1)
	defer client.Delete(ctx, sample2)

	// not found
	obj, err := client.Query(&Sample{}).Where("Value", "==", 1002).ReturnFirst(ctx)
	assert.Nil(err)
	assert.NotNil(obj)

	// found
	list, err := client.List(ctx, &Sample{}, 2)
	assert.Nil(err)
	assert.True(len(list) >= 2)

	list, err = client.Query(&Sample{}).Return(ctx)
	assert.Nil(err)
	assert.True(len(list) >= 2)

	obj, err = client.Query(&Sample{}).Where("Value", "==", 1002).ReturnFirst(ctx)
	assert.Nil(err)
	assert.Equal(name2, (obj.(*Sample)).Name)
}

func TestConnectionContextCanceled(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	client := sampleClient()

	ctx := util.CanceledCtx()
	sample := &Sample{}
	sample.SetID("id")

	err := client.Set(ctx, sample)
	assert.NotNil(err)
	_, err = client.Get(ctx, &Sample{}, "no id")
	assert.NotNil(err)
	err = client.Delete(ctx, sample)
	assert.NotNil(err)
	_, err = client.List(ctx, &Sample{}, 1)
	assert.NotNil(err)
	_, err = client.Exists(ctx, &Sample{}, "no id")
	assert.NotNil(err)
	_, err = client.Select(ctx, &Sample{}, "not id", "Value")
	assert.NotNil(err)
	err = client.Update(ctx, sample, map[string]interface{}{
		"Name":  "Sample2",
		"Value": "2",
	})
	assert.NotNil(err)
	_, err = client.Clear(ctx, &Sample{}, 10)
	assert.NotNil(err)
	_, err = client.Query(&Sample{}).Return(ctx)
	assert.NotNil(err)
	err = client.Increment(ctx, sample, "Value", 2)
	assert.NotNil(err)
}

func TestGdbClear(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	ctx := context.Background()
	client := sampleClient()

	sample := &SampleClear{
		Name: "sampleClear",
	}
	err := client.Set(ctx, sample)
	assert.Nil(err)

	cleared, err := client.Clear(ctx, sample, 100)
	assert.Nil(err)
	assert.True(cleared)
}

func BenchmarkSetSpeed(b *testing.B) {
	ctx := context.Background()
	client := sampleClient()
	sample := &Sample{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sample.Name = "gdb-benchmark"
		err := client.Set(ctx, sample)
		if err != nil {
			return
		}
	}
	client.Delete(ctx, sample)
}

func BenchmarkUpdateSpeed(b *testing.B) {
	ctx := context.Background()
	client := sampleClient()

	sample := &Sample{}
	err := client.Set(ctx, sample)
	if err != nil {
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client.Update(ctx, sample, map[string]interface{}{
			"Name": "hello" + strconv.Itoa(i),
		})
	}
	client.Delete(ctx, sample)
}
