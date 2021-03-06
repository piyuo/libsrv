package server

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestServerTaskLock(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	db, err := New(ctx)
	assert.Nil(err)
	defer db.Close()
	lockID := "testLock"

	err = db.CreateTaskLock(ctx, lockID)
	assert.Nil(err)

	found, createTime, err := db.IsTaskLockExists(ctx, lockID)
	assert.Nil(err)
	assert.True(found)
	assert.True(createTime.Before(time.Now().UTC()))

	err = db.DeleteTaskLock(ctx, lockID)
	assert.Nil(err)

	found, createTime, err = db.IsTaskLockExists(ctx, lockID)
	assert.Nil(err)
	assert.False(found)
	assert.True(createTime.Before(time.Now().UTC()))
}

func TestServerLockTask(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	db, err := New(ctx)
	assert.Nil(err)
	defer db.Close()
	lockID := "testLock"

	// when task lock not exists
	ready, err := db.LockTask(ctx, lockID, 15*time.Minute)
	assert.Nil(err)
	assert.True(ready)
	err = db.DeleteTaskLock(ctx, lockID)
	assert.Nil(err)

	// when a expired task lock exists
	lock := &TaskLock{}
	lock.SetID(lockID)
	lock.SetCreateTime(time.Now().UTC().Add(-16 * time.Minute))
	err = db.TaskLockTable().Set(ctx, lock)

	ready, err = db.LockTask(ctx, lockID, 15*time.Minute)
	assert.Nil(err)
	assert.True(ready)
	err = db.DeleteTaskLock(ctx, lockID)
	assert.Nil(err)

	// when a not expired task lock exist
	lockInProgress := &TaskLock{}
	lockInProgress.SetID(lockID)
	lockInProgress.SetCreateTime(time.Now().UTC().Add(-10 * time.Minute))
	err = db.TaskLockTable().Set(ctx, lockInProgress)

	ready, err = db.LockTask(ctx, lockID, 15*time.Minute)
	assert.Nil(err)
	assert.False(ready)
	err = db.DeleteTaskLock(ctx, lockID)
	assert.Nil(err)
}
