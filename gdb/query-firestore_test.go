package gdb

import (
	"context"
	"testing"
	"time"

	"github.com/piyuo/libsrv/db"
	"github.com/piyuo/libsrv/identifier"
	"github.com/stretchr/testify/assert"
)

func TestQuery(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	ctx := context.Background()
	client := sampleClient()
	rand := identifier.RandomString(8)
	name1 := "test-query-1-" + rand
	name2 := "test-query-2-" + rand

	sample1 := &Sample{
		Name:  name1,
		Value: 1,
		Tag:   rand,
	}
	sample2 := &Sample{
		Name:  name2,
		Value: 2,
		Tag:   rand,
	}
	err := client.Set(ctx, sample1)
	assert.Nil(err)
	defer client.Delete(ctx, sample1)
	err = client.Set(ctx, sample2)
	assert.Nil(err)
	defer client.Delete(ctx, sample2)

	// no obj will result error
	_, err = client.Query(nil).Where("ID", "==", sample1.ID()).ReturnFirstID(ctx)
	assert.NotNil(err)

	// get full object from id
	firstID, err := client.Query(&Sample{}).Where("ID", "==", sample1.ID()).ReturnFirstID(ctx)
	assert.Nil(err)
	assert.Equal(sample1.ID(), firstID)

	// get full object from name
	list, err := client.Query(&Sample{}).Where("Name", "==", name1).Return(ctx)
	assert.Nil(err)
	assert.Equal(1, len(list))
	assert.Equal(name1, (list[0].(*Sample)).Name)

	list, err = client.Query(&Sample{}).Where("Value", "==", 2).Return(ctx)
	assert.Nil(err)
	assert.True(len(list) >= 1)

	//OrderBy,OrderByDesc
	count, err := client.Query(&Sample{}).OrderBy("Name").ReturnCount(ctx)
	assert.Nil(err)
	assert.True(count >= 2)

	count, err = client.Query(&Sample{}).OrderByDesc("Name").ReturnCount(ctx)
	assert.Nil(err)
	assert.True(count >= 2)

	//limit
	list, err = client.Query(&Sample{}).Where("Tag", "==", rand).Limit(1).Return(ctx)
	assert.Nil(err)
	assert.Equal(1, len(list))

	//startAt,startAfter,endAt,endBefore
	count, err = client.Query(&Sample{}).OrderBy("Name").StartAt(name2).ReturnCount(ctx)
	assert.Nil(err)
	assert.True(count >= 1)

	count, err = client.Query(&Sample{}).OrderBy("Name").StartAfter(name1).ReturnCount(ctx)
	assert.Nil(err)
	assert.True(count >= 1)

	count, err = client.Query(&Sample{}).OrderBy("Name").EndAt(name2).ReturnCount(ctx)
	assert.Nil(err)
	assert.True(count >= 1)

	count, err = client.Query(&Sample{}).OrderBy("Name").EndBefore(name2).ReturnCount(ctx)
	assert.Nil(err)
	assert.True(count >= 1)

	// return count/is empty/is exists
	count, err = client.Query(&Sample{}).Where("Name", "==", name1).ReturnCount(ctx)
	assert.Nil(err)
	assert.Equal(1, count)

	isEmpty, err := client.Query(&Sample{}).Where("Name", "==", name1).ReturnEmpty(ctx)
	assert.Nil(err)
	assert.False(isEmpty)

	isExist, err := client.Query(&Sample{}).Where("Name", "==", name1).ReturnExists(ctx)
	assert.Nil(err)
	assert.True(isExist)
}

func TestQueryList(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	ctx := context.Background()
	client := sampleClient()
	rand := identifier.RandomString(8)
	name1 := "test-query-list-1-" + rand
	name2 := "test-query-list-2-" + rand

	sample1 := &Sample{
		Name:  name1,
		Value: 1,
		Tag:   rand,
	}
	sample2 := &Sample{
		Name:  name2,
		Value: 2,
		Tag:   rand,
	}
	err := client.Set(ctx, sample1)
	assert.Nil(err)
	defer client.Delete(ctx, sample1)
	err = client.Set(ctx, sample2)
	assert.Nil(err)
	defer client.Delete(ctx, sample2)

	// no obj will result error
	_, err = client.Query(nil).Where("Tag", "==", rand).ReturnID(ctx)
	assert.NotNil(err)

	// get id only
	list, err := client.Query(&Sample{}).Where("Tag", "==", rand).ReturnID(ctx)
	assert.Nil(err)
	assert.Equal(2, len(list))
	assert.NotEmpty(list[0])
	assert.NotEmpty(list[1])
	assert.NotEqual(list[1], list[0])

	// no obj will result error
	obj, err := client.Query(nil).Where("Name", "==", name1).ReturnFirst(ctx)
	assert.NotNil(err)
	assert.Nil(obj)

	// first
	obj, err = client.Query(&Sample{}).Where("Name", "==", name1).ReturnFirst(ctx)
	assert.Nil(err)
	assert.Equal(sample1.ID(), obj.ID())

	id, err := client.Query(&Sample{}).Where("Name", "==", name1).ReturnFirstID(ctx)
	assert.Nil(err)
	assert.Equal(sample1.ID(), id)

	// object not found will return empty string
	id, err = client.Query(&Sample{}).Where("Name", "==", "not exist").ReturnFirstID(ctx)
	assert.Nil(err)
	assert.Empty(id)
}

func TestQueryNotExistFieldWillNotError(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	ctx := context.Background()
	client := sampleClient()
	// get top one object only
	obj, err := client.Query(&Sample{}).Where("notExist", "<", time.Now().UTC()).ReturnFirst(ctx)
	assert.Nil(err)
	assert.Nil(obj)
}

func TestQueryTime(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	ctx := context.Background()
	client := sampleClient()
	rand := identifier.RandomString(8)
	name := "test-query-time-" + rand

	sample := &Sample{
		Name:  name,
		Value: 1,
	}

	err := client.Set(ctx, sample)
	assert.Nil(err)
	defer client.Delete(ctx, sample)

	// get top one object only
	obj, err := client.Query(&Sample{}).Where("CreateTime", "<=", time.Now().Add(5*time.Second).UTC()).ReturnFirst(ctx)
	assert.Nil(err)
	assert.NotNil(obj)
}

func TestQueryDelete(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	ctx := context.Background()
	client := sampleClient()
	rand := identifier.RandomString(8)
	name := "test-query-delete-" + rand

	sample := &Sample{
		Name:  name,
		Value: 1,
	}

	err := client.Set(ctx, sample)
	assert.Nil(err)

	found, err := client.Query(&Sample{}).Where("Name", "==", name).ReturnExists(ctx)
	assert.Nil(err)
	assert.True(found)

	done, count, err := client.Query(&Sample{}).Where("Name", "==", name).Delete(ctx, 100)
	assert.Nil(err)
	assert.True(done)
	assert.True(count > 0)

	found, err = client.Query(&Sample{}).Where("Name", "==", name).ReturnExists(ctx)
	assert.Nil(err)
	assert.False(found)
}

func TestQueryDeleteInTransaction(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	ctx := context.Background()
	client := sampleClient()

	// delete query is not support in transaction
	client.Transaction(ctx, func(ctx context.Context, tx db.Transaction) error {
		done, numDeleted, err := tx.Query(&Sample{}).Where("Tag", "==", "1").Delete(ctx, 10)
		assert.NotNil(err)
		assert.False(done)
		assert.Equal(0, numDeleted)
		return nil
	})
}

func TestQueryCleanup(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	ctx := context.Background()
	client := sampleClient()
	rand := identifier.RandomString(8)
	name := "test-query-cleanup-" + rand

	for i := 0; i < 30; i++ {
		sample := &Sample{
			Name: name,
		}
		err := client.Set(ctx, sample)
		assert.Nil(err)
	}

	done, err := client.Query(&Sample{}).Where("Name", "==", name).Cleanup(ctx)
	assert.Nil(err)
	assert.True(done)

	found, err := client.Query(&Sample{}).Where("Name", "==", name).ReturnExists(ctx)
	assert.Nil(err)
	assert.False(found)
}

func BenchmarkCleanupSetup(b *testing.B) {
	ctx := context.Background()
	client := sampleClient()
	name := "test-query-benchmark1"
	rand := identifier.RandomString(8)
	for i := 0; i < 1000; i++ {
		sample := &Sample{
			Name: name,
			Tag:  rand,
		}
		client.Set(ctx, sample)
	}
}

func BenchmarkCleanup(b *testing.B) {
	ctx := context.Background()
	client := sampleClient()
	name := "test-query-benchmark"
	client.Query(&Sample{}).Where("Name", "==", name).Cleanup(ctx)
}
