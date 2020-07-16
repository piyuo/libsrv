package gmdl

import (
	data "github.com/piyuo/libsrv/data"
)

// Account represent single account
//
type Account struct {
	data.Object `firestore:"-"`
}

// AccountTable return account table
//
//	counter := db.AccountTable()
//
func (db *DB) AccountTable() *data.Table {
	return &data.Table{
		Connection: db.Connection,
		TableName:  "account",
		Factory: func() data.ObjectRef {
			return &Account{}
		},
	}
}

// AccountTotal return account total counter
//
//	counter := d.AccountCounter()
//
func (c *Counters) AccountTotal() data.CounterRef {
	return c.Counter("AccountTotal", 100)
}

// AccountID return account id coder
//
//	coder := d.AccountID()
//
func (c *Coders) AccountID() data.CoderRef {
	return c.Coder("AccountID", 100)
}
