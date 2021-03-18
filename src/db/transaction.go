package db

import (
	"context"
)

// TransactionFunc define a transaction function
//
type TransactionFunc func(ctx context.Context, tx Transaction) error

// Transaction define transaction operation
//
type Transaction interface {

	// Get data object from table, return nil if object does not exist
	//
	//	object, err := Get(ctx, &Sample{}, "id")
	//
	Get(ctx context.Context, obj Object, id string) (Object, error)

	// Exists return true if object with id exist
	//
	//	found,err := Exists(ctx, &Sample{}, "id")
	//
	Exists(ctx context.Context, obj Object, id string) (bool, error)

	// List return object list, use max to specific return object count
	//
	//	list,err := List(ctx, &Sample{},10)
	//
	List(ctx context.Context, obj Object, max int) ([]Object, error)

	// Select return object field from data store, return nil if object does not exist
	//
	//	return Select(ctx, &Sample{}, id, field)
	//
	Select(ctx context.Context, obj Object, id, field string) (interface{}, error)

	// Query create query
	//
	//	c.Query(ctx, &Sample{}).Execute(ctx)
	//
	Query(obj Object) Query

	// Set object into table, If the document not exist, it will be created. If the document does exist, its contents will be overwritten with the newly provided data, if object does not have id, it will created using UUID
	//
	//	 err := Set(ctx, object)
	//
	Set(ctx context.Context, obj Object) error

	// Update partial object field, create new one if object does not exist, this function is significant slow than Set()
	//
	//	err = Update(ctx, Sample, map[string]interface{}{
	//		"desc": "hi",
	//	})
	//
	Update(ctx context.Context, obj Object, fields map[string]interface{}) error

	// Increment value on object field, return error if object does not exist
	//
	//	err := Increment(ctx,sample, "Value", 2)
	//
	Increment(ctx context.Context, obj Object, field string, value int) error

	// Delete object, no error if id not exist
	//
	//	Delete(ctx, sample)
	//
	Delete(ctx context.Context, obj Object) error

	// Clear delete all document in collection. try 10 batch each batch only delete document specific in max. return true if collection is cleared
	//
	//	cleared, err := Clear(ctx, &Sample{}, 50)
	//
	Clear(ctx context.Context, obj Object, max int) (bool, error)
}
