package gdb

import (
	"context"
	"testing"

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

/*
func TestListQueryFindCountClear(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	g, err := NewSampleGlobalDB(ctx)
	assert.Nil(err)
	defer g.Close()
	table := g.SampleTable()

	sample1 := &Sample{
		Name:  "sample1",
		Value: 1,
	}
	sample2 := &Sample{
		Name:  "sample2",
		Value: 2,
	}
	err = table.Set(ctx, sample1)
	assert.Nil(err)
	err = table.Set(ctx, sample2)
	assert.Nil(err)

	list, err := table.All(ctx)
	assert.Nil(err)
	assert.Equal(2, len(list))
	assert.Contains(list[0].(*Sample).Name, "sample")
	assert.Contains(list[1].(*Sample).Name, "sample")

	// factory has no object return must error
	bakFactory := table.Factory
	table.Factory = func() data.Object {
		return nil
	}
	listX, err := table.All(ctx)
	assert.NotNil(err)
	assert.Nil(listX)
	table.Factory = bakFactory

	obj, err := table.Find(ctx, "Name", "==", "sample1")
	assert.Nil(err)
	assert.Equal("sample1", (obj.(*Sample)).Name)

	list, err = table.Query().OrderBy("Name").Execute(ctx)
	assert.Nil(err)
	assert.Equal(2, len(list))
	assert.Equal(sample1.Name, list[0].(*Sample).Name)
	assert.Equal(sample2.Name, list[1].(*Sample).Name)

	obj, err = table.Find(ctx, "Value", "==", 2)
	assert.Nil(err)
	assert.Equal("sample2", (obj.(*Sample)).Name)

	err = table.Clear(ctx)
	assert.Nil(err)

	obj, err = table.Find(ctx, "Value", "==", 2)
	assert.Nil(err)
	assert.Nil(obj)
}

func TestSearchCountIsEmpty(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	g, err := NewSampleGlobalDB(ctx)
	assert.Nil(err)
	defer g.Close()
	table := g.SampleTable()

	sample := &Sample{
		Name:  "sample",
		Value: 0,
	}
	err = table.Set(ctx, sample)
	assert.Nil(err)

	objects, err := table.List(ctx, "Name", "==", "sample")
	assert.Nil(err)
	assert.Equal(1, len(objects))

	count, err := table.Count(ctx)
	assert.Nil(err)
	assert.Equal(1, count)

	empty, err := table.IsEmpty(ctx)
	assert.Nil(err)
	assert.False(empty)

	err = table.DeleteObject(ctx, sample)
	assert.Nil(err)
}

func TestDelete(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	g, err := NewSampleGlobalDB(ctx)
	assert.Nil(err)
	defer g.Close()
	table := g.SampleTable()

	sample := &Sample{
		Name:  "sample",
		Value: 0,
	}
	err = table.DeleteObject(ctx, sample)
	assert.Nil(err)
	err = table.Delete(ctx, "NotExistID")
	assert.Nil(err)
	err = table.Delete(ctx, "NotExistID")
	assert.Nil(err)

	err = table.Set(ctx, sample)
	assert.Nil(err)
	exist, err := table.IsExists(ctx, sample.ID)
	assert.Nil(err)
	assert.True(exist)

	sample2 := &Sample{}
	sample2.ID = sample.ID
	err = table.DeleteObject(ctx, sample2)
	assert.Nil(err)
	exist, err = table.IsExists(ctx, sample.ID)
	assert.Nil(err)
	assert.False(exist)

	//delete batch
	//delete empty batch
	ids := []string{}
	err = table.DeleteBatch(ctx, ids)
	assert.Nil(err)

	err = table.Set(ctx, sample)
	assert.Nil(err)
	exist, err = table.IsExists(ctx, sample.ID)
	assert.Nil(err)
	assert.True(exist)

	ids = []string{sample.ID}
	err = table.DeleteBatch(ctx, ids)
	assert.Nil(err)
	exist, err = table.IsExists(ctx, sample.ID)
	assert.Nil(err)
	assert.False(exist)

}

func TestConnectionContextCanceled(t *testing.T) {
	assert := assert.New(t)
	g, err := NewSampleGlobalDB(context.Background())
	assert.Nil(err)
	defer g.Close()
	table := g.SampleTable()

	ctx := util.CanceledCtx()
	sample := &Sample{}

	err = table.Set(ctx, sample)
	assert.NotNil(err)
	_, err = table.Get(ctx, "notexist")
	assert.NotNil(err)
	err = table.Delete(ctx, "notexist")
	assert.NotNil(err)
	err = table.DeleteObject(ctx, sample)
	assert.NotNil(err)
	err = table.DeleteBatch(ctx, []string{})
	assert.NotNil(err)
	_, err = table.All(ctx)
	assert.NotNil(err)
	_, err = table.IsExists(ctx, "notexist")
	assert.NotNil(err)
	_, err = table.Select(ctx, "notexist", "Value")
	assert.NotNil(err)
	err = table.Update(ctx, "notexist", map[string]interface{}{
		"Name":  "Sample2",
		"Value": "2",
	})
	assert.NotNil(err)
	err = table.Clear(ctx)
	assert.NotNil(err)
	_, err = table.Query().Execute(ctx)
	assert.NotNil(err)
	_, err = table.Find(ctx, "Value", "==", "2")
	assert.NotNil(err)
	_, err = table.Count(ctx)
	assert.NotNil(err)
	err = table.Increment(ctx, "notexist", "Value", 2)
	assert.NotNil(err)
	_, err = table.List(ctx, "Name", "==", "1")
	assert.NotNil(err)
	_, err = table.SortList(ctx, "Name", "==", "1", "", data.ASC)
	assert.NotNil(err)
	_, err = table.All(ctx)
	assert.NotNil(err)
	_, err = table.IsEmpty(ctx)
	assert.NotNil(err)
	err = table.Clear(ctx)
	assert.NotNil(err)
}

func BenchmarkPutSpeed(b *testing.B) {
	ctx := context.Background()
	dbG, err := NewSampleGlobalDB(ctx)
	defer dbG.Close()

	table := dbG.SampleTable()
	dbR, err := NewSampleRegionalDB(ctx)
	defer dbR.Close()

	sample := &Sample{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sample.Name = "hello" + strconv.Itoa(i)
		err = table.Set(ctx, sample)
		if err != nil {
			return
		}
	}
	table.DeleteObject(ctx, sample)
}

func BenchmarkUpdateSpeed(b *testing.B) {
	ctx := context.Background()
	dbG, err := NewSampleGlobalDB(ctx)
	defer dbG.Close()
	table := dbG.SampleTable()
	dbR, err := NewSampleRegionalDB(ctx)
	defer dbR.Close()

	sample := &Sample{}
	err = table.Set(ctx, sample)
	if err != nil {
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		table.Update(ctx, sample.ID, map[string]interface{}{
			"Name": "hello" + strconv.Itoa(i),
		})
	}
	table.DeleteObject(ctx, sample)
}
*/
