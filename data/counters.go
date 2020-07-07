package data

import (
	"context"

	"github.com/pkg/errors"
)

// Counters represent usage table
//
type Counters interface {
	SetConnection(conn Connection)
	SetTableName(tablename string)
	TableName()
	Counter(ctx context.Context, countername string, numshards int) (Counter, error)
}

// DocCounters store all counter
//
type DocCounters struct {
	Counters
	conn      Connection
	tablename string
}

// SetConnection set connection for table
//
//	table.SetConnection(conn)
//
func (dcs *DocCounters) SetConnection(conn Connection) {
	dcs.conn = conn
}

// SetTableName set table name
//
//	table.SetTableName("sample")
//
func (dcs *DocCounters) SetTableName(tablename string) {
	dcs.tablename = tablename
}

// TableName return table name
//
//	table.TableName()
//
func (dcs *DocCounters) TableName() string {
	return dcs.tablename
}

// Counter return counter from data store, create one if not exist
//
//	counter,err = db.Counter(ctx,"", "myCounter",10)
//
func (dcs *DocCounters) Counter(ctx context.Context, countername string, numshards int) (Counter, error) {
	if numshards <= 0 {
		numshards = 10
	}
	if numshards >= 100 {
		numshards = 100
	}

	if dcs.tablename == "" {
		return nil, errors.New("table name can not be empty")
	}
	if countername == "" {
		return nil, errors.New("counter name can not be empty")
	}

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	return dcs.conn.Counter(ctx, dcs.tablename, countername, numshards)
}
