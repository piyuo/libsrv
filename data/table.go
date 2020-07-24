package data

import (
	"context"

	identifier "github.com/piyuo/libsrv/identifier"
)

// Table represent collection of document in document database, you can do operation like get/set/query on documents
//
//	return &data.Table{
//		Connection: db.Connection,
//		TableName:  "account",
//		Factory:    func() data.Object {
//			return &Account{}
//		},
//	}
//
type Table struct {
	CurrentConnection Connection
	Factory           func() Object
	TableName         string
}

// NewObject use factory to create new Object
//
//	obj:=table.NewObject()
//	account:=obj.(*Account)
//
func (t *Table) NewObject() Object {
	return t.Factory()
}

// UUID is a help function help you create UUID
//
//	id := table.UUID()
//
func (t *Table) UUID() string {
	return identifier.UUID()
}

// Get object by id, return nil if object is not exist
//
//	object, err := table.Get(ctx, sample.ID)
//	if object != nil && err == nil{
//		sample := object.(*Sample)
//	}
//
func (t *Table) Get(ctx context.Context, id string) (Object, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	object, err := t.CurrentConnection.Get(ctx, t.TableName, id, t.Factory)
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
func (t *Table) Set(ctx context.Context, object Object) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	if err := t.CurrentConnection.Set(ctx, t.TableName, object); err != nil {
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
	return t.CurrentConnection.Exist(ctx, t.TableName, id)
}

// List return objects in table, max 10 object, if you need more! using query instead
//
//	list, err := table.List(ctx)
//	So(len(list), ShouldEqual, 2)
//	So(list[0].(*Sample).Name, ShouldStartWith, "sample")
//	So(list[1].(*Sample).Name, ShouldStartWith, "sample")
//
func (t *Table) List(ctx context.Context) ([]Object, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	return t.CurrentConnection.List(ctx, t.TableName, t.Factory)
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
	return t.CurrentConnection.Select(ctx, t.TableName, id, field)
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
	return t.CurrentConnection.Update(ctx, t.TableName, id, fields)
}

// Delete object using id
//
//	err = table.Delete(ctx, "12345")
//
func (t *Table) Delete(ctx context.Context, id string) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	return t.CurrentConnection.Delete(ctx, t.TableName, id)
}

// DeleteObject delete object
//
//	err = table.DeleteObject(ctx, sample)
//
func (t *Table) DeleteObject(ctx context.Context, object Object) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	return t.CurrentConnection.DeleteObject(ctx, t.TableName, object)
}

// Clear delete all object in specific time, 500 documents at a time, if in transaction , only 10 documents can be delete
//
//	err = table.Clear(ctx)
//
func (t *Table) Clear(ctx context.Context) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	return t.CurrentConnection.Clear(ctx, t.TableName)
}

// Query create a query
//
//	list, err = table.Query().OrderBy("Name").Execute(ctx)
//	So(len(list), ShouldEqual, 2)
//	So(list[0].(*Sample).Name, ShouldEqual, sample1.Name)
//	So(list[1].(*Sample).Name, ShouldEqual, sample2.Name)
//
func (t *Table) Query() Query {
	return t.CurrentConnection.Query(t.TableName, t.Factory)
}

// Find return first object in table
//
//	obj, err = table.Find(ctx, "Value", "==", 1)
//	So((obj.(*Sample)).Name, ShouldEqual, "sample")
//
func (t *Table) Find(ctx context.Context, field, operator string, value interface{}) (Object, error) {
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

// Search return max 10 result set,cause firestore are charged for a read each time a document in the result set, we need keep result set as small as possible, if you need more please use query
//
//	objects, err := table.Search(ctx, "Name", "==", "sample")
//	So(len(objects), ShouldEqual, 1)
//
func (t *Table) Search(ctx context.Context, field, operator string, value interface{}) ([]Object, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	list, err := t.Query().Where(field, operator, value).Limit(10).Execute(ctx)
	if err != nil {
		return nil, err
	}
	return list, nil
}

// Count can only count max 10 in table, because data store charge by document count, you need use counter keep large count
//
//	count, err := table.Count(ctx)
//	So(count, ShouldEqual, 1)
//
func (t *Table) Count(ctx context.Context) (int, error) {
	if ctx.Err() != nil {
		return 0, ctx.Err()
	}
	return t.Query().Count(ctx)
}

// IsEmpty return true if table is empty
//
//	empty, err := table.IsEmpty(ctx)
//	So(empty, ShouldEqual, false)
//
func (t *Table) IsEmpty(ctx context.Context) (bool, error) {
	if ctx.Err() != nil {
		return false, ctx.Err()
	}
	return t.Query().IsEmpty(ctx)
}

// Increment value on object field
//
//	err = table.Increment(ctx, sample.ID, "Value", 1)
//
func (t *Table) Increment(ctx context.Context, id, field string, value int) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	return t.CurrentConnection.Increment(ctx, t.TableName, id, field, value)
}
