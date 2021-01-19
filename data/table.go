package data

import (
	"context"
	"time"

	"github.com/piyuo/libsrv/env"
	"github.com/piyuo/libsrv/identifier"
	"github.com/pkg/errors"
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
	Connection Connection
	Factory    func() Object
	TableName  string
}

// NewObject use factory to create new Object
//
//	obj:=table.NewObject()
//	account:=obj.(*Account)
//
func (c *Table) NewObject() Object {
	return c.Factory()
}

// UUID is a help function help you create UUID
//
//	id := table.UUID()
//
func (c *Table) UUID() string {
	return identifier.UUID()
}

// Get object by id, return nil if object is not exist
//
//	object, err := table.Get(ctx, sample.ID)
//	if object != nil && err == nil{
//		sample := object.(*Sample)
//	}
//
func (c *Table) Get(ctx context.Context, id string) (Object, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	object, err := c.Connection.Get(ctx, c.TableName, id, c.Factory)
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
func (c *Table) Set(ctx context.Context, object Object) error {
	if object == nil {
		return errors.New("object can not be nil: " + c.TableName)
	}

	if ctx.Err() != nil {
		return ctx.Err()
	}

	t := time.Now().UTC()
	object.SetCreateTime(t) // create time will not change if it's not empty
	object.SetUpdateTime(t)
	object.SetAccountID(env.GetAccountID(ctx))
	object.SetUserID(env.GetUserID(ctx))

	if err := c.Connection.Set(ctx, c.TableName, object); err != nil {
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
func (c *Table) Exist(ctx context.Context, id string) (bool, error) {
	if ctx.Err() != nil {
		return false, ctx.Err()
	}
	return c.Connection.Exist(ctx, c.TableName, id)
}

// All return objects in table, max 10 object, if you need more! using query instead
//
//	list, err := table.All(ctx)
//	So(len(list), ShouldEqual, 2)
//	So(list[0].(*Sample).Name, ShouldStartWith, "sample")
//	So(list[1].(*Sample).Name, ShouldStartWith, "sample")
//
func (c *Table) All(ctx context.Context) ([]Object, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	return c.Connection.All(ctx, c.TableName, c.Factory)
}

// Select return object specific field from database
//
//	name, err := table.Select(ctx, sample.ID, "Name")
//	So(name, ShouldEqual, "sample")
//
func (c *Table) Select(ctx context.Context, id, field string) (interface{}, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	return c.Connection.Select(ctx, c.TableName, id, field)
}

// Update partial object field without overwriting the entire document
//
//	err = table.Update(ctx, sample.ID, map[string]interface{}{
//		"Name":  "sample2",
//		"Value": 2,
//	})
//
func (c *Table) Update(ctx context.Context, id string, fields map[string]interface{}) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	return c.Connection.Update(ctx, c.TableName, id, fields)
}

// Delete object using id
//
//	err = table.Delete(ctx, "12345")
//
func (c *Table) Delete(ctx context.Context, id string) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	return c.Connection.Delete(ctx, c.TableName, id)
}

// DeleteObject delete object
//
//	err = table.DeleteObject(ctx, sample)
//
func (c *Table) DeleteObject(ctx context.Context, object Object) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	return c.Connection.DeleteObject(ctx, c.TableName, object)
}

// DeleteBatch delete list of id use batch mode, no error if id not exist
//
//	table.DeleteBatch(ctx, ids)
//
func (c *Table) DeleteBatch(ctx context.Context, ids []string) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	return c.Connection.DeleteBatch(ctx, c.TableName, ids)
}

// Clear delete all object in specific time, 500 documents at a time, if in transaction , only 10 documents can be delete
//
//	err = table.Clear(ctx)
//
func (c *Table) Clear(ctx context.Context) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	return c.Connection.Clear(ctx, c.TableName)
}

// Query create a query
//
//	list, err = table.Query().OrderBy("Name").Execute(ctx)
//	So(len(list), ShouldEqual, 2)
//	So(list[0].(*Sample).Name, ShouldEqual, sample1.Name)
//	So(list[1].(*Sample).Name, ShouldEqual, sample2.Name)
//
func (c *Table) Query() Query {
	return c.Connection.Query(c.TableName, c.Factory)
}

// Find return first object in table, return nil if not found
//
//	obj, err = table.Find(ctx, "Value", "==", 1)
//	So((obj.(*Sample)).Name, ShouldEqual, "sample")
//
func (c *Table) Find(ctx context.Context, field, operator string, value interface{}) (Object, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	list, err := c.Query().Where(field, operator, value).Execute(ctx)
	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return list[0], nil
	}
	return nil, nil
}

// GetFirstObject return first object in table, return nil if not found
//
//	obj, err = table.GetFirstObject(ctx)
//	So((obj.(*Sample)).Name, ShouldEqual, "sample")
//
func (c *Table) GetFirstObject(ctx context.Context) (Object, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	return c.Query().GetFirstObject(ctx)
}

// GetFirstID return first id in table, return nil if not found
//
//	id, err = table.GetFirstID(ctx)
//
func (c *Table) GetFirstID(ctx context.Context) (string, error) {
	if ctx.Err() != nil {
		return "", ctx.Err()
	}
	return c.Query().GetFirstID(ctx)
}

// List return max 10 result set,cause firestore are charged for a read each time a document in the result set, we need keep result set as small as possible, if you need more please use query
//
//	list, err := table.Search(ctx, "Name", "==", "sample")
//	So(len(list), ShouldEqual, 1)
//
func (c *Table) List(ctx context.Context, field, operator string, value interface{}) ([]Object, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	query := c.Query().Where(field, operator, value).Limit(10)
	list, err := query.Execute(ctx)
	if err != nil {
		return nil, err
	}
	return list, nil
}

// SortList work like list but add sort function, please be aware you need create index on firestore in order to sort
//
//	list, err := table.SortList(ctx, "Name", "==", "sample","Value",ASC)
//	So(len(list), ShouldEqual, 1)
//
func (c *Table) SortList(ctx context.Context, field, operator string, value interface{}, orderby string, order orderby) ([]Object, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	query := c.Query().Where(field, operator, value).Limit(10)
	if orderby != "" {
		if order == ASC {
			query = query.OrderBy(orderby)
		} else {
			query = query.OrderByDesc(orderby)
		}
	}
	list, err := query.Execute(ctx)
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
func (c *Table) Count(ctx context.Context) (int, error) {
	if ctx.Err() != nil {
		return 0, ctx.Err()
	}
	return c.Query().Count(ctx)
}

// IsEmpty return true if table is empty
//
//	empty, err := table.IsEmpty(ctx)
//	So(empty, ShouldEqual, false)
//
func (c *Table) IsEmpty(ctx context.Context) (bool, error) {
	if ctx.Err() != nil {
		return false, ctx.Err()
	}
	return c.Query().IsEmpty(ctx)
}

// Increment value on object field
//
//	err = table.Increment(ctx, sample.ID, "Value", 1)
//
func (c *Table) Increment(ctx context.Context, id, field string, value int) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	return c.Connection.Increment(ctx, c.TableName, id, field, value)
}
