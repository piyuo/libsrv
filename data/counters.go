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

// Counter return counter from data store, create one if not exist
//
//	counter,err = counters.Counter(ctx,"", "myCounter",10)
//
func (cs *Counters) Counter(ctx context.Context, countername string, numshards int) (CounterRef, error) {
	if numshards <= 0 {
		numshards = 10
	}

	if cs.TableName == "" {
		return nil, errors.New("table name can not be empty")
	}
	if countername == "" {
		return nil, errors.New("counter name can not be empty")
	}

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	return cs.Connection.Counter(ctx, cs.TableName, countername, numshards)
}

// Delete counter
//
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
