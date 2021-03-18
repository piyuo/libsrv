package gdb

import (
	"context"
	"strconv"

	"cloud.google.com/go/firestore"
	"github.com/piyuo/libsrv/src/db"
	"github.com/piyuo/libsrv/src/mapping"
	"github.com/pkg/errors"
	"google.golang.org/api/iterator"
)

// TransactionFirestore implement firestore transaction
//
type TransactionFirestore struct {
	db.Transaction

	// client is gdb connection
	//
	client *ClientFirestore

	//tx is curenet transacton, it is nil if not in transaction
	//
	tx *firestore.Transaction
}

// Get data object from table, return nil if object does not exist
//
//	object, err := Get(ctx, &Sample{}, "id")
//
func (c *TransactionFirestore) Get(ctx context.Context, obj db.Object, id string) (db.Object, error) {
	if err := db.Check(ctx, obj, false); err != nil {
		return nil, err
	}
	if err := db.CheckID(id); err != nil {
		return nil, err
	}
	docRef := c.client.getDocRef(obj.Collection(), id)
	snapshot, err := c.tx.Get(docRef)
	return snapshotToObject(obj, docRef, snapshot, err)
}

// Exists return true if object with id exist
//
//	found,err := Exists(ctx, &Sample{}, "id")
//
func (c *TransactionFirestore) Exists(ctx context.Context, obj db.Object, id string) (bool, error) {
	if err := db.Check(ctx, obj, false); err != nil {
		return false, err
	}
	if err := db.CheckID(id); err != nil {
		return false, err
	}
	docRef := c.client.getDocRef(obj.Collection(), id)
	snapshot, err := c.tx.Get(docRef)
	return snapshotExists(obj, id, snapshot, err)
}

// List return object list, use max to specific return object count
//
//	list,err := List(ctx, &Sample{},10)
//
func (c *TransactionFirestore) List(ctx context.Context, obj db.Object, max int) ([]db.Object, error) {
	if err := db.Check(ctx, obj, false); err != nil {
		return nil, err
	}
	collectionRef := c.client.getCollectionRef(obj.Collection())
	iter := c.tx.Documents(collectionRef.Query.Limit(max))
	defer iter.Stop()
	return iterObjects(obj, iter)
}

// Select return object field from data store, return nil if object does not exist
//
//	return Select(ctx, &Sample{}, id, field)
//
func (c *TransactionFirestore) Select(ctx context.Context, obj db.Object, id, field string) (interface{}, error) {
	if err := db.Check(ctx, obj, false); err != nil {
		return false, err
	}
	if err := db.CheckID(id); err != nil {
		return false, err
	}
	docRef := c.client.getDocRef(obj.Collection(), id)
	snapshot, err := c.tx.Get(docRef)
	return snapshotToField(obj, id, field, snapshot, err)
}

// Query create query
//
//	c.Query(ctx, &Sample{}).Return(ctx)
//
func (c *TransactionFirestore) Query(obj db.Object) db.Query {
	return (&QueryFirestore{
		BaseQuery: db.BaseQuery{
			QueryTransaction: c,
			QueryObject:      obj,
		},
		query:  c.client.getCollectionRef(obj.Collection()).Query,
		client: c.client,
	}).Limit(20)
}

// Set object into table, If the document not exist, it will be created. If the document does exist, its contents will be overwritten with the newly provided data, if object does not have id, it will created using UUID
//
//	 err := Set(ctx, object)
//
func (c *TransactionFirestore) Set(ctx context.Context, obj db.Object) error {
	if err := db.Check(ctx, obj, false); err != nil {
		return err
	}
	c.client.BaseClient.BeforeSet(ctx, obj)
	docRef := c.client.refFromObj(ctx, obj)
	err := c.tx.Set(docRef, obj)
	if err != nil {
		return errors.Wrapf(err, "tx set doc %v-%v", obj.Collection(), obj.ID())
	}
	return nil
}

// Update partial object field, create new one if object does not exist, this function is significant slow than Set()
//
//	err = Update(ctx, Sample, map[string]interface{}{
//		"desc": "hi",
//	})
//
func (c *TransactionFirestore) Update(ctx context.Context, obj db.Object, fields map[string]interface{}) error {
	if err := db.Check(ctx, obj, true); err != nil {
		return err
	}
	docRef := c.client.getDocRef(obj.Collection(), obj.ID())
	err := c.tx.Set(docRef, fields, firestore.MergeAll)
	if err != nil {
		fieldStr := mapping.ToString(fields)
		return errors.Wrapf(err, "tx update field %v %v-%v"+fieldStr, obj.Collection(), obj.ID())
	}
	return nil
}

// Increment value on object field, return error if object does not exist
//
//	err := Increment(ctx,sample, "Value", 2)
//
func (c *TransactionFirestore) Increment(ctx context.Context, obj db.Object, field string, value int) error {
	if err := db.Check(ctx, obj, true); err != nil {
		return err
	}
	docRef := c.client.getDocRef(obj.Collection(), obj.ID())
	err := c.tx.Update(docRef, []firestore.Update{
		{Path: field, Value: firestore.Increment(value)},
	})
	if err != nil {
		return errors.Wrapf(err, "tx inc field "+strconv.Itoa(value)+" %v-%v", obj.Collection(), obj.ID())
	}
	return nil
}

// Delete object, no error if id not exist
//
//	Delete(ctx, sample)
//
func (c *TransactionFirestore) Delete(ctx context.Context, obj db.Object) error {
	if err := db.Check(ctx, obj, true); err != nil {
		return err
	}
	docRef := c.client.objDeleteRef(obj)
	err := c.tx.Delete(docRef)
	if err != nil {
		return errors.Wrapf(err, "tx delete %v-%v", obj.Collection(), obj.ID())
	}
	return nil
}

// DeleteCollection delete document in collection. delete max doc count. return true if no doc left in collection
//
//	cleared, err := DeleteCollection(ctx, &Sample{}, 50, iter)
//
func (c *TransactionFirestore) DeleteCollection(ctx context.Context, obj db.Object, max int, iter *firestore.DocumentIterator) (bool, error) {
	numDeleted := 0
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return false, errors.Wrap(err, "tx iter next "+obj.Collection())
		}
		if err := c.tx.Delete(doc.Ref); err != nil {
			return false, errors.Wrap(err, "tx delete "+obj.Collection()+"-"+doc.Ref.ID)
		}
		numDeleted++
	}
	if numDeleted < max {
		return true, nil
	}
	return false, nil
}

// Clear delete all document in collection. try 10 batch each batch only delete document specific in max. return true if collection is cleared
//
//	cleared, err := Clear(ctx, &Sample{}, 50)
//
func (c *TransactionFirestore) Clear(ctx context.Context, obj db.Object, max int) (bool, error) {
	if err := db.Check(ctx, obj, false); err != nil {
		return false, err
	}
	collectionRef := c.client.getCollectionRef(obj.Collection())
	iter := c.tx.Documents(collectionRef.Query.Limit(max))
	defer iter.Stop()
	return c.DeleteCollection(ctx, obj, max, iter)
}

// isShardExists return true if shard already exist
//
func (c *TransactionFirestore) snapshot(ctx context.Context, ref *firestore.DocumentRef) (*firestore.DocumentSnapshot, error) {
	snapshot, err := c.tx.Get(ref)
	if snapshot != nil && !snapshot.Exists() {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Wrap(err, "get snapshot")
	}
	return snapshot, nil
}

// isShardExists return true if shard already exist
//
func (c *TransactionFirestore) isShardExists(ctx context.Context, ref *firestore.DocumentRef) (bool, error) {
	snapshot, err := c.tx.Get(ref)
	if snapshot != nil && !snapshot.Exists() {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// createShard create a shard
//
func (c *TransactionFirestore) createShard(ref *firestore.DocumentRef, shard map[string]interface{}) error {
	if shard[db.MetaValue] == nil {
		return errors.New("N must not nil")
	}

	err := c.tx.Set(ref, shard, firestore.MergeAll)
	if err != nil {
		return errors.Wrap(err, "tx set shard")
	}
	return nil
}

// incrementShard increment shard count
//
func (c *TransactionFirestore) incrementShard(ref *firestore.DocumentRef, value interface{}) error {
	if value == nil {
		return errors.New("value must not nil")
	}

	err := c.tx.Update(ref, []firestore.Update{
		{Path: db.MetaValue, Value: firestore.Increment(value)},
	})
	if err != nil {
		return errors.Wrap(err, "tx update shard")
	}
	return nil
}
