package gmdl

import (
	"context"

	data "github.com/piyuo/libsrv/data"
)

// DB represent global database
//
type DB struct {
	data.DB
}

// NewDB create db instance
//
func NewDB(ctx context.Context) (*DB, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	conn, err := data.FirestoreGlobalConnection(ctx)
	if err != nil {
		return nil, err
	}

	db := &DB{
		DB: data.DB{Connection: conn},
	}
	return db, nil
}

// Counters return collection of counter
//
func (db *DB) Counters() *Counters {
	return &Counters{
		Counters: data.Counters{
			Connection: db.Connection,
			TableName:  "count",
		},
	}
}

// Serials return collection of serial
//
func (db *DB) Serials() *Serials {
	return &Serials{
		Serials: data.Serials{
			Connection: db.Connection,
			TableName:  "serial",
		}}
}

// Coders return collection of coder
//
func (db *DB) Coders() *Coders {
	return &Coders{
		Coders: data.Coders{
			Connection: db.Connection,
			TableName:  "serial",
		}}
}

// Counters is collection of global usage counters
//
type Counters struct {
	data.Counters `firestore:"-"`
}

// Serials keep all serial numbers
//
type Serials struct {
	data.Serials `firestore:"-"`
}

// Coders keep all coders
//
type Coders struct {
	data.Coders `firestore:"-"`
}
