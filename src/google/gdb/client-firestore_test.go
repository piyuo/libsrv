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

func TestClientClose(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	ctx := context.Background()
	cred, err := gaccount.GlobalCredential(ctx)
	assert.Nil(err)
	client, err := NewClient(ctx, cred)
	assert.Nil(err)
	client.Close()
}

func TestClientCRUD(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	ctx := context.Background()
	client := sampleClient()
	name := "test-client-CRUD-" + identifier.RandomString(8)
	sample := &Sample{
		Name:  name,
		Value: 1,
	}
	assert.Empty(sample.ID())

	// return error if no id
	_, err := client.Get(ctx, &Sample{}, "")
	assert.NotNil(err)

	// return nil if object not exists
	o, err := client.Get(ctx, &Sample{}, "no id")
	assert.Nil(err)
	assert.Nil(o)

	// return error if no id
	_, err = client.Exists(ctx, &Sample{}, "")
	assert.NotNil(err)

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
		Name:  "sample-manual-id",
		Value: 1,
	}
	sample.SetID("my-id")
	err = client.Set(ctx, sample)
	assert.Nil(err)
	defer client.Delete(ctx, sample)
	assert.Equal("my-id", sample.ID())

	sample3, err := client.Get(ctx, &Sample{}, "my-id")
	assert.Nil(err)
	assert.NotNil(sample3)
	assert.Equal(sample3.(*Sample).Name, sample.Name)

	// delete not exists object
	err = client.Delete(ctx, &Sample{})
	assert.NotNil(err)
}

func TestClientUpdate(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	ctx := context.Background()
	client := sampleClient()

	name1 := "test-client-update-" + identifier.RandomString(8)
	name2 := "test-client-update-2-" + identifier.RandomString(8)
	sample := &Sample{
		Name:  name1,
		Value: 6,
	}

	err := client.Set(ctx, sample)
	assert.Nil(err)
	defer client.Delete(ctx, sample)

	// return error if no id
	_, err = client.Select(ctx, &Sample{}, "", "Value")
	assert.NotNil(err)

	// not exists
	value, err := client.Select(ctx, &Sample{}, "no id", "Value")
	assert.Nil(err)
	assert.Nil(value)

	// found
	value, err = client.Select(ctx, &Sample{}, sample.ID(), "Value")
	assert.Nil(err)
	assert.Equal(int64(6), value)

	// update

	// nil fields will not error
	err = client.Update(ctx, sample, nil)
	assert.Nil(err)

	// field not exist will not error
	err = client.Update(ctx, sample, map[string]interface{}{
		"NotExist": name2,
	})
	assert.Nil(err)

	// nothing to update will not error
	err = client.Update(ctx, sample, map[string]interface{}{})
	assert.Nil(err)

	// success
	err = client.Update(ctx, sample, map[string]interface{}{
		"Name":  name2,
		"Value": 2,
	})
	assert.Nil(err)

	// select
	name, err := client.Select(ctx, &Sample{}, sample.ID(), "Name")
	assert.Nil(err)
	assert.Equal(name2, name)

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

}

func TestClientList(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	ctx := context.Background()
	client := sampleClient()
	name1 := "test-client-list-" + identifier.RandomString(8)
	name2 := "test-client-list-" + identifier.RandomString(8)
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

func TestClientContextCanceled(t *testing.T) {
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
		"Name":  "ctx cancel",
		"Value": "2",
	})
	assert.NotNil(err)
	_, err = client.Query(&Sample{}).Return(ctx)
	assert.NotNil(err)
	err = client.Increment(ctx, sample, "Value", 2)
	assert.NotNil(err)
}

func TestClientDeleteAll(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	ctx := context.Background()
	client := sampleClient()

	name := "test-client-delete-all-" + identifier.RandomString(8)
	sample := &SampleEmpty{
		Name: name,
	}
	err := client.Set(ctx, sample)
	assert.Nil(err)

	cleared, err := client.(*ClientFirestore).deleteAll(ctx, sample, 100)
	assert.Nil(err)
	assert.True(cleared)
}

func BenchmarkClientSet(b *testing.B) {
	ctx := context.Background()
	client := sampleClient()
	name := "benchmark-client-set-" + identifier.RandomString(8)
	sample := &Sample{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sample.Name = name
		err := client.Set(ctx, sample)
		if err != nil {
			return
		}
	}
	client.Delete(ctx, sample)
}

func BenchmarkClientUpdate(b *testing.B) {
	ctx := context.Background()
	client := sampleClient()
	name := "benchmark-client-update-" + identifier.RandomString(8)
	sample := &Sample{}
	err := client.Set(ctx, sample)
	if err != nil {
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client.Update(ctx, sample, map[string]interface{}{
			"Name": name + strconv.Itoa(i),
		})
	}
	client.Delete(ctx, sample)
}
