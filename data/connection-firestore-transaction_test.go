package data

import (
	"context"
	"testing"

	util "github.com/piyuo/libsrv/util"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TransactionTest(t *testing.T) {
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

	assert.False(g.InTransaction())
	//success transaction
	err = g.Transaction(ctx, func(ctx context.Context) error {
		assert.True(g.InTransaction())
		err := table.Set(ctx, sample1)
		assert.Nil(err)
		err = table.Set(ctx, sample2)
		assert.Nil(err)
		return nil
	})
	assert.Nil(err)

	list, err := table.Query().OrderBy("Name").Execute(ctx)
	assert.Nil(err)
	assert.Equal(2, len(list))
	assert.Equal("sample1", (list[0].(*Sample)).Name)
	assert.Equal("sample2", (list[1].(*Sample)).Name)
	isEmpty, err := table.IsEmpty(ctx)
	assert.False(isEmpty)
	err = table.Clear(ctx)
	assert.Nil(err)

	//fail transaction
	err = g.Transaction(ctx, func(ctx context.Context) error {
		err = table.Set(ctx, sample1)
		assert.Nil(err)
		return errors.New("something wrong")
	})
	assert.NotNil(err)

	isEmpty, err = table.IsEmpty(ctx)
	assert.True(isEmpty)

	// success delete
	err = g.Transaction(ctx, func(ctx context.Context) error {
		err = table.Set(ctx, sample1)
		assert.Nil(err)
		err = table.DeleteObject(ctx, sample1)
		assert.Nil(err)
		return nil
	})
	assert.Nil(err)

	isEmpty, err = table.IsEmpty(ctx)
	assert.True(isEmpty)
	err = table.Clear(ctx)
	assert.Nil(err)

	// failed delete
	err = g.Transaction(ctx, func(ctx context.Context) error {
		err = table.Set(ctx, sample1)
		assert.Nil(err)
		err = table.DeleteObject(ctx, sample1)
		assert.Nil(err)
		return errors.New("something wrong")
	})
	assert.NotNil(err)

	isEmpty, err = table.IsEmpty(ctx)
	assert.True(isEmpty)
	err = table.Clear(ctx)
	assert.Nil(err)
}

func TestMethodTest(t *testing.T) {
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

	// get & deleteObject
	err = table.Set(ctx, sample1)
	assert.Nil(err)
	err = g.Transaction(ctx, func(ctx context.Context) error {
		sample, err := table.Get(ctx, sample1.ID)
		assert.Nil(err)
		err = table.DeleteObject(ctx, sample)
		assert.Nil(err)
		return nil
	})
	assert.Nil(err)
	isEmpty, err := table.IsEmpty(ctx)
	assert.True(isEmpty)

	// exist & list & delete
	err = table.Set(ctx, sample1)
	assert.Nil(err)
	err = g.Transaction(ctx, func(ctx context.Context) error {
		exist, err := table.Exist(ctx, sample1.ID)
		assert.Nil(err)
		assert.True(exist)
		objects, err := table.All(ctx)
		assert.Nil(err)
		assert.Equal(1, len(objects))
		err = table.Delete(ctx, sample1.ID)
		assert.Nil(err)
		return nil
	})
	assert.Nil(err)
	isEmpty, err = table.IsEmpty(ctx)
	assert.True(isEmpty)

	// select & update & Increment
	err = table.Set(ctx, sample1)
	assert.Nil(err)
	err = g.Transaction(ctx, func(ctx context.Context) error {
		name, err := table.Select(ctx, sample1.ID, "Name")
		assert.Nil(err)
		assert.Equal("sample1", name.(string))
		err = table.Update(ctx, sample1.ID, map[string]interface{}{
			"Name": "sample",
		})
		assert.Nil(err)
		err = table.Increment(ctx, sample1.ID, "Value", 1)
		assert.Nil(err)
		return nil
	})
	assert.Nil(err)
	name, err := table.Select(ctx, sample1.ID, "Name")
	assert.Nil(err)
	assert.Equal("sample", name.(string))
	value, err := table.Select(ctx, sample1.ID, "Value")
	assert.Nil(err)
	intValue, err := util.ToInt(value)
	assert.Nil(err)
	assert.Equal(2, intValue)
	table.DeleteObject(ctx, sample1)

	// query & clear
	err = table.Set(ctx, sample1)
	assert.Nil(err)
	err = table.Set(ctx, sample2)
	assert.Nil(err)
	err = g.Transaction(ctx, func(ctx context.Context) error {
		obj, err := table.Find(ctx, "Name", "==", "sample1")
		assert.Nil(err)
		assert.Equal("sample1", (obj.(*Sample)).Name)

		idList, err := table.Query().OrderBy("Name").GetIDs(ctx)
		assert.Nil(err)
		assert.Equal(2, len(idList))

		list, err := table.Query().OrderBy("Name").Execute(ctx)
		assert.Nil(err)
		assert.Equal(2, len(list))
		assert.Equal(sample1.Name, list[0].(*Sample).Name)
		assert.Equal(sample2.Name, list[1].(*Sample).Name)

		err = table.Clear(ctx)
		assert.Nil(err)
		return nil
	})
	assert.Nil(err)
	obj, err := table.Find(ctx, "Value", "==", 2)
	assert.Nil(err)
	assert.Nil(obj)
	isEmpty, err = table.IsEmpty(ctx)
	assert.True(isEmpty)

	// search & count & is empty
	err = table.Set(ctx, sample1)
	assert.Nil(err)
	err = g.Transaction(ctx, func(ctx context.Context) error {

		objects, err := table.List(ctx, "Name", "==", "sample1")
		assert.Nil(err)
		assert.Equal(1, len(objects))

		count, err := table.Count(ctx)
		assert.Nil(err)
		assert.Equal(1, count)

		empty, err := table.IsEmpty(ctx)
		assert.Nil(err)
		assert.False(empty)

		err = table.DeleteObject(ctx, sample1)
		assert.Nil(err)
		return nil
	})
	assert.Nil(err)
	isEmpty, err = table.IsEmpty(ctx)
	assert.True(isEmpty)
}
