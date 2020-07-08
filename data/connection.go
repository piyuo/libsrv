package data

import (
	"context"
)

const limitQueryDefault = 10
const limitTransactionClear = 10
const limitClear = 500

// Connection define how to connect and manipulate database
//
type Connection interface {
	// Close database connection
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

	// Get data object from data store, return nil if object does not exist
	//
	//	object, err := conn.Get(ctx, tablename, id, factory)
	//
	Get(ctx context.Context, tablename, id string, factory func() Object) (Object, error)

	// Set object into data store, If the document does not exist, it will be created. If the document does exist, its contents will be overwritten with the newly provided data,
	//
	// if object does not have id, it will created using UUID
	//
	//	if err := conn.Set(ctx, tablename, object); err != nil {
	//		return err
	//	}
	//
	Set(ctx context.Context, tablename string, obj Object) error

	// Exist return true if object with id exist
	//
	//	return conn.Exist(ctx, tablename, id)
	//
	Exist(ctx context.Context, tablename, id string) (bool, error)

	// List return max 10 object, if you need more! using query instead
	//
	//	return conn.List(ctx, tablename, factory)
	//
	List(ctx context.Context, tablename string, factory func() Object) ([]Object, error)

	// Select return object field from data store, return nil if object does not exist
	//
	//	return conn.Select(ctx, tablename, id, field)
	//
	Select(ctx context.Context, tablename, id, field string) (interface{}, error)

	// Update partial object field, create new one if object does not exist,  this function is significant slow than Set()
	//
	//	err = conn.Update(ctx, tablename, greet.ID(), map[string]interface{}{
	//		"Description": "helloworld",
	//	})
	//
	Update(ctx context.Context, tablename, id string, fields map[string]interface{}) error

	// Delete object using table name and id, no error if id not exist
	//
	//	conn.Delete(ctx, tablename, id)
	//
	Delete(ctx context.Context, tablename, id string) error

	// DeleteObject delete object, no error if id not exist
	//
	//	conn.DeleteObject(ctx, dt.tablename, object)
	//
	DeleteObject(ctx context.Context, tablename string, obj Object) error

	// Clear delete all object in specific time, 500 documents at a time, return false if still has object need to be delete
	//	if in transaction , only 500 documents can be delete
	//
	//	err := conn.Clear(ctx, tablename)
	//
	Clear(ctx context.Context, tablename string) error

	// Query create query
	//
	//	conn.Query(ctx, tablename, factory)
	//
	Query(ctx context.Context, tablename string, factory func() Object) Query

	// Transaction start a transaction
	//
	//	err := conn.Transaction(ctx, func(ctx context.Context) error {
	//		return nil
	//	})
	//
	Transaction(ctx context.Context, callback func(ctx context.Context) error) error

	// Increment value on object field, return error if object does not exist
	//
	//	err := conn.Increment(ctx,"", GreetModelName, greet.ID(), "Value", 2)
	//
	Increment(ctx context.Context, tablename, id, field string, value int) error

	// Counter return counter from data store, create one if not exist
	//
	//	counter,err = conn.Counter(ctx, tablename, countername, numshards)
	//
	Counter(ctx context.Context, tablename, countername string, numShards int) (Counter, error)
}
