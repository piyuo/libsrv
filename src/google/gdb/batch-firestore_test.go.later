package gdb

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGstoreBatchDeleteObjectTest(t *testing.T) {
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

	count, err := table.Query().Count(ctx)
	assert.Equal(0, count)

	batch := g.Batch()
	batch.Set(ctx, sample1) //batch mode do not return error
	batch.Set(ctx, sample2)
	err = g.BatchCommit(ctx)
	assert.Nil(err)
	list, err := table.Query().Execute(ctx)
	assert.Nil(err)
	assert.Equal(2, len(list))
	assert.False(g.InBatch())

	s1 := list[0].(*Sample)
	s2 := list[1].(*Sample)
	g.BatchBegin()
	table.Update(ctx, s1.ID, map[string]interface{}{
		"Value": 9,
	})
	table.Update(ctx, s2.ID, map[string]interface{}{
		"Value": 9,
	})
	err = g.BatchCommit(ctx)
	assert.Nil(err)

	list, err = table.Query().Where("Value", "==", 9).Execute(ctx)
	assert.Nil(err)
	assert.Equal(2, len(list))
	g1 := list[0].(*Sample)
	g2 := list[1].(*Sample)
	assert.Equal(9, g1.Value)
	assert.Equal(9, g2.Value)

	g.BatchBegin()
	table.Increment(ctx, s1.ID, "Value", 1)
	table.Increment(ctx, s2.ID, "Value", 1)
	err = g.BatchCommit(ctx)
	assert.Nil(err)
	list, err = table.Query().Where("Value", "==", 10).Execute(ctx)
	assert.Nil(err)
	assert.Equal(2, len(list))
	i1 := list[0].(*Sample)
	i2 := list[1].(*Sample)
	assert.Equal(10, i1.Value)
	assert.Equal(10, i2.Value)

	g.BatchBegin()
	table.DeleteObject(ctx, sample1) //batch mode do not return error
	table.DeleteObject(ctx, sample2)
	err = g.BatchCommit(ctx)
	assert.Nil(err)
	count, err = table.Query().Count(ctx)
	assert.Nil(err)
	assert.Equal(0, count)
}

func TestGstoreBatchDeleteTest(t *testing.T) {
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

	g.BatchBegin()
	table.Set(ctx, sample1) //batch mode do not return error
	table.Set(ctx, sample2)
	err = g.BatchCommit(ctx)
	assert.Nil(err)

	idList, err := table.Query().GetIDs(ctx)
	assert.Nil(err)
	assert.Equal(2, len(idList))

	g.BatchBegin()
	table.Delete(ctx, idList[0]) //batch mode do not return error
	table.Delete(ctx, idList[1])
	err = g.BatchCommit(ctx)
	assert.Nil(err)
	count, err := table.Query().Count(ctx)
	assert.Nil(err)
	assert.Equal(0, count)
}
