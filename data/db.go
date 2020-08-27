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
	//	c.Close()
	//
	Close()

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
	//	err := c.Transaction(ctx, func(ctx context.Context, tx Transaction) error {
	//		tx.Put(ctx, &greet1)
	//		return nil
	//	})
	//
	Transaction(ctx context.Context, callback func(ctx context.Context) error) error

	// InTransaction return true if connection is in transaction
	//
	//	inTx := c.InTransaction()
	//
	InTransaction() bool

	// GetConnection return current connection
	//
	//	conn := c.GetConnection()
	//
	GetConnection() Connection

	// Usage return usage object
	//
	//	usage := c.Usage()
	//
	//	Usage() Usage
}

// BaseDB represent document database
//
type BaseDB struct {
	DB

	// Conn is current database connection
	//
	Connection Connection
}

// GetConnection return current connection
//
//	conn := c.GetConnection()
//
func (c *BaseDB) GetConnection() Connection {
	return c.Connection
}

// Close connection
//
//	c.Close()
//
func (c *BaseDB) Close() {
	if c.Connection != nil {
		c.Connection.Close()
		c.Connection = nil
	}
}

// BatchBegin put connection into batch mode. Set/Update/Delete will hold operation until CommitBatch
//
//	err := conn.BatchBegin(ctx)
//
func (c *BaseDB) BatchBegin() {
	c.Connection.BatchBegin()
}

// BatchCommit commit batch operation
//
//	err := conn.BatchCommit(ctx)
//
func (c *BaseDB) BatchCommit(ctx context.Context) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	return c.Connection.BatchCommit(ctx)
}

// InBatch return true if connection is in batch mode
//
//	inBatch := conn.InBatch()
//
func (c *BaseDB) InBatch() bool {
	return c.Connection.InBatch()
}

// Transaction start a transaction
//
//	err := conn.Transaction(ctx, func(ctx context.Context) error {
//		accounts.Set(ctx, &greet1)
//		return nil
//	})
//
func (c *BaseDB) Transaction(ctx context.Context, callback func(ctx context.Context) error) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	return c.Connection.Transaction(ctx, callback)
}

// InTransaction return true if connection is in transaction
//
//	inTx := conn.InTransaction()
//
func (c *BaseDB) InTransaction() bool {
	return c.Connection.InTransaction()
}
