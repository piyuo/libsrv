package gstore

import (
	"context"
	"fmt"
	"strconv"

	"cloud.google.com/go/firestore"
	"github.com/piyuo/libsrv/src/data"
	"github.com/piyuo/libsrv/src/identifier"
	"github.com/pkg/errors"
	"google.golang.org/api/iterator"
)

// TransactionFirestore implement firestore transaction
//
type TransactionFirestore struct {
	data.Transaction

	// conn is firestore connection
	//
	conn *ConnectionFirestore

	//tx is curenet transacton, it is nil if not in transaction
	//
	tx *firestore.Transaction
}

// Begin a transaction
//
//	err := transaction.Begin(ctx, func(ctx context.Context) error {
//		return nil
//	})
//
func (c *TransactionFirestore) Begin(ctx context.Context, callback func(ctx context.Context) error) error {
	var stopTransaction = func() {
		c.tx = nil
	}

	return c.conn.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		c.tx = tx
		defer stopTransaction()
		return callback(ctx)
	})
}

// IsBegin return true if connection is in transaction
//
//	begin := c.IsBegin()
//
func (c *TransactionFirestore) IsBegin() bool {
	return c.tx != nil
}

// Get data object from table, return nil if object does not exist
//
//	object, err := Get(ctx, &Sample{}, "id")
//
func (c *TransactionFirestore) Get(ctx context.Context, obj data.Object, id string) (data.Object, error) {
	if obj == nil {
		return nil, errors.New(fmt.Sprintf("obj must not nil %v", id))
	}
	if id == "" {
		return nil, errors.New(fmt.Sprintf("id must not empty %v", obj.TableName()))
	}
	docRef := c.conn.getDocRef(obj.TableName(), id)
	snapshot, err := c.tx.Get(docRef)
	return snapshotToObject(obj, id, docRef, snapshot, err)
}

// Set object into table, If the document does not exist, it will be created. If the document does exist, its contents will be overwritten with the newly provided data, if object does not have id, it will created using UUID
//
//	 err := Set(ctx, object)
//
func (c *TransactionFirestore) Set(ctx context.Context, obj data.Object) error {
	if obj == nil {
		return errors.New("Set() obj must not nil")
	}
	var docRef *firestore.DocumentRef
	if obj.Ref() == nil { // new object
		if obj.ID() == "" {
			obj.SetID(identifier.UUID())
		}
		docRef = c.conn.getDocRef(obj.TableName(), obj.ID())
		obj.SetRef(docRef)
	} else { // object already exist
		docRef = obj.Ref().(*firestore.DocumentRef)
	}

	err := c.tx.Set(docRef, obj)

	if err != nil {
		return errors.Wrapf(err, "Set(%v,%v)", errorID(obj.TableName(), obj.ID()))
	}
	return nil
}

// Exists return true if object with id exist
//
//	found,err := Exists(ctx, &Sample{}, "id")
//
func (c *TransactionFirestore) Exists(ctx context.Context, obj data.Object, id string) (bool, error) {
	if obj == nil {
		return false, errors.New("Exists() obj must not nil")
	}
	if id == "" {
		return false, errors.New(fmt.Sprintf("Exists(%v) id must not empty", obj.TableName()))
	}
	docRef := c.conn.getDocRef(obj.TableName(), id)
	snapshot, err := c.tx.Get(docRef)

	if snapshot != nil && !snapshot.Exists() {
		return false, nil
	}
	if err != nil {
		return false, errors.Wrap(err, "failed to get document: "+errorID(tablename, id))
	}
	return true, nil
}

// All return max 10 object, if you need more! using query instead
//
//	list,err := All(ctx, &Sample{})
//
func (c *TransactionFirestore) All(ctx context.Context, obj data.Object) ([]data.Object, error) {
	collectionRef := c.conn.getCollectionRef(obj.TableName())
	list := []data.Object{}
	iter := c.tx.Documents(collectionRef.Query.Limit(data.LimitQueryDefault))
	defer iter.Stop()

	for {
		snapshot, err := iter.Next()
		if err == iterator.Done {
			break
		} else if err != nil {
			return nil, errors.Wrapf(err, "iter(%v)", obj.TableName())
		}
		object := obj.Factory()
		if object == nil {
			return nil, errors.New(fmt.Sprint("%v not implement Factory()", obj.TableName()))
		}

		_, err = snapshotToObject(object, snapshot.Ref.ID, snapshot.Ref, snapshot, err)
		if err != nil {
			return nil, errors.WithMessagef(err, "iter snapshotToObject(%v,%v)", obj.TableName(), snapshot.Ref.ID)
		}
		list = append(list, object)
	}
	return list, nil
}

// Select return object field from data store, return nil if object does not exist
//
//	return c.Select(ctx, "sample", "sample-id", "Name")
//
func (c *TransactionFirestore) Select(ctx context.Context, tablename, id, field string) (interface{}, error) {
	docRef := c.conn.getDocRef(tablename, id)
	snapshot, err := c.tx.Get(docRef)

	if snapshot != nil && !snapshot.Exists() {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to get document: "+errorID(tablename, id))
	}
	value, err := snapshot.DataAt(field)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get value from document: "+errorID(tablename, id))
	}
	return value, nil
}

// Update partial object field, create new one if object does not exist,  this function is significant slow than Set()
//
//	err = c.Update(ctx, "sample", "sample-id", map[string]interface{}{
//		"Name": "helloworld",
//	})
//
func (c *TransactionFirestore) Update(ctx context.Context, tablename, id string, fields map[string]interface{}) error {
	docRef := c.conn.getDocRef(tablename, id)
	err := c.tx.Set(docRef, fields, firestore.MergeAll)
	if err != nil {
		return errors.Wrap(err, "failed to update field in transaction: "+errorID(tablename, id))
	}
	return nil
}

// Increment value on object field, return error if object does not exist
//
//	err := c.Increment(ctx,"sample", "sample-id", "Value", 1)
//
func (c *TransactionFirestore) Increment(ctx context.Context, tablename, id, field string, value int) error {
	docRef := c.conn.getDocRef(tablename, id)
	err := c.tx.Update(docRef, []firestore.Update{
		{Path: field, Value: firestore.Increment(value)},
	})
	if err != nil {
		return errors.Wrap(err, "failed to increment "+field+" with "+strconv.Itoa(value)+" in transaction: "+errorID(tablename, id))
	}
	return nil
}

// Delete object using table name and id, no error if id did not exist
//
//	c.Delete(ctx, "sample", "sample-id")
//
func (c *TransactionFirestore) Delete(ctx context.Context, tablename, id string) error {
	docRef := c.conn.getDocRef(tablename, id)
	err := c.tx.Delete(docRef)
	if err != nil {
		return errors.Wrap(err, "failed to delete in transaction: "+errorID(tablename, id))
	}
	return nil
}

// DeleteObject delete object, no error if id did not exist
//
//	c.DeleteObject(ctx, "sample", object)
//
func (c *TransactionFirestore) DeleteObject(ctx context.Context, tablename string, object data.Object) error {
	if object == nil || object.ID() == "" {
		return nil
	}
	var docRef *firestore.DocumentRef
	if object.Ref() == nil {
		docRef = c.conn.getDocRef(tablename, object.ID())
	} else {
		docRef = object.Ref().(*firestore.DocumentRef)
	}

	err := c.tx.Delete(docRef)
	if err != nil {
		return errors.Wrap(err, "failed to delete in transaction: "+errorID(tablename, object.ID()))
	}
	object.SetRef(nil)
	object.SetID("")
	return nil
}

// Clear keep delete all object in table until ctx timeout or all object deleted. it delete 500 documents at a time, if in transaction only 10 documents can be delete
//
//	err := c.Clear(ctx, tablename)
//
func (c *TransactionFirestore) Clear(ctx context.Context, tablename string) error {
	collectionRef := c.conn.getCollectionRef(tablename)
	iter := c.tx.Documents(collectionRef.Query.Limit(data.LimitTransactionClear))
	defer iter.Stop()
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return errors.Wrap(err, "failed to iterator documents: "+errorID(tablename, ""))
		}
		c.tx.Delete(doc.Ref)
	}
	return nil
}
