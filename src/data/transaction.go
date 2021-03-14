package data

import (
	"context"
)

// Transaction define transaction operation
//
type Transaction interface {

	// Get data object from data store, return nil if object does not exist
	//
	//	object, err := Get(ctx, &Sample{}, "id")
	//
	Get(ctx context.Context, obj Object, id string) (Object, error)

	// Set object into data store, If the document does not exist, it will be created. If the document does exist, its contents will be overwritten with the newly provided data,
	// if object does not have id, it will created using UUID
	//
	//	object, err := Get(ctx, &Sample{}, "id")
	//
	Set(ctx context.Context, obj Object) error

	// Exists return true if object with id exist
	//
	//	found,err := Exists(ctx, &Sample{}, "id")
	//
	Exists(ctx context.Context, obj Object, id string) (bool, error)

	// All return max 10 object, if you need more! using query instead
	//
	//	list,err := All(ctx, &Sample{})
	//
	All(ctx context.Context, obj Object) ([]Object, error)

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

	// DeleteBatch delete list of id, no error if id not exist
	//
	//	conn.DeleteBatch(ctx, dt.tablename, ids)
	//
	DeleteBatch(ctx context.Context, tablename string, ids []string) error

	// Clear delete all object in specific time, 500 documents at a time, return false if still has object need to be delete
	//	if in transaction , only 500 documents can be delete
	//
	//	err := conn.Clear(ctx, tablename)
	//
	Clear(ctx context.Context, tablename string) error

	// Begin a transaction
	//
	//	err := transaction.Begin(ctx, func(ctx context.Context) error {
	//		return nil
	//	})
	//
	Begin(ctx context.Context, callback func(ctx context.Context) error) error

	// IsBegin return true if connection is in transaction
	//
	//	begin := conn.InTransaction()
	//
	IsBegin() bool

	// Increment value on object field, return error if object does not exist
	//
	//	err := conn.Increment(ctx,"", GreetModelName, greet.ID(), "Value", 2)
	//
	Increment(ctx context.Context, tablename, id, field string, value int) error
}
