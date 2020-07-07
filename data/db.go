package data

import (
	"context"
)

// DB represent database
//
type DB interface {
	// Connection return db connection
	//
	//	conn.DB()
	//
	Connection() Connection

	// Close connection
	//
	//	conn.Close()
	//
	Close()

	// SetConnection set current db connection
	//
	//	db := &GlobalDB{}
	//	db.SetConnection(conn)
	//
	SetConnection(connection Connection)

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
}

// DocDB represent document database
//
type DocDB struct {
	connection Connection
}

// Connection return database connection
//
//	db.Connection()
//
func (db *DocDB) Connection() Connection {
	return db.connection
}

// SetConnection set database connection
//
//	conn.SetConnection(donn)
//
func (db *DocDB) SetConnection(connection Connection) {
	db.connection = connection
}

// Close connection
//
//	conn.Close()
//
func (db *DocDB) Close() {
	if db.connection != nil {
		db.connection.Close()
		db.connection = nil
	}
}

// Transaction start a transaction
//
//	err := conn.Transaction(ctx, func(ctx context.Context) error {
//		accounts.Set(ctx, &greet1)
//		return nil
//	})
//
func (db *DocDB) Transaction(ctx context.Context, callback func(ctx context.Context) error) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	return db.connection.Transaction(ctx, callback)
}

// CreateNamespace create namespace, create new one if not exist
//
//	dbRef, err := conn.CreateNamespace(ctx)
//
func (db *DocDB) CreateNamespace(ctx context.Context) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	return db.connection.CreateNamespace(ctx)
}

// DeleteNamespace delete namespace
//
//	err := db.DeleteNamespace(ctx)
//
func (db *DocDB) DeleteNamespace(ctx context.Context) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	return db.connection.DeleteNamespace(ctx)
}
