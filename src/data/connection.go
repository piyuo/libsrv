package data

import (
	"context"
)

// Connection define how to connect and manipulate database
//
type Connection interface {
	// Close database connection
	//
	//	conn.Close()
	//
	Close()

	// Get data object from data store, return nil if object does not exist
	//
	//	object, err := Get(ctx, &Sample{}, "id")
	//
	Get(ctx context.Context, obj Object, id string) (Object, error)

	// Set object into data store, If the document does not exist, it will be created. If the document does exist, its contents will be overwritten with the newly provided data,
	// if object does not have id, it will created using UUID
	//
	//	 err := Set(ctx, object)
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

	// Query create query
	//
	//	conn.Query(ctx, tablename, factory)
	//
	Query(tablename string, factory func() Object) Query

	// Increment value on object field, return error if object does not exist
	//
	//	err := conn.Increment(ctx,"", GreetModelName, greet.ID(), "Value", 2)
	//
	Increment(ctx context.Context, tablename, id, field string, value int) error

	// CreateTransaction create transaction
	//
	CreateTransaction() Transaction

	// CreateBatch create batch
	//
	CreateBatch() Batch

	// CreateCoder return coder from database, set numshards 100 times of concurrent usage. for example if you think concurrent use is 10/seconds then set numshards to 1000 to avoid too much retention error
	//
	//	productCoder,err = conn.CreateCoder("tableName","coderName",100)
	//
	CreateCoder(tableName, coderName string, numshards int) Coder

	// Counter return counter from database, create one if not exist, set numshards 100 times of concurrent usage. for example if you think concurrent use is 10/seconds then set numshards to 1000 to avoid too much retention error
	// if keepDateHierarchy is true, counter will automatically generate year/month/day/hour hierarchy in utc timezone
	//
	//	orderCountCounter,err = conn.CreateCounter("tableName","coderName",100,true)
	//
	CreateCounter(tableName, counterName string, numshards int, hierarchy DateHierarchy) Counter

	// Serial return serial from database, create one if not exist, please be aware Serial can only generate 1 number per second, use serial with high frequency will cause too much retention error
	//
	//	productNo,err = conn.CreateSerial("tableName","serialName")
	//
	CreateSerial(tableName, serialName string) Serial
}
