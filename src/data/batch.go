package data

import (
	"context"
)

// Batch define batch operation
//
type Batch interface {

	// Set object into data store, If the document does not exist, it will be created. If the document does exist, its contents will be overwritten with the newly provided data,
	//
	// if object does not have id, it will created using UUID
	//
	//	if err := conn.Set(ctx, tablename, object); err != nil {
	//		return err
	//	}
	//
	Set(ctx context.Context, tablename string, object Object) error

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

	// Increment value on object field, return error if object does not exist
	//
	//	err := conn.Increment(ctx,"", GreetModelName, greet.ID(), "Value", 2)
	//
	Increment(ctx context.Context, tablename, id, field string, value int) error

	// Begin batch operation Set/Update/Delete will hold operation until CommitBatch
	//
	//	err := conn.BatchBegin()
	//
	Begin()

	// InBatch return true if connection is in batch mode
	//
	//	inBatch := conn.InBatch()
	//
	InBatch() bool

	// BatchCommit commit batch operation
	//
	//	err := conn.BatchCommit(ctx)
	//
	Commit(ctx context.Context) error
}
