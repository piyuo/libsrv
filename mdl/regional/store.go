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
	return db.newTable("store", func() data.ObjectRef {
		return &Store{}
	})
}
