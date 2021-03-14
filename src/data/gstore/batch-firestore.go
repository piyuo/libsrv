package gstore

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/piyuo/libsrv/src/data"
	"github.com/piyuo/libsrv/src/identifier"
	"github.com/pkg/errors"
)

// BatchFirestore implement firestore batch
//
type BatchFirestore struct {
	data.Batch

	// conn is firestore connection
	//
	conn *ConnectionFirestore

	//batch is curenet batch, it is nil if not in batch
	//
	batch *firestore.WriteBatch
}

// Begin put connection into batch mode. Set/Update/Delete will hold operation until CommitBatch
//
//	err := c.Begin()
//
func (c *BatchFirestore) Begin() {
	c.batch = c.conn.client.Batch()
}

// Commit batch operation
//
//	err := c.Commit(ctx)
//
func (c *BatchFirestore) Commit(ctx context.Context) error {
	batch := c.batch
	c.batch = nil
	_, err := batch.Commit(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to commit batch")
	}
	return nil
}

// Set object into table, If the document does not exist, it will be created. If the document does exist, its contents will be overwritten with the newly provided data, if object does not have id, it will created using UUID
//
//	if err := c.Set(ctx, tablename, object); err != nil {
//		return err
//	}
//
func (c *BatchFirestore) Set(ctx context.Context, tablename string, object data.Object) error {
	if object == nil {
		return errors.New("object can not be nil: " + errorID(tablename, ""))
	}
	var docRef *firestore.DocumentRef
	if object.Ref() == nil { // this is new object
		if object.ID() == "" {
			object.SetID(identifier.UUID())
		}
		docRef = c.conn.getDocRef(tablename, object.ID())
		object.SetRef(docRef)
	} else {
		docRef = object.Ref().(*firestore.DocumentRef)
	}
	c.batch.Set(docRef, object)
	return nil
}

// Update partial object field, create new one if object does not exist,  this function is significant slow than Set()
//
//	err = c.Update(ctx, "sample", "sample-id", map[string]interface{}{
//		"Name": "helloworld",
//	})
//
func (c *BatchFirestore) Update(ctx context.Context, tablename, id string, fields map[string]interface{}) error {
	docRef := c.conn.getDocRef(tablename, id)
	c.batch.Set(docRef, fields, firestore.MergeAll)
	return nil
}

// Increment value on object field, return error if object does not exist
//
//	err := c.Increment(ctx,"sample", "sample-id", "Value", 1)
//
func (c *BatchFirestore) Increment(ctx context.Context, tablename, id, field string, value int) error {
	docRef := c.conn.getDocRef(tablename, id)
	c.batch.Update(docRef, []firestore.Update{
		{Path: field, Value: firestore.Increment(value)},
	})
	return nil
}

// Delete object using table name and id, no error if id did not exist
//
//	c.Delete(ctx, "sample", "sample-id")
//
func (c *BatchFirestore) Delete(ctx context.Context, tablename, id string) error {
	docRef := c.conn.getDocRef(tablename, id)
	c.batch.Delete(docRef)
	return nil
}

// DeleteBatch delete list of id use batch mode, no error if id not exist
//
//	c.DeleteBatch(ctx, dt.tablename, ids)
//
func (c *BatchFirestore) DeleteBatch(ctx context.Context, tablename string, ids []string) error {
	if len(ids) == 0 {
		return nil
	}

	c.Begin()
	for _, id := range ids {
		c.Delete(ctx, tablename, id)
	}
	if err := c.Commit(ctx); err != nil {
		return err
	}
	return nil
}

// DeleteObject delete object, no error if id did not exist
//
//	c.DeleteObject(ctx, "sample", object)
//
func (c *BatchFirestore) DeleteObject(ctx context.Context, tablename string, object data.Object) error {
	if object == nil || object.ID() == "" {
		return nil
	}
	var docRef *firestore.DocumentRef
	if object.Ref() == nil {
		docRef = c.conn.getDocRef(tablename, object.ID())
	} else {
		docRef = object.Ref().(*firestore.DocumentRef)
	}

	c.batch.Delete(docRef)
	object.SetRef(nil)
	object.SetID("")
	return nil
}
