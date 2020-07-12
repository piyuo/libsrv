package data

import (
	"context"

	"github.com/pkg/errors"
)

// Counters store all counter
//
type Counters struct {
	Connection ConnectionRef

	//TableName is counter table name
	//
	TableName string
}

// Counter return counter from database, create one if not exist
//
//	counters := db.Counters()
//	orderCountCounter,err = counters.Counter("order-count",10)
//
func (cs *Counters) Counter(countername string, numshards int) CounterRef {
	if numshards <= 0 {
		numshards = 10
	}
	return cs.Connection.Counter(cs.TableName, countername, numshards)
}

// Delete counter from database
//
//	counters := db.Counter()
//	err = counters.Delete(ctx, "myCounter")
//
func (cs *Counters) Delete(ctx context.Context, countername string) error {

	if cs.TableName == "" {
		return errors.New("table name can not be empty")
	}
	if countername == "" {
		return errors.New("counter name can not be empty")
	}

	if ctx.Err() != nil {
		return ctx.Err()
	}
	return cs.Connection.DeleteCounter(ctx, cs.TableName, countername)
}
