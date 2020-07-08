package data

import (
	"context"

	util "github.com/piyuo/libsrv/util"
)

// Table represent a table in document database
//
type Table struct {
	Connection ConnectionRef
	Factory    func() ObjectRef
	TableName  string
}

// NewObject use factory to create new Object
//
//	obj:=table.NewObject()
//
func (t *Table) NewObject() ObjectRef {
	return t.Factory()
}

// ID create new id for empty object
//
//
//	id := table.ID()
//
func (t *Table) ID() string {
	return util.UUID()
}

// Get object by id, return nil if object is not exist
//
func (t *Table) Get(ctx context.Context, id string) (ObjectRef, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	object, err := t.Connection.Get(ctx, t.TableName, id, t.Factory)
	if err != nil {
		return nil, err
	}
	return object, nil
}

// Set object to table
//
func (t *Table) Set(ctx context.Context, object ObjectRef) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	if err := t.Connection.Set(ctx, t.TableName, object); err != nil {
		return err
	}
	return nil
}

// Exist return true if object with id exist
//
func (t *Table) Exist(ctx context.Context, id string) (bool, error) {
	if ctx.Err() != nil {
		return false, ctx.Err()
	}
	return t.Connection.Exist(ctx, t.TableName, id)
}

// List return max 10 object, if you need more! using query instead
//
func (t *Table) List(ctx context.Context) ([]ObjectRef, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	return t.Connection.List(ctx, t.TableName, t.Factory)
}

// Select return object field from data store
//
//
func (t *Table) Select(ctx context.Context, id, field string) (interface{}, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	return t.Connection.Select(ctx, t.TableName, id, field)
}

// Update partial object field without overwriting the entire document
//
//
func (t *Table) Update(ctx context.Context, id string, fields map[string]interface{}) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	return t.Connection.Update(ctx, t.TableName, id, fields)
}

// Delete object using id
//
//
func (t *Table) Delete(ctx context.Context, id string) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	return t.Connection.Delete(ctx, t.TableName, id)
}

// DeleteObject delete object
//
//
func (t *Table) DeleteObject(ctx context.Context, object ObjectRef) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	return t.Connection.DeleteObject(ctx, t.TableName, object)
}

// Clear delete all object in specific time, 1000 documents at a time, return false if still has object need to be delete
//
//
func (t *Table) Clear(ctx context.Context) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	return t.Connection.Clear(ctx, t.TableName)
}

// Query create query
//
//	query := db.Query(ctx, func() Object {
//		return new(Greet)
//	})
//
func (t *Table) Query(ctx context.Context) QueryRef {
	return t.Connection.Query(ctx, t.TableName, t.Factory)
}

// Find return first object
//
//	exist, err := db.Available(ctx, "From", "==", "1")
//
func (t *Table) Find(ctx context.Context, field, operator string, value interface{}) (ObjectRef, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	list, err := t.Query(ctx).Where(field, operator, value).Execute(ctx)
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
func (t *Table) Search(ctx context.Context, field, operator string, value interface{}) ([]ObjectRef, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	list, err := t.Query(ctx).Where(field, operator, value).Execute(ctx)
	if err != nil {
		return nil, err
	}

	return list, nil
}

// Count can only return max 10 count in table, because data store charge by document count, we need use counter to get real count
//
//	count, err := db.Count(ctx,"", GreetModelName, "From", "==", "1")
//
func (t *Table) Count(ctx context.Context) (int, error) {
	if ctx.Err() != nil {
		return 0, ctx.Err()
	}
	return t.Query(ctx).Count(ctx)
}

// IsEmpty check
//
//	count, err := db.Count(ctx,"", GreetModelName, "From", "==", "1")
//
func (t *Table) IsEmpty(ctx context.Context) (bool, error) {
	if ctx.Err() != nil {
		return false, ctx.Err()
	}
	return t.Query(ctx).IsEmpty(ctx)
}

// Increment value on object field
//
//	err := db.Increment(ctx, greet.ID(), "Value", 2)
//
func (t *Table) Increment(ctx context.Context, id, field string, value int) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	return t.Connection.Increment(ctx, t.TableName, id, field, value)
}
