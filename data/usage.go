package data

/*
import (
	"context"
	"time"

	"github.com/pkg/errors"
)

// Usage can track usage in certain duration
//
type Usage interface {
	// Count return usage of duration
	//
	//	err = usage.Get(ctx, "aaa@mail.com", 1 day)
	//
	Count(ctx context.Context, key string, duration time.Duration) (int, error)

	// Add usage
	//
	//	err = usage.Add(ctx, "aaa@mail.com", 10,)
	//
	Add(ctx context.Context, key string) error

	// Remove usage
	//
	//	err = usage.Add(ctx, "aaa@mail.com")
	//
	Remove(ctx context.Context, key string) error
}

// baseUsage implement Usage
//
type baseUsage struct {
	Usage

	// table is usage table
	//
	table *Table
}

type usage struct {
	BaseObject `firestore:"-"`
	// Group is group name, key are separate by group
	//
	Group string

	// Key name
	Key string

	// Time is record created time
	Time time.Time
}

// NewUsage return Usage
//
func NewUsage(conn Connection) Usage {
	table := &Table{
		Connection: conn,
		TableName:  "usage",
		Factory: func() Object {
			return &usage{}
		},
	}
	return &baseUsage{
		table: table,
	}
}

// Count return usage of duration
//
//	count,lastDuration,err = usage.Count(ctx, "aaa@mail.com", time.Duration(24)*time.Hour)
//
func (c *baseUsage) Count(ctx context.Context, group, key string, duration time.Duration) (int, time.Duration, error) {
	q := c.table.Query().Where("Group", "==", group).Where("Key", "==", key).Limit(10).OrderByDesc("Time")
	list, err := q.Execute(ctx)
	if err != nil {
		return 0, 0, errors.Wrap(err, "failed to count usage group: "+group+",key: "+key)
	}
	count := len(list)
	if count > 0 {
		u := list[0].(*usage)
		diff := time.Now().UTC().Sub(u.Time)
		return count, diff, nil
	}
	return 0, 0, nil
}

// Add usage
//
//	err = usage.Add(ctx,"email", "aaa@mail.com")
//
func (c *baseUsage) Add(ctx context.Context, group, key string) error {
	u := &usage{
		Group: group,
		Key:   key,
		Time:  time.Now().UTC(),
	}
	err := c.table.Set(ctx, u)
	if err != nil {
		return errors.Wrap(err, "failed to add usage group: "+group+", key= "+key)
	}
	return nil
}

// Remove usage
//
//	err = usage.Add(ctx,"email", "aaa@mail.com")
//
func (c *baseUsage) Remove(ctx context.Context, group, key string) error {
	q := c.table.Query().Where("Group", "==", group).Where("Key", "==", key).Limit(10)
	list, err := q.Execute(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to list usage group: "+group+",key: "+key)
	}
	for _, u := range list {
		if err := c.table.DeleteObject(ctx, u); err != nil {
			return err
		}
	}
	return nil
}

// Maintenance remove usage that is over 1 month, return true if no more usage record need to be delete
//
//	err = usage.Maintenance(ctx, "aaa@mail.com",time.Duration(1) * time.Second)
//
func Maintenance(ctx context.Context) bool {
	checkPoint := time.Now().UTC().Add(time.Duration(1) * time.Second)

	q := c.table.Query().Where("Time", "<", checkPoint).Limit(1000)
	list, err := q.Execute(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to list usage group: "+group+",key: "+key)
	}
	for _, u := range list {
		if err := c.table.DeleteObject(ctx, u); err != nil {
			return err
		}
	}
	return nil
}
*/
