package data

import (
	"context"
	"strconv"
	"testing"

	"github.com/piyuo/libsrv/src/gaccount"
	"github.com/piyuo/libsrv/src/util"
	"github.com/stretchr/testify/assert"
)

func TestFirestoreNewDB(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	cred, err := gaccount.GlobalCredential(ctx)
	assert.Nil(err)
	db, err := firestoreNewConnection(ctx, cred)
	defer db.Close()
	assert.Nil(err)
	assert.NotNil(db)
}

func TestFirestoreGlobalDB(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	conn, err := FirestoreGlobalConnection(ctx)
	defer conn.Close()
	assert.Nil(err)
	assert.NotNil(conn)

	firestoreConn := conn.(*ConnectionFirestore)
	id := firestoreConn.errorID("tablename", "")
	assert.Equal("tablename", id)
	id = firestoreConn.errorID("tablename", "id")
	assert.Equal("tablename-id", id)
}

func TestFirestoreID(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	g, err := NewSampleGlobalDB(ctx)
	assert.Nil(err)
	defer g.Close()
	table := g.SampleTable()

	sample := &Sample{
		Name:  "sample",
		Value: 1,
	}
	assert.Empty(sample.ID)

	o, err := table.Get(ctx, "")
	assert.Nil(err)
	assert.Nil(o)

	// auto id
	err = table.Set(ctx, sample)
	assert.Nil(err)
	assert.NotEmpty(sample.ID)

	sample2, err := table.Get(ctx, sample.ID)
	assert.Nil(err)
	assert.NotNil(sample2)
	assert.Equal(sample2.(*Sample).Name, sample.Name)
	sampleCreateTime := sample2.GetCreateTime()
	assert.False(sampleCreateTime.IsZero())
	assert.False(sample2.GetUpdateTime().IsZero())

	// factory has no object return must error
	bakFactory := table.Factory
	table.Factory = func() Object {
		return nil
	}
	sampleX, err := table.Get(ctx, sample.ID)
	assert.NotNil(err)
	assert.Nil(sampleX)
	table.Factory = bakFactory

	// set sample again
	sample.Name = "modified"
	err = table.Set(ctx, sample)
	assert.Nil(err)

	m, err := table.Get(ctx, sample.ID)
	sampleM := m.(*Sample)
	assert.Nil(err)
	assert.NotNil(sampleM)
	assert.Equal("modified", sampleM.Name)

	// set nil object
	err = table.Set(ctx, nil)
	assert.NotNil(err)

	err = table.DeleteObject(ctx, sample2)
	assert.Nil(err)

	// manual id
	sample = &Sample{
		Name:  "sample",
		Value: 1,
	}
	sample.ID = "sample-id"
	err = table.Set(ctx, sample)
	assert.Nil(err)
	assert.Equal("sample-id", sample.ID)

	sample3, err := table.Get(ctx, "sample-id")
	assert.Nil(err)
	assert.NotNil(sample3)
	assert.Equal(sample3.(*Sample).Name, sample.Name)

	err = table.DeleteObject(ctx, sample3)
	assert.Nil(err)

}

func TestSetGetExistDelete(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	g, err := NewSampleGlobalDB(ctx)
	assert.Nil(err)
	defer g.Close()
	table := g.SampleTable()

	sample := &Sample{
		Name:  "sample",
		Value: 1,
	}

	err = table.Set(ctx, sample)
	assert.Nil(err)
	sampleID := sample.ID
	sample2, err := table.Get(ctx, sampleID)
	assert.Nil(err)
	assert.NotNil(sample2)
	assert.Equal(sample2.(*Sample).Name, sample.Name)

	exist, err := table.IsExists(ctx, sampleID)
	assert.Nil(err)
	assert.True(exist)

	exist, err = table.IsExists(ctx, "")
	assert.Nil(err)
	assert.False(exist)

	err = table.Delete(ctx, sampleID)
	assert.Nil(err)

	exist, err = table.IsExists(ctx, sampleID)
	assert.Nil(err)
	assert.False(exist)

	sample3, err := table.Get(ctx, sampleID)
	assert.Nil(err)
	assert.Nil(sample3)

	err = table.Clear(ctx)
	assert.Nil(err)
}

func TestSelectUpdateIncrementDelete(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	g, err := NewSampleGlobalDB(ctx)
	assert.Nil(err)
	defer g.Close()
	table := g.SampleTable()

	sample := &Sample{
		Name:  "sample",
		Value: 6,
	}
	err = table.Set(ctx, sample)
	assert.Nil(err)

	value, err := table.Select(ctx, "NotExistID", "Value")
	assert.Nil(err)
	assert.Nil(value)
	value, err = table.Select(ctx, sample.ID, "Value")
	assert.Nil(err)
	assert.Equal(int64(6), value)

	err = table.Update(ctx, "NotExistID", map[string]interface{}{
		"Name":  "sample2",
		"Value": 2,
	})
	assert.Nil(err)

	err = table.Delete(ctx, "NotExistID")
	assert.Nil(err)

	err = table.Update(ctx, sample.ID, map[string]interface{}{
		"Name":  "sample2",
		"Value": 2,
	})
	assert.Nil(err)

	name, err := table.Select(ctx, sample.ID, "Name")
	assert.Nil(err)
	assert.Equal("sample2", name)

	value, err = table.Select(ctx, sample.ID, "Value")
	assert.Nil(err)
	assert.Equal(int64(2), value)

	err = table.Increment(ctx, "NotExistID", "Value", 3)
	assert.NotNil(err)

	err = table.Delete(ctx, "NotExistID")
	assert.Nil(err)

	err = table.Increment(ctx, sample.ID, "Value", 3)
	assert.Nil(err)

	value, err = table.Select(ctx, sample.ID, "Value")
	assert.Nil(err)
	assert.Equal(int64(5), value)

	err = table.DeleteObject(ctx, sample)
	assert.Nil(err)

}

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
	table.Factory = func() Object {
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
	_, err = table.SortList(ctx, "Name", "==", "1", "", ASC)
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
