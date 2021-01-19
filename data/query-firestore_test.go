package data

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestQueryTest(t *testing.T) {
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
	defer table.DeleteObject(ctx, sample1)

	err = table.Set(ctx, sample2)
	assert.Nil(err)
	defer table.DeleteObject(ctx, sample2)

	// get full object
	list, err := table.Query().Where("Name", "==", "sample1").Execute(ctx)
	assert.Nil(err)
	assert.Equal(1, len(list))
	assert.Equal("sample1", (list[0].(*Sample)).Name)

	// factory has no object return must error
	bakFactory := table.Factory
	table.Factory = func() Object {
		return nil
	}
	listX, err := table.Query().Where("Name", "==", "sample1").Execute(ctx)
	assert.NotNil(err)
	assert.Nil(listX)
	table.Factory = bakFactory

	list, err = table.Query().Where("Name", "==", "sample2").Execute(ctx)
	assert.Nil(err)
	assert.Equal(1, len(list))
	assert.Equal("sample2", (list[0].(*Sample)).Name)

	list, err = table.Query().Where("Value", "==", 1).Execute(ctx)
	assert.Nil(err)
	assert.Equal(1, len(list))
	assert.Equal("sample1", (list[0].(*Sample)).Name)

	list, err = table.Query().Where("Value", "==", 2).Execute(ctx)
	assert.Nil(err)
	assert.Equal(1, len(list))
	assert.Equal("sample2", (list[0].(*Sample)).Name)

	//OrderBy,OrderByDesc
	list, err = table.Query().OrderBy("Name").Execute(ctx)
	assert.Nil(err)
	assert.Equal(2, len(list))
	assert.Equal("sample1", (list[0].(*Sample)).Name)

	list, err = table.Query().OrderByDesc("Name").Execute(ctx)
	assert.Nil(err)
	assert.Equal(2, len(list))
	assert.Equal("sample2", (list[0].(*Sample)).Name)

	//limit
	list, err = table.Query().OrderBy("Name").Limit(1).Execute(ctx)
	assert.Nil(err)
	assert.Equal(1, len(list))
	assert.Equal("sample1", (list[0].(*Sample)).Name)

	list, err = table.Query().OrderByDesc("Name").Limit(1).Execute(ctx)
	assert.Nil(err)
	assert.Equal(1, len(list))
	assert.Equal("sample2", (list[0].(*Sample)).Name)

	//startAt,startAfter,endAt,endBefore
	list, err = table.Query().OrderBy("Name").StartAt("sample2").Execute(ctx)
	assert.Nil(err)
	assert.Equal(1, len(list))
	assert.Equal("sample2", (list[0].(*Sample)).Name)

	list, err = table.Query().OrderBy("Name").StartAfter("sample1").Execute(ctx)
	assert.Nil(err)
	assert.Equal(1, len(list))
	assert.Equal("sample2", (list[0].(*Sample)).Name)

	list, err = table.Query().OrderBy("Name").EndAt("sample2").Execute(ctx)
	assert.Nil(err)
	assert.Equal(2, len(list))
	assert.Equal("sample1", (list[0].(*Sample)).Name)

	list, err = table.Query().OrderBy("Name").EndBefore("sample2").Execute(ctx)
	assert.Nil(err)
	assert.Equal(1, len(list))
	assert.Equal("sample1", (list[0].(*Sample)).Name)

	count, err := table.Query().Where("Name", "==", "sample1").Count(ctx)
	assert.Nil(err)
	assert.Equal(1, count)

	isEmpty, err := table.Query().Where("Name", "==", "sample1").IsEmpty(ctx)
	assert.Nil(err)
	assert.False(isEmpty)

	isExist, err := table.Query().Where("Name", "==", "sample1").IsExist(ctx)
	assert.Nil(err)
	assert.True(isExist)

}

func TestList(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	g, err := NewSampleGlobalDB(ctx)
	assert.Nil(err)
	defer g.Close()
	table := g.SampleTable()

	sample1 := &Sample{
		Name:  "sample",
		Value: 1,
	}
	sample2 := &Sample{
		Name:  "sample",
		Value: 2,
	}
	err = table.Set(ctx, sample1)
	assert.Nil(err)
	defer table.DeleteObject(ctx, sample1)
	err = table.Set(ctx, sample2)
	assert.Nil(err)
	defer table.DeleteObject(ctx, sample2)

	// get id only
	list, err := table.Query().Where("Name", "==", "sample").GetIDs(ctx)
	assert.Nil(err)
	assert.Equal(2, len(list))
	assert.NotEmpty(list[0])
	assert.NotEmpty(list[1])
	assert.NotEqual(list[1], list[0])

}

func TestGetFirstObjectTest(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	g, err := NewSampleGlobalDB(ctx)
	assert.Nil(err)
	defer g.Close()
	table := g.SampleTable()

	obj, err := table.Query().Where("Name", "==", "sample").GetFirstObject(ctx)
	assert.Nil(err)
	assert.Nil(obj)

	id, err := table.Query().Where("Name", "==", "sample").GetFirstID(ctx)
	assert.Nil(err)
	assert.Empty(id)

	sample1 := &Sample{
		Name:  "sample",
		Value: 1,
	}
	sample2 := &Sample{
		Name:  "sample",
		Value: 2,
	}
	err = table.Set(ctx, sample1)
	assert.Nil(err)
	defer table.DeleteObject(ctx, sample1)
	err = table.Set(ctx, sample2)
	assert.Nil(err)
	defer table.DeleteObject(ctx, sample2)

	// get top one object only
	obj, err = table.Query().Where("Name", "==", "sample").GetFirstObject(ctx)
	assert.Nil(err)
	assert.NotNil(obj)

	id, err = table.Query().Where("Name", "==", "sample").GetFirstID(ctx)
	assert.Nil(err)
	assert.NotEmpty(id)

	// set limit 2 still get 1 object
	obj, err = table.Query().Where("Name", "==", "sample").Limit(2).GetFirstObject(ctx)
	assert.Nil(err)
	assert.NotNil(obj)

}

func TestExecuteQueryID(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	g, err := NewSampleGlobalDB(ctx)
	assert.Nil(err)
	defer g.Close()
	table := g.SampleTable()

	obj, err := table.Query().Where("Name", "==", "sample").GetFirstObject(ctx)
	assert.Nil(err)
	assert.Nil(obj)

	id, err := table.Query().Where("Name", "==", "sample").GetFirstID(ctx)
	assert.Nil(err)
	assert.Empty(id)

	sample1 := &Sample{
		BaseObject: BaseObject{
			ID: "s1",
		},
		Name:  "sample",
		Value: 1,
	}
	sample2 := &Sample{
		BaseObject: BaseObject{
			ID: "s2",
		},
		Name:  "sample",
		Value: 2,
	}
	err = table.Set(ctx, sample1)
	assert.Nil(err)
	defer table.DeleteObject(ctx, sample1)

	err = table.Set(ctx, sample2)
	assert.Nil(err)
	defer table.DeleteObject(ctx, sample2)

	// get top one object only
	obj, err = table.Query().Where("ID", "==", "s1").GetFirstObject(ctx)
	assert.Nil(err)
	assert.NotNil(obj)

}

func TestQueryNotExistFieldWillNotCauseError(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	g, err := NewSampleGlobalDB(ctx)
	assert.Nil(err)
	defer g.Close()
	table := g.SampleTable()

	sample1 := &Sample{
		BaseObject: BaseObject{
			ID: "s1",
		},
		Name:  "sample",
		Value: 1,
	}

	err = table.Set(ctx, sample1)
	assert.Nil(err)
	defer table.DeleteObject(ctx, sample1)

	// get top one object only
	obj, err := table.Query().Where("notExist", "<", time.Now().UTC()).GetFirstObject(ctx)
	assert.Nil(err)
	assert.Nil(obj)
}

func TestCreateUpdateTimeTest(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	g, err := NewSampleGlobalDB(ctx)
	assert.Nil(err)
	defer g.Close()
	table := g.SampleTable()

	sample1 := &Sample{
		BaseObject: BaseObject{
			ID: "s1",
		},
		Name:  "sample",
		Value: 1,
	}

	err = table.Set(ctx, sample1)
	assert.Nil(err)
	defer table.DeleteObject(ctx, sample1)

	// get top one object only
	obj, err := table.Query().Where("CreateTime", "<=", time.Now().UTC()).GetFirstObject(ctx)
	assert.Nil(err)
	assert.NotNil(obj)
}
