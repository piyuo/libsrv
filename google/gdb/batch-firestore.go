package gdb

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/piyuo/libsrv/db"
)

// BatchFirestore implement firestore batch
//
type BatchFirestore struct {
	db.Batch

	// client gdb client
	//
	client *ClientFirestore

	// batch is curenet batch, it is nil if not in batch
	//
	batch *firestore.WriteBatch

	// hasSomethingToCommit set to true when batch operation has been called like set/update/delete
	//
	hasSomethingToCommit bool
}

// Set object into table, If the document not exist, it will be created. If the document does exist, its contents will be overwritten with the newly provided data, if object does not have id, it will created using UUID
//
//	 Set(ctx, object)
//
func (c *BatchFirestore) Set(ctx context.Context, obj db.Object) {
	c.client.BaseClient.BeforeSet(ctx, obj)
	docRef := c.client.refFromObj(ctx, obj)
	c.batch.Set(docRef, obj)
	c.hasSomethingToCommit = true
}

// Update partial object field, create new one if object does not exist, this function is significant slow than Set()
//
//	Update(ctx, Sample, map[string]interface{}{
//		"desc": "hi",
//	})
//
func (c *BatchFirestore) Update(ctx context.Context, obj db.Object, fields map[string]interface{}) {
	docRef := c.client.getDocRef(obj.Collection(), obj.ID())
	c.batch.Set(docRef, fields, firestore.MergeAll)
	c.hasSomethingToCommit = true
}

// Increment value on object field, return error if object does not exist
//
//	Increment(ctx,sample, "Value", 2)
//
func (c *BatchFirestore) Increment(ctx context.Context, obj db.Object, field string, value int) {
	docRef := c.client.getDocRef(obj.Collection(), obj.ID())
	c.batch.Update(docRef, []firestore.Update{
		{Path: field, Value: firestore.Increment(value)},
	})
	c.hasSomethingToCommit = true
}

// Delete object, no error if id not exist
//
//	Delete(ctx, sample)
//
func (c *BatchFirestore) Delete(ctx context.Context, obj db.Object) {
	docRef := c.client.objDeleteRef(obj)
	c.batch.Delete(docRef)
	c.hasSomethingToCommit = true
}

// DeleteList delete object use list of id, no error if id not exist
//
//	DeleteList(ctx, &Sample{}, []string{"1","2"})
//
func (c *BatchFirestore) DeleteList(ctx context.Context, obj db.Object, list []string) {
	obj.SetRef(nil)
	for _, id := range list {
		obj.SetID(id)
		c.Delete(ctx, obj)
	}
	if len(list) > 0 {
		c.hasSomethingToCommit = true
	}
	obj.SetID("")
}

// DeleteRef delete object use document ref
//
//	DeleteRef(ref)
//
func (c *BatchFirestore) DeleteRef(ref *firestore.DocumentRef) {
	c.batch.Delete(ref)
	c.hasSomethingToCommit = true
}
