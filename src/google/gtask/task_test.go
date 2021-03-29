package gtask

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTask(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	task := newTask("TestTask", "TestTask", 1800, 3)
	assert.NotNil(task)
	assert.Equal(3, task.MaxRetry)
	assert.Equal("TestTask", task.ID())

	assert.Equal("TestTask", task.Collection())
	assert.NotNil(task.Factory())
}

func TestTaskLock(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	taskID := "TestTaskLock"
	task := newTask(taskID, "TestTask", 1800, 3)

	err := task.Lock()
	assert.Nil(err)
	assert.Equal(1, task.Retry)

	// max retry exceeded
	task = newTask(taskID, "TestTask", 1800, 3)
	task.MaxRetry = 3
	task.Retry = 3
	err = task.Lock()
	assert.NotNil(err)

	// task running
	task = newTask(taskID, "TestTask", 1800, 3)
	task.LockTime = time.Now().UTC()
	err = task.Lock()
	assert.NotNil(err)
}
