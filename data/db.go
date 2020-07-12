package data

import (
	"context"
)

const limitQueryDefault = 10
const limitTransactionClear = 10
const limitClear = 500

// DBRef represent DB public method
//
type DBRef interface {

	// Close connection
	//
	//	conn.Close()
	//
	Close()

	// CreateNamespace create namespace, create new one if not exist
	//
	//	dbRef, err := conn.CreateNamespace(ctx)
	//
	CreateNamespace(ctx context.Context) error

	// DeleteNamespace delete namespace
	//
	//	err := db.DeleteNamespace(ctx)
	//
	DeleteNamespace(ctx context.Context) error

	// Transaction start a transaction
	//
	//	err := conn.Transaction(ctx, func(ctx context.Context, tx Transaction) error {
	//		tx.Put(ctx, &greet1)
	//		return nil
	//	})
	//
	Transaction(ctx context.Context, callback func(ctx context.Context) error) error

	// IsInTransaction return true if connection is in transaction
	//
	//	inTx := conn.IsInTransaction()
	//
	IsInTransaction() bool
}

// DB represent document database
//
type DB struct {
	DBRef

	// Connection is database connection
	//
	Connection ConnectionRef
}

// Close connection
//
//	conn.Close()
//
func (db *DB) Close() {
	if db.Connection != nil {
		db.Connection.Close()
		db.Connection = nil
	}
}

// Transaction start a transaction
//
//	err := conn.Transaction(ctx, func(ctx context.Context) error {
//		accounts.Set(ctx, &greet1)
//		return nil
//	})
//
func (db *DB) Transaction(ctx context.Context, callback func(ctx context.Context) error) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	return db.Connection.Transaction(ctx, callback)
}

// IsInTransaction return true if connection is in transaction
//
//	inTx := conn.IsInTransaction()
//
func (db *DB) IsInTransaction() bool {
	return db.Connection.IsInTransaction()
}

// CreateNamespace create namespace, create new one if not exist
//
//	dbRef, err := conn.CreateNamespace(ctx)
//
func (db *DB) CreateNamespace(ctx context.Context) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	return db.Connection.CreateNamespace(ctx)
}

// DeleteNamespace delete namespace
//
//	err := db.DeleteNamespace(ctx)
//
func (db *DB) DeleteNamespace(ctx context.Context) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	return db.Connection.DeleteNamespace(ctx)
}
