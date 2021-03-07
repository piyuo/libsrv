package server

import (
	"context"
	"time"

	data "github.com/piyuo/libsrv/src/data"
	"github.com/pkg/errors"
)

// TaskLock keep task lock records
//
type TaskLock struct {
	data.BaseObject
}

// Database is firestore database
//
type Database struct {
	data.BaseDB
}

// TaskLockTable return TaskLock table
//
//	table := db.TaskLockTable()
//
func (c *Database) TaskLockTable() *data.Table {
	return &data.Table{
		Connection: c.Connection,
		TableName:  "TaskLock",
		Factory: func() data.Object {
			return &TaskLock{}
		},
	}
}

// New global db instance
//
//	db, err := Database.New(ctx)
//	if err != nil {
//		return err
//	}
//	defer db.Close()
//
func New(ctx context.Context) (*Database, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	conn, err := data.ConnectGlobalFirestore(ctx)
	if err != nil {
		return nil, err
	}

	c := &Database{
		BaseDB: data.BaseDB{Connection: conn},
	}
	return c, nil
}

// CreateTaskLock create lock on task
//
//	err = db.CreateTaskLock(ctx, "lock id")
//
func (c *Database) CreateTaskLock(ctx context.Context, lockID string) error {
	lock := &TaskLock{}
	lock.SetID(lockID)
	return c.TaskLockTable().Set(ctx, lock)
}

// DeleteTaskLock delete lock on task
//
//	err = db.DeleteTaskLock(ctx, "lock id")
//
func (c *Database) DeleteTaskLock(ctx context.Context, lockID string) error {

	return c.TaskLockTable().Delete(ctx, lockID)
}

// IsTaskLockExists check lock exists
//
//	found, createTime, err := db.IsTaskLockExists(ctx, "lock id")
//
func (c *Database) IsTaskLockExists(ctx context.Context, lockID string) (bool, time.Time, error) {
	obj, err := c.TaskLockTable().Get(ctx, lockID)
	if err != nil {
		return false, time.Time{}, errors.Wrap(err, "failed to get task lock:"+lockID)
	}
	if obj != nil {
		return true, obj.GetCreateTime(), nil
	}
	return false, time.Time{}, nil
}

// LockTask lock task for 15 mins
//
//	locked, err := db.LockTask(ctx, "lock id")
//
func (c *Database) LockTask(ctx context.Context, lockID string, duration time.Duration) (bool, error) {
	found, createTime, err := c.IsTaskLockExists(ctx, lockID)
	if err != nil {
		return false, errors.Wrap(err, "failed to check lock exist:"+lockID)
	}
	if !found {
		if err := c.CreateTaskLock(ctx, lockID); err != nil {
			return false, errors.Wrap(err, "failed to create task lock")
		}
		return true, nil
	}

	deadline := time.Now().UTC().Add(-duration)
	if createTime.Before(deadline) {
		// this target is too old, maybe something went wrong
		if err := c.TaskLockTable().Update(ctx, lockID, map[string]interface{}{"CreateTime": time.Now().UTC()}); err != nil {
			return false, errors.Wrap(err, "failed to update task lock")
		}
		return true, nil
	}
	return false, nil
}
