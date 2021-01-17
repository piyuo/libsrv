package data

import (
	"context"
	"testing"

	"github.com/piyuo/libsrv/session"
	"github.com/stretchr/testify/assert"
)

func TestTable(t *testing.T) {
	ctx := context.Background()
	dbG, dbR := createSampleDB()
	defer removeSampleDB(dbG, dbR)
	tableG, tableR := createSampleTable(dbG, dbR)
	defer removeSampleTable(tableG, tableR)

	noErrorTest(ctx, t, tableG)
	searchTest(ctx, t, tableG)
	firstObjectTest(ctx, t, tableG)
}

func noErrorTest(ctx context.Context, t *testing.T, table *Table) {
	assert := assert.New(t)
	assert.NotNil(table.Factory)
	assert.NotEmpty(table.UUID())

	obj := table.NewObject()
	assert.Equal("Sample", table.TableName)
	assert.NotNil(obj)
	assert.Empty((obj.(*Sample)).Name)

	obj2 := table.Factory
	assert.NotNil(obj2)
}

func firstObjectTest(ctx context.Context, t *testing.T, table *Table) {
	assert := assert.New(t)
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

func searchTest(ctx context.Context, t *testing.T, table *Table) {
	assert := assert.New(t)

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

func TestChangedBy(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	dbG, dbR := createSampleDB()
	defer removeSampleDB(dbG, dbR)
	tableG, tableR := createSampleTable(dbG, dbR)
	defer removeSampleTable(tableG, tableR)
	sample := &Sample{
		Name:  "a",
		Value: 1,
	}
	tableG.Set(ctx, sample)
	assert.Empty(sample.GetBy())

	ctx = session.SetUserID(ctx, "user1")
	tableG.Set(ctx, sample)
	assert.Equal("user1", sample.GetBy())
}
