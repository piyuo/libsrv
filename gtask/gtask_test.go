package gtask

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMock(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	ctx := context.WithValue(context.Background(), MockNoError, "")
	_, err := New(ctx, "task", "http://not-exists", []byte{}, "my-task", 1800, 3)
	assert.Nil(err)
	err = Lock(ctx, "")
	assert.Nil(err)
	err = Delete(ctx, "")
	assert.Nil(err)

	ctx = context.WithValue(context.Background(), MockError, "")
	_, err = New(ctx, "task", "http://not-exists", []byte{}, "my-task", 1800, 3)
	assert.NotNil(err)
	err = Lock(ctx, "")
	assert.NotNil(err)
	err = Delete(ctx, "")
	assert.NotNil(err)
}

func TestNew(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	ctx := context.Background()

	taskID, err := New(ctx, "task", "http://not-exists", []byte{}, "task-create", 1800, 3)
	assert.Nil(err)
	assert.NotEmpty(taskID)

	client := sampleClient()
	taskObj, err := client.Get(ctx, &Task{}, taskID)
	defer client.Delete(ctx, taskObj)
	assert.Nil(err)
	task := taskObj.(*Task)
	assert.Equal(3, task.MaxRetry)
	assert.True(task.LockTime.IsZero())
	assert.Equal(0, task.Retry)
}

func TestLock(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	ctx := context.Background()

	taskID, err := New(ctx, "task", "http://not-exists", []byte{}, "task-lock", 1800, 3)
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
}

func TestDelete(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	ctx := context.Background()

	taskID, err := New(ctx, "task", "http://not-exists", []byte{}, "task-delete", 1800, 3)
	assert.Nil(err)

	client := sampleClient()
	found, err := client.Exists(ctx, &Task{}, taskID)
	assert.Nil(err)
	assert.True(found)

	err = Delete(ctx, taskID)
	assert.Nil(err)
	found, err = client.Exists(ctx, &Task{}, taskID)
	assert.Nil(err)
	assert.False(found)
}
