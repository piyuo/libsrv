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

	// BatchBegin put connection into batch mode. Set/Update/Delete will hold operation until CommitBatch
	//
	//	err := conn.BatchBegin()
	//
	BatchBegin()

	// InBatch return true if connection is in batch mode
	//
	//	inBatch := conn.InBatch()
	//
	InBatch() bool

	// BatchCommit commit batch operation
	//
	//	err := conn.BatchCommit(ctx)
	//
	BatchCommit(ctx context.Context) error

	// Transaction start a transaction
	//
	//	err := db.Transaction(ctx, func(ctx context.Context, tx Transaction) error {
	//		tx.Put(ctx, &greet1)
	//		return nil
	//	})
	//
	Transaction(ctx context.Context, callback func(ctx context.Context) error) error

	// InTransaction return true if connection is in transaction
	//
	//	inTx := db.InTransaction()
	//
	InTransaction() bool

	// Connection return current connection
	//
	//	conn := db.Connection()
	//
	Connection() Connection

	// Usage return usage object
	//
	//	usage := db.Usage()
	//
	//	Usage() Usage
}

// BaseDB represent document database
//
type BaseDB struct {
	DB

	// Conn is current database connection
	//
	Conn Connection
}

// Connection return current connection
//
//	conn := db.Connection()
//
func (db *BaseDB) Connection() Connection {
	return db.Conn
}

// Close connection
//
//	db.Close()
//
func (db *BaseDB) Close() {
	if db.Conn != nil {
		db.Conn.Close()
		db.Conn = nil
	}
}

// BatchBegin put connection into batch mode. Set/Update/Delete will hold operation until CommitBatch
//
//	err := conn.BatchBegin(ctx)
//
func (db *BaseDB) BatchBegin() {
	db.Conn.BatchBegin()
}

// BatchCommit commit batch operation
//
//	err := conn.BatchCommit(ctx)
//
func (db *BaseDB) BatchCommit(ctx context.Context) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	return db.Conn.BatchCommit(ctx)
}

// InBatch return true if connection is in batch mode
//
//	inBatch := conn.InBatch()
//
func (db *BaseDB) InBatch() bool {
	return db.Conn.InBatch()
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
	return db.Conn.Transaction(ctx, callback)
}

// InTransaction return true if connection is in transaction
//
//	inTx := conn.InTransaction()
//
func (db *BaseDB) InTransaction() bool {
	return db.Conn.InTransaction()
}

// CreateNamespace create namespace, create new one if not exist
//
//	db, err := conn.CreateNamespace(ctx)
//
func (db *BaseDB) CreateNamespace(ctx context.Context) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	return db.Conn.CreateNamespace(ctx)
}

// DeleteNamespace delete namespace
//
//	err := db.DeleteNamespace(ctx)
//
func (db *BaseDB) DeleteNamespace(ctx context.Context) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	return db.Conn.DeleteNamespace(ctx)
}
