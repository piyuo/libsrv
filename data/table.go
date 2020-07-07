package data

import (
	"context"

	util "github.com/piyuo/libsrv/util"
)

// Table represent database table
//
type Table interface {
	SetConnection(conn Connection)
	SetFactory(factory func() Object)
	Factory() func() Object
	NewObject() Object
	SetTableName(tablename string)
	TableName() string
	ID() string
	Get(ctx context.Context, id string) (Object, error)
	Set(ctx context.Context, object Object) error
	Exist(ctx context.Context, id string) (bool, error)
	List(ctx context.Context) ([]Object, error)
	Select(ctx context.Context, id, field string) (interface{}, error)
	Update(ctx context.Context, id string, fields map[string]interface{}) error
	Delete(ctx context.Context, id string) error
	DeleteObject(ctx context.Context, object Object) error
	Count(ctx context.Context) (int, error)
	IsEmpty(ctx context.Context) (bool, error)
	Clear(ctx context.Context) error
	Query(ctx context.Context) Query
	Find(ctx context.Context, field, operator string, value interface{}) (Object, error)
	Search(ctx context.Context, field, operator string, value interface{}) ([]Object, error)
	Increment(ctx context.Context, id, field string, value int) error
}

// DocTable represent a table in document database
//
type DocTable struct {
	Table
	conn      Connection
	factory   func() Object
	tablename string
}

// SetConnection set connection for table
//
//	table.SetConnection(conn)
//
func (dt *DocTable) SetConnection(conn Connection) {
	dt.conn = conn
}

// SetFactory set factory function to create object
//
//	table.SetFactory(f)
//
func (dt *DocTable) SetFactory(factory func() Object) {
	dt.factory = factory
}

// Factory return factory function to create object
//
//	table.Factory()
//
func (dt *DocTable) Factory() func() Object {
	return dt.factory
}

// NewObject use factory to create new Object
//
//	obj:=table.NewObject()
//
func (dt *DocTable) NewObject() Object {
	return dt.factory()
}

// SetTableName set table name
//
//	table.SetTableName("sample")
//
func (dt *DocTable) SetTableName(tablename string) {
	dt.tablename = tablename
}

// TableName return table name
//
//	table.TableName()
//
func (dt *DocTable) TableName() string {
	return dt.tablename
}

// ID create new id for empty object
//
//
//	id := table.ID()
//
func (dt *DocTable) ID() string {
	return util.UUID()
}

// Get object by id, return nil if object is not exist
//
func (dt *DocTable) Get(ctx context.Context, id string) (Object, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	object, err := dt.conn.Get(ctx, dt.tablename, id, dt.factory)
	if err != nil {
		return nil, err
	}
	return object, nil
}

// Set object to table
//
func (dt *DocTable) Set(ctx context.Context, object Object) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	if err := dt.conn.Set(ctx, dt.tablename, object); err != nil {
		return err
	}
	return nil
}

// Exist return true if object with id exist
//
func (dt *DocTable) Exist(ctx context.Context, id string) (bool, error) {
	if ctx.Err() != nil {
		return false, ctx.Err()
	}
	return dt.conn.Exist(ctx, dt.tablename, id)
}

// List return max 10 object, if you need more! using query instead
//
func (dt *DocTable) List(ctx context.Context) ([]Object, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	return dt.conn.List(ctx, dt.tablename, dt.factory)
}

// Select return object field from data store
//
//
func (dt *DocTable) Select(ctx context.Context, id, field string) (interface{}, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	return dt.conn.Select(ctx, dt.tablename, id, field)
}

// Update partial object field without overwriting the entire document
//
//
func (dt *DocTable) Update(ctx context.Context, id string, fields map[string]interface{}) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	return dt.conn.Update(ctx, dt.tablename, id, fields)
}

// Delete object using id
//
//
func (dt *DocTable) Delete(ctx context.Context, id string) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	return dt.conn.Delete(ctx, dt.tablename, id)
}

// DeleteObject delete object
//
//
func (dt *DocTable) DeleteObject(ctx context.Context, object Object) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	return dt.conn.DeleteObject(ctx, dt.tablename, object)
}

// Clear delete all object in specific time, 1000 documents at a time, return false if still has object need to be delete
//
//
func (dt *DocTable) Clear(ctx context.Context) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	return dt.conn.Clear(ctx, dt.tablename)
}

// Query create query
//
//	query := db.Query(ctx, func() Object {
//		return new(Greet)
//	})
//
func (dt *DocTable) Query(ctx context.Context) Query {
	return dt.conn.Query(ctx, dt.tablename, dt.factory)
}

// Find return first object
//
//	exist, err := db.Available(ctx, "From", "==", "1")
//
func (dt *DocTable) Find(ctx context.Context, field, operator string, value interface{}) (Object, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	list, err := dt.Query(ctx).Where(field, operator, value).Execute(ctx)
	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return list[0], nil
	}
	return nil, nil
}

// Search return max 10 result set,cause firestore are charged for a read each time a document in the result set, we need keep result set as small as possible
//
//	count, err := db.Count(ctx,"", GreetModelName, "From", "==", "1")
//
func (dt *DocTable) Search(ctx context.Context, field, operator string, value interface{}) ([]Object, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	list, err := dt.Query(ctx).Where(field, operator, value).Execute(ctx)
	if err != nil {
		return nil, err
	}

	return list, nil
}

// Count can only return max 10 count in table, because data store charge by document count, we need use counter to get real count
//
//	count, err := db.Count(ctx,"", GreetModelName, "From", "==", "1")
//
func (dt *DocTable) Count(ctx context.Context) (int, error) {
	if ctx.Err() != nil {
		return 0, ctx.Err()
	}
	return dt.Query(ctx).Count(ctx)
}

// IsEmpty check
//
//	count, err := db.Count(ctx,"", GreetModelName, "From", "==", "1")
//
func (dt *DocTable) IsEmpty(ctx context.Context) (bool, error) {
	if ctx.Err() != nil {
		return false, ctx.Err()
	}
	return dt.Query(ctx).IsEmpty(ctx)
}

// Increment value on object field
//
//	err := db.Increment(ctx, greet.ID(), "Value", 2)
//
func (dt *DocTable) Increment(ctx context.Context, id, field string, value int) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	return dt.conn.Increment(ctx, dt.tablename, id, field, value)
}
