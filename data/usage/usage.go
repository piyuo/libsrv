package usage

import (
	"context"
	"time"

	data "github.com/piyuo/libsrv/data"
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

// usage implement Usage
//
type usage struct {
	Usage
	table *data.Table
}

type record struct {
	data.BaseObject `firestore:"-"`
	Key             string
	Time            time.Time
}

// NewUsage return Usage
//
func NewUsage(tablename string, conn data.Connection) Usage {
	table := &data.Table{
		Connection: conn,
		TableName:  tablename,
		Factory: func() data.Object {
			return &record{}
		},
	}
	return &usage{
		table: table,
	}
}

// Count return usage of duration
//
//	err = usage.Get(ctx, "aaa@mail.com", 1 day)
//
func (c *usage) Count(ctx context.Context, key string, duration time.Duration) (int, error) {
	list, err := c.table.Search(ctx, "Key", "==", key)
	if err != nil {
		return 0, errors.Wrap(err, "failed to search usage where key="+key)
	}
	So((obj.(*Sample)).Name, ShouldEqual, "sample")
	return 0
}

// Add usage
//
//	err = usage.Add(ctx, "aaa@mail.com", 10,)
//
func Add(ctx context.Context, key string) error {
	return nil
}

// Remove usage
//
//	err = usage.Add(ctx, "aaa@mail.com")
//
func Remove(ctx context.Context, key string) error {

}

// Maintenance remove usage that is over 1 month, return true if no more usage record need to be delete
//
//	err = usage.Add(ctx, "aaa@mail.com")
//
func Maintenance(ctx context.Context) bool {
	return false
}
