package data

import (
	"context"
)

const limitQueryDefault = 10
const limitTransactionClear = 10
const limitClear = 500

// DB represent DB public method
//
type DB interface {

	// Close connection
	//
	//	db.Close()
	//
	Close()

	// CreateNamespace create namespace, create new one if not exist
	//
	//	db, err := db.CreateNamespace(ctx)
	//
	CreateNamespace(ctx context.Context) error

	// DeleteNamespace delete namespace
	//
	//	err := db.DeleteNamespace(ctx)
	//
	DeleteNamespace(ctx context.Context) error

	// Transaction start a transaction
	//
	//	err := db.Transaction(ctx, func(ctx context.Context, tx Transaction) error {
	//		tx.Put(ctx, &greet1)
	//		return nil
	//	})
	//
	Transaction(ctx context.Context, callback func(ctx context.Context) error) error

	// IsInTransaction return true if connection is in transaction
	//
	//	inTx := db.IsInTransaction()
	//
	IsInTransaction() bool

	// Connection return current connection
	//
	//	conn := db.Connection()
	//
	Connection() Connection
}

// BaseDB represent document database
//
type BaseDB struct {
	DB

	// conn is current database connection
	//
	conn Connection
}

// Connection return current connection
//
//	conn := db.Connection()
//
func (db *BaseDB) Connection() Connection {
	return db.conn
}

// Close connection
//
//	db.Close()
//
func (db *BaseDB) Close() {
	if db.conn != nil {
		db.conn.Close()
		db.conn = nil
	}
}

// Transaction start a transaction
//
//	err := conn.Transaction(ctx, func(ctx context.Context) error {
//		accounts.Set(ctx, &greet1)
//		return nil
//	})
//
func (db *BaseDB) Transaction(ctx context.Context, callback func(ctx context.Context) error) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	return db.conn.Transaction(ctx, callback)
}

// IsInTransaction return true if connection is in transaction
//
//	inTx := conn.IsInTransaction()
//
func (db *BaseDB) IsInTransaction() bool {
	return db.conn.IsInTransaction()
}

// CreateNamespace create namespace, create new one if not exist
//
//	db, err := conn.CreateNamespace(ctx)
//
func (db *BaseDB) CreateNamespace(ctx context.Context) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	return db.conn.CreateNamespace(ctx)
}

// DeleteNamespace delete namespace
//
//	err := db.DeleteNamespace(ctx)
//
func (db *BaseDB) DeleteNamespace(ctx context.Context) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	return db.conn.DeleteNamespace(ctx)
}
