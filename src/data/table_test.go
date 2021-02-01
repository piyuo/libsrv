package data

import (
	"context"
	"testing"

	"github.com/piyuo/libsrv/src/env"
	"github.com/stretchr/testify/assert"
)

func TestNoErrorTest(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	g, err := NewSampleGlobalDB(ctx)
	assert.Nil(err)
	defer g.Close()
	table := g.SampleTable()

	assert.NotNil(table.Factory)
	assert.NotEmpty(table.UUID())

	obj := table.NewObject()
	assert.Equal("Sample", table.TableName)
	assert.NotNil(obj)
	assert.Empty((obj.(*Sample)).Name)

	obj2 := table.Factory
	assert.NotNil(obj2)
}

func TestFirstObjectTest(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	g, err := NewSampleGlobalDB(ctx)
	assert.Nil(err)
	defer g.Close()
	table := g.SampleTable()

	sample1 := &Sample{
		Name:  "a",
		Value: 1,
	}
	table.Set(ctx, sample1)

	obj, err := table.GetFirstObject(ctx)
	assert.Nil(err)
	assert.NotNil(obj)

	id, err := table.GetFirstID(ctx)
	assert.Nil(err)
	assert.NotEmpty(id)

	err = table.Delete(ctx, id)
	assert.Nil(err)
}

func TestSearch(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	g, err := NewSampleGlobalDB(ctx)
	assert.Nil(err)
	defer g.Close()
	table := g.SampleTable()

	sample1 := &Sample{
		Name:  "a",
		Value: 1,
	}
	sample2 := &Sample{
		Name:  "a",
		Value: 2,
	}
	table.Set(ctx, sample1)
	table.Set(ctx, sample2)

	list, err := table.SortList(ctx, "Name", "==", "a", "Value", DESC)
	assert.Nil(err)
	assert.Equal(2, len(list))
	obj1 := list[0].(*Sample)
	obj2 := list[1].(*Sample)
	assert.Equal(2, obj1.Value)
	assert.Equal(1, obj2.Value)

	list, err = table.SortList(ctx, "Name", "==", "a", "Value", ASC)
	assert.Nil(err)
	assert.Equal(2, len(list))
	obj1 = list[0].(*Sample)
	obj2 = list[1].(*Sample)
	assert.Equal(1, obj1.Value)
	assert.Equal(2, obj2.Value)
	table.Delete(ctx, obj1.ID)
	table.Delete(ctx, obj2.ID)
}

func TestObjectUserID(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	g, err := NewSampleGlobalDB(ctx)
	assert.Nil(err)
	defer g.Close()
	table := g.SampleTable()

	sample := &Sample{
		Name:  "a",
		Value: 1,
	}
	table.Set(ctx, sample)
	defer table.DeleteObject(ctx, sample)
	assert.Empty(sample.GetUserID())
	assert.Empty(sample.GetAccountID())

	ctx = env.SetUserID(ctx, "user1")
	ctx = env.SetAccountID(ctx, "account1")
	table.Set(ctx, sample)
	assert.Equal("user1", sample.GetUserID())
	assert.Equal("account1", sample.GetAccountID())
}

func TestEmptyEnvAccountIDUserID(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	g, err := NewSampleGlobalDB(ctx)
	assert.Nil(err)
	defer g.Close()
	table := g.SampleTable()

	sample := &Sample{
		DomainObject: DomainObject{
			UserID:    "myUserID",
			AccountID: "myAccountID",
		},
		Name:  "a",
		Value: 1,
	}
	table.Set(ctx, sample)
	defer table.DeleteObject(ctx, sample)

	assert.Equal("myUserID", sample.GetUserID())
	assert.Equal("myAccountID", sample.GetAccountID())

	ctx = env.SetUserID(ctx, "")
	ctx = env.SetAccountID(ctx, "")
	table.Set(ctx, sample)
	assert.Equal("myUserID", sample.GetUserID())
	assert.Equal("myAccountID", sample.GetAccountID())
}
