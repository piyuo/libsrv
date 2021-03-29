package gtask

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	ctx := context.Background()

	taskID, err := New(ctx, "task", "http://it-is-not-exist.com", []byte{}, "task-create", 1800, 3)
	assert.Nil(err)
	assert.NotEmpty(taskID)

	client := sampleClient()
	taskObj, err := client.Get(ctx, &Task{}, taskID)
	assert.Nil(err)
	task := taskObj.(*Task)
	assert.Equal(3, task.MaxRetry)
	assert.True(task.LockTime.IsZero())
	assert.Equal(0, task.Retry)
	defer client.Delete(ctx, task)

	defer TestModeBackNormal()
	TestModeAlwaySuccess()
	_, err = New(ctx, "task", "http://notExist", []byte{}, "my-task", 1800, 3)
	assert.Nil(err)

	TestModeAlwayFail()
	_, err = New(ctx, "task", "http://notExist", []byte{}, "my-task", 1800, 3)
	assert.NotNil(err)

}

func TestLock(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	ctx := context.Background()

	taskID, err := New(ctx, "task", "http://it-is-not-exist.com", []byte{}, "task-lock", 1800, 3)
	assert.Nil(err)

	err = Lock(ctx, taskID)
	assert.Nil(err)

	client := sampleClient()
	taskObj, err := client.Get(ctx, &Task{}, taskID)
	assert.Nil(err)
	task := taskObj.(*Task)
	assert.Equal(1, task.Retry)
	assert.False(task.LockTime.IsZero())
	defer client.Delete(ctx, task)

	// lock again
	err = Lock(ctx, taskID)
	assert.NotNil(err)

	defer TestModeBackNormal()
	TestModeAlwaySuccess()
	err = Lock(ctx, taskID)
	assert.Nil(err)

	TestModeAlwayFail()
	err = Lock(ctx, taskID)
	assert.NotNil(err)
}

func TestDelete(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	ctx := context.Background()

	taskID, err := New(ctx, "task", "http://it-is-not-exist.com", []byte{}, "task-delete", 1800, 3)
	assert.Nil(err)

	client := sampleClient()
	found, err := client.Exists(ctx, &Task{}, taskID)
	assert.Nil(err)
	assert.True(found)

	err = Delete(ctx, taskID)
	found, err = client.Exists(ctx, &Task{}, taskID)
	assert.Nil(err)
	assert.False(found)

	defer TestModeBackNormal()
	TestModeAlwaySuccess()
	err = Delete(ctx, taskID)
	assert.Nil(err)

	TestModeAlwayFail()
	err = Delete(ctx, taskID)
	assert.NotNil(err)
}
