package rmdl

import (
	data "github.com/piyuo/libsrv/data"
)

// Store represent single store
//
type Store struct {
	data.Object `firestore:"-"`
}

// StoreTable return store table
//
func (db *DB) StoreTable() *data.Table {
	return &data.Table{
		Connection: db.Connection,
		TableName:  "store",
		Factory: func() data.ObjectRef {
			return &Store{}
		},
	}
}
