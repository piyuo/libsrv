package gstore

import (
	"context"
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
	data.Connection

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
//	factory := func() Object {
//		return &Sample{}
//	}
//	object, err := c.Get(ctx, "sample", id, factory)
//
func (c *TransactionFirestore) Get(ctx context.Context, tablename, id string, factory func() data.Object) (data.Object, error) {
	if id == "" {
		return nil, nil
	}
	docRef := c.conn.getDocRef(tablename, id)
	snapshot, err := c.tx.Get(docRef)

	if snapshot != nil && !snapshot.Exists() {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to get document: "+errorID(tablename, id))
	}

	object := factory()
	if object == nil {
		return nil, errors.New("failed to create object from factory: " + errorID(tablename, id))
	}

	err = c.conn.snapshotToObject(tablename, docRef, snapshot, object)
	if err != nil {
		return nil, err
	}
	return object, nil
}

// Set object into table, If the document does not exist, it will be created. If the document does exist, its contents will be overwritten with the newly provided data, if object does not have id, it will created using UUID
//
//	if err := c.Set(ctx, tablename, object); err != nil {
//		return err
//	}
//
func (c *TransactionFirestore) Set(ctx context.Context, tablename string, object data.Object) error {
	if object == nil {
		return errors.New("object can not be nil: " + errorID(tablename, ""))
	}
	var docRef *firestore.DocumentRef
	if object.GetRef() == nil { // this is new object
		if object.GetID() == "" {
			object.SetID(identifier.UUID())
		}
		docRef = c.conn.getDocRef(tablename, object.GetID())
		object.SetRef(docRef)
	} else {
		docRef = object.GetRef().(*firestore.DocumentRef)
	}

	err := c.tx.Set(docRef, object)

	if err != nil {
		return errors.Wrap(err, "failed to set object: "+errorID(tablename, object.GetID()))
	}
	return nil
}

// IsExists return true if object with id exist
//
//	return c.IsExists(ctx, tablename, id)
//
func (c *TransactionFirestore) IsExists(ctx context.Context, tablename, id string) (bool, error) {
	if id == "" {
		return false, nil
	}
	docRef := c.conn.getDocRef(tablename, id)
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
//	return c.All(ctx, tablename, factory)
//
func (c *TransactionFirestore) All(ctx context.Context, tablename string, factory func() data.Object) ([]data.Object, error) {
	collectionRef := c.conn.getCollectionRef(tablename)
	list := []data.Object{}
	iter := c.tx.Documents(collectionRef.Query.Limit(data.LimitQueryDefault))
	defer iter.Stop()

	for {
		snapshot, err := iter.Next()
		if err == iterator.Done {
			break
		} else if err != nil {
			return nil, errors.Wrap(err, "failed to iterator documents: "+errorID(tablename, ""))
		}
		object := factory()
		if object == nil {
			return nil, errors.New("failed to create object from factory: " + errorID(tablename, ""))
		}

		err = snapshot.DataTo(object)
		if err != nil {
			return nil, errors.Wrap(err, "failed to convert document to object: "+errorID(tablename, ""))
		}
		c.conn.snapshotToObject(tablename, snapshot.Ref, snapshot, object)
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
	if object == nil || object.GetID() == "" {
		return nil
	}
	var docRef *firestore.DocumentRef
	if object.GetRef() == nil {
		docRef = c.conn.getDocRef(tablename, object.GetID())
	} else {
		docRef = object.GetRef().(*firestore.DocumentRef)
	}

	err := c.tx.Delete(docRef)
	if err != nil {
		return errors.Wrap(err, "failed to delete in transaction: "+errorID(tablename, object.GetID()))
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
