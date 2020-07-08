package data

import (
	"context"

	"github.com/pkg/errors"
)

// Counters store all counter
//
type Counters struct {
	conn      Connection
	tablename string
}

// SetConnection set connection for table
//
//	table.SetConnection(conn)
//
func (cs *Counters) SetConnection(conn Connection) {
	cs.conn = conn
}

// SetTableName set table name
//
//	table.SetTableName("sample")
//
func (cs *Counters) SetTableName(tablename string) {
	cs.tablename = tablename
}

// TableName return table name
//
//	table.TableName()
//
func (cs *Counters) TableName() string {
	return cs.tablename
}

// Counter return counter from data store, create one if not exist
//
//	counter,err = db.Counter(ctx,"", "myCounter",10)
//
func (cs *Counters) Counter(ctx context.Context, countername string, numshards int) (Counter, error) {
	if numshards <= 0 {
		numshards = 10
	}
	if numshards >= 100 {
		numshards = 100
	}

	if cs.tablename == "" {
		return nil, errors.New("table name can not be empty")
	}
	if countername == "" {
		return nil, errors.New("counter name can not be empty")
	}

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	return cs.conn.Counter(ctx, cs.tablename, countername, numshards)
}
