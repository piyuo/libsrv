package gmdl

import (
	"context"

	data "github.com/piyuo/libsrv/data"
)

// Account represent single account
//
type Account struct {
	data.Object `firestore:"-"`
}

// AccountTable return account table
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

// AccountTotal return total account count
//
//	id := d.AccountTotal(ctx)
//
func (c *Counters) AccountTotal(ctx context.Context) (data.CounterRef, error) {
	return c.Counter(ctx, "accountTotal", 4)
}

// AccountID generate new account serial id
//
//	id := d.TableName()
//
func (s *Serial) AccountID(ctx context.Context) (string, error) {
	return s.Code(ctx, "accountID")
}
