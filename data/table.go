package data

import (
	"context"

	util "github.com/piyuo/libsrv/util"
)

// Table represent collection of document in document database, you can do operation like get/set/query on documents
//
//	return &data.Table{
//		Connection: db.Connection,
//		TableName:  "account",
//		Factory:    func() data.ObjectRef {
//			return &Account{}
//		},
//	}
//
type Table struct {
	Connection ConnectionRef
	Factory    func() ObjectRef
	TableName  string
}

// NewObject use factory to create new Object
//
//	obj:=table.NewObject()
//	account:=obj.(*Account)
//
func (t *Table) NewObject() ObjectRef {
	return t.Factory()
}

// UUID is a help function help you create UUID
//
//	id := table.UUID()
//
func (t *Table) UUID() string {
	return util.UUID()
}

// Get object by id, return nil if object is not exist
//
//	object, err := table.Get(ctx, sample.ID)
//	if object != nil && err == nil{
//		sample := object.(*Sample)
//	}
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

// Set save object to database, if document with same id exist it will be overwrite
//
//	sample := &Sample{
//		Name:  "sample",
//		Value: 1,
//	}
//	err = table.Set(ctx, sample)
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
//	exist, err := table.Exist(ctx, sample.ID)
//	if exist {
//		fmt.Printf("object exist")
//	}
//
func (t *Table) Exist(ctx context.Context, id string) (bool, error) {
	if ctx.Err() != nil {
		return false, ctx.Err()
	}
	return t.Connection.Exist(ctx, t.TableName, id)
}

// List return objects in table, max 10 object, if you need more! using query instead
//
//	list, err := table.List(ctx)
//	So(len(list), ShouldEqual, 2)
//	So(list[0].(*Sample).Name, ShouldStartWith, "sample")
//	So(list[1].(*Sample).Name, ShouldStartWith, "sample")
//
func (t *Table) List(ctx context.Context) ([]ObjectRef, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	return t.Connection.List(ctx, t.TableName, t.Factory)
}

// Select return object specific field from database
//
//	name, err := table.Select(ctx, sample.ID, "Name")
//	So(name, ShouldEqual, "sample")
//
func (t *Table) Select(ctx context.Context, id, field string) (interface{}, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	return t.Connection.Select(ctx, t.TableName, id, field)
}

// Update partial object field without overwriting the entire document
//
//	err = table.Update(ctx, sample.ID, map[string]interface{}{
//		"Name":  "sample2",
//		"Value": 2,
//	})
//
func (t *Table) Update(ctx context.Context, id string, fields map[string]interface{}) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	return t.Connection.Update(ctx, t.TableName, id, fields)
}

// Delete object using id
//
//	err = table.Delete(ctx, "12345")
//
func (t *Table) Delete(ctx context.Context, id string) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	return t.Connection.Delete(ctx, t.TableName, id)
}

// DeleteObject delete object
//
//	err = table.DeleteObject(ctx, sample)
//
func (t *Table) DeleteObject(ctx context.Context, object ObjectRef) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	return t.Connection.DeleteObject(ctx, t.TableName, object)
}

// Clear delete all object in specific time, 500 documents at a time, if in transaction , only 10 documents can be delete
//
//	err = table.Clear(ctx)
//
func (t *Table) Clear(ctx context.Context) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	return t.Connection.Clear(ctx, t.TableName)
}

// Query create a query
//
//	list, err = table.Query().OrderBy("Name").Execute(ctx)
//	So(len(list), ShouldEqual, 2)
//	So(list[0].(*Sample).Name, ShouldEqual, sample1.Name)
//	So(list[1].(*Sample).Name, ShouldEqual, sample2.Name)
//
func (t *Table) Query() QueryRef {
	return t.Connection.Query(t.TableName, t.Factory)
}

// Find return first object
//
//	exist, err := db.Available(ctx, "From", "==", "1")
//
func (t *Table) Find(ctx context.Context, field, operator string, value interface{}) (ObjectRef, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	list, err := t.Query().Where(field, operator, value).Execute(ctx)
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

	list, err := t.Query().Where(field, operator, value).Execute(ctx)
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
	return t.Query().Count(ctx)
}

// IsEmpty check
//
//	count, err := db.Count(ctx,"", GreetModelName, "From", "==", "1")
//
func (t *Table) IsEmpty(ctx context.Context) (bool, error) {
	if ctx.Err() != nil {
		return false, ctx.Err()
	}
	return t.Query().IsEmpty(ctx)
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
