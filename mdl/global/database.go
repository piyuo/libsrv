package gmdl

import (
	"context"

	data "github.com/piyuo/libsrv/data"
)

// Database represent global database
//
type Database struct {
	data.DocDB
}

// NewDatabase create database
//
func NewDatabase(ctx context.Context) (*Database, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	conn, err := data.FirestoreGlobalConnection(ctx, "")
	if err != nil {
		return nil, err
	}

	db := &Database{}
	db.SetConnection(conn)
	return db, nil
}

// AccountTable return account table
//
func (db *Database) AccountTable() data.Table {
	table := &data.DocTable{}
	table.SetConnection(db.Connection())
	table.SetTableName("account")
	return table
}

// Counters return global counters
//
func (db *Database) Counters() *Counters {
	counters := &Counters{}
	counters.SetConnection(db.Connection())
	counters.SetTableName("counter")
	return counters
}

// Serial return serial
//
func (db *Database) Serial() *Serial {
	serial := &Serial{}
	serial.SetConnection(db.Connection())
	serial.SetTableName("serial")
	return serial
}

// Counters is collection of global usage counters
//
type Counters struct {
	data.Counters `firestore:"-"`
}

// AccountTotal return total account count
//
//	id := d.AccountTotal(ctx)
//
func (c *Counters) AccountTotal(ctx context.Context) (data.Counter, error) {
	return c.Counter(ctx, "accountTotal", 4)
}

// Serial keep serial numbers
//
type Serial struct {
	data.Serial `firestore:"-"`
}

// AccountID generate new account serial id
//
//	id := d.TableName()
//
func (s *Serial) AccountID(ctx context.Context) (string, error) {
	return s.Code(ctx, "accountID")
}
