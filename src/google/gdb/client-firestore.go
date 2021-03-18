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

// ClientFirestore implement firestore connection
//
type ClientFirestore struct {
	db.BaseClient

	// client is firestore native client
	//
	firestoreClient *firestore.Client
}

// Close client
//
//	client, err := NewClient(ctx,cred)
//	defer c.Close()
//
func (c *ClientFirestore) Close() {
	if c.firestoreClient != nil {
		c.firestoreClient.Close()
		c.firestoreClient = nil
	}
}

// Batch start a batch operation
//
//	err := Batch(ctx, func(ctx context.Context,bc db.Batch) error {
//		return nil
//	})
//
func (c *ClientFirestore) Batch(ctx context.Context, f db.BatchFunc) error {
	batch := c.firestoreClient.Batch()
	bc := &BatchFirestore{
		client: c,
		batch:  batch,
	}
	err := f(ctx, bc)
	if err != nil {
		return errors.Wrapf(err, "run batch func")
	}

	_, err = batch.Commit(ctx)
	if err != nil {
		return errors.Wrapf(err, "commit batch")
	}
	return err

}

// Transaction start a transaction
//
//	err := Transaction(ctx, func(ctx context.Context,tx db.Transaction) error {
//		return nil
//	})
//
func (c *ClientFirestore) Transaction(ctx context.Context, f db.TransactionFunc) error {
	return c.firestoreClient.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		trans := &TransactionFirestore{
			client: c,
			tx:     tx,
		}
		return f(ctx, trans)
	})
}

// getCollectionRef return collection reference in table
//
//	collectionRef, err := c.getCollectionRef(tablename)
//
func (c *ClientFirestore) getCollectionRef(tablename string) *firestore.CollectionRef {
	return c.firestoreClient.Collection(tablename)
}

// getDocRef return document reference in table
//
//	docRef, err := c.getDocRef( tablename, id)
//
func (c *ClientFirestore) getDocRef(tablename, id string) *firestore.DocumentRef {
	return c.getCollectionRef(tablename).Doc(id)
}

// snapshotToObject convert snap shot to object
//
func snapshotToObject(obj db.Object, docRef *firestore.DocumentRef, snapshot *firestore.DocumentSnapshot, err error) (db.Object, error) {
	if snapshot != nil && !snapshot.Exists() {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Wrapf(err, "get snapshot from %v-%v", obj.Collection(), docRef.ID)
	}

	if err := snapshot.DataTo(obj); err != nil {
		return nil, errors.Wrapf(err, "make snapshot to object %v-%v", obj.Collection(), docRef.ID)
	}
	obj.SetRef(docRef)
	obj.SetID(docRef.ID)
	return obj, nil
}

// snapshotExists return true if snapshot exists
//
func snapshotExists(obj db.Object, id string, snapshot *firestore.DocumentSnapshot, err error) (bool, error) {
	if snapshot != nil && !snapshot.Exists() {
		return false, nil
	}
	if err != nil {
		return false, errors.Wrapf(err, "get snapshot %v-%v", obj.Collection(), obj.ID())
	}
	return true, nil
}

// select return object field from data store, return nil if object does not exist
//
func snapshotToField(obj db.Object, id, field string, snapshot *firestore.DocumentSnapshot, err error) (interface{}, error) {
	if snapshot != nil && !snapshot.Exists() {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Wrapf(err, "get snapshot %v-%v ", obj.Collection(), id)
	}
	value, err := snapshot.DataAt(field)
	if err != nil {
		return nil, errors.Wrapf(err, "get data at field %v %v-%v ", field, obj.Collection(), id)
	}
	return value, nil
}

// iterObjects convert list of snapshot to list of object
//
func iterObjects(obj db.Object, iter *firestore.DocumentIterator) ([]db.Object, error) {
	list := []db.Object{}
	for {
		snapshot, err := iter.Next()
		if err == iterator.Done {
			break
		} else if err != nil {
			return nil, errors.Wrapf(err, "iter next %v", obj.Collection())
		}
		newObj := obj.Factory()
		if newObj == nil {
			return nil, errors.New(obj.Collection() + " not implement Factory()")
		}

		_, err = snapshotToObject(newObj, snapshot.Ref, snapshot, err)
		if err != nil {
			return nil, errors.Wrapf(err, "iter snapshot to object %v-%v", obj.Collection(), snapshot.Ref.ID)
		}
		list = append(list, newObj)
	}
	return list, nil
}

// objDeleteRef return reference to delete object
//
func (c *ClientFirestore) objDeleteRef(obj db.Object) *firestore.DocumentRef {
	var docRef *firestore.DocumentRef
	if obj.Ref() == nil {
		docRef = c.getDocRef(obj.Collection(), obj.ID())
	} else {
		docRef = obj.Ref().(*firestore.DocumentRef)
	}
	obj.SetRef(nil)
	obj.SetID("")
	return docRef
}

// refFromObj get docRef for object
//
func (c *ClientFirestore) refFromObj(ctx context.Context, obj db.Object) *firestore.DocumentRef {
	var docRef *firestore.DocumentRef
	if obj.Ref() == nil { // new object
		docRef = c.getDocRef(obj.Collection(), obj.ID())
		obj.SetRef(docRef)
	} else { // object already exist
		docRef = obj.Ref().(*firestore.DocumentRef)
	}
	return docRef
}

// Get data object from table, return nil if object does not exist
//
//	object, err := Get(ctx, &Sample{}, "id")
//
func (c *ClientFirestore) Get(ctx context.Context, obj db.Object, id string) (db.Object, error) {
	if err := db.Check(ctx, obj, false); err != nil {
		return nil, err
	}
	if err := db.CheckID(id); err != nil {
		return nil, err
	}
	docRef := c.getDocRef(obj.Collection(), id)
	snapshot, err := docRef.Get(ctx)
	return snapshotToObject(obj, docRef, snapshot, err)
}

// Exists return true if object with id exist
//
//	found,err := Exists(ctx, &Sample{}, "id")
//
func (c *ClientFirestore) Exists(ctx context.Context, obj db.Object, id string) (bool, error) {
	if err := db.Check(ctx, obj, false); err != nil {
		return false, err
	}
	if err := db.CheckID(id); err != nil {
		return false, err
	}
	docRef := c.getDocRef(obj.Collection(), id)
	snapshot, err := docRef.Get(ctx)
	return snapshotExists(obj, id, snapshot, err)
}

// List return object list, use max to specific return object count
//
//	list,err := List(ctx, &Sample{},10)
//
func (c *ClientFirestore) List(ctx context.Context, obj db.Object, max int) ([]db.Object, error) {
	if err := db.Check(ctx, obj, false); err != nil {
		return nil, err
	}
	collectionRef := c.getCollectionRef(obj.Collection())
	iter := collectionRef.Limit(max).Documents(ctx)
	defer iter.Stop()
	return iterObjects(obj, iter)
}

// Select return object field from data store, return nil if object does not exist
//
//	return Select(ctx, &Sample{}, id, field)
//
func (c *ClientFirestore) Select(ctx context.Context, obj db.Object, id, field string) (interface{}, error) {
	if err := db.Check(ctx, obj, false); err != nil {
		return false, err
	}
	if err := db.CheckID(id); err != nil {
		return false, err
	}
	docRef := c.getDocRef(obj.Collection(), id)
	snapshot, err := docRef.Get(ctx)
	return snapshotToField(obj, id, field, snapshot, err)
}

// Query create query
//
//	c.Query(ctx, &Sample{}).Return(ctx)
//
func (c *ClientFirestore) Query(obj db.Object) db.Query {
	return (&QueryFirestore{
		BaseQuery: db.BaseQuery{
			QueryObject: obj,
		},
		query:  c.getCollectionRef(obj.Collection()).Query,
		client: c,
	}).Limit(20)
}

// Set object into table, If the document not exist, it will be created. If the document does exist, its contents will be overwritten with the newly provided data, if object does not have id, it will created using UUID
//
//	 err := Set(ctx, object)
//
func (c *ClientFirestore) Set(ctx context.Context, obj db.Object) error {
	if err := db.Check(ctx, obj, false); err != nil {
		return err
	}
	c.BaseClient.BeforeSet(ctx, obj)
	docRef := c.refFromObj(ctx, obj)
	_, err := docRef.Set(ctx, obj)
	if err != nil {
		return errors.Wrapf(err, "set doc %v-%v", obj.Collection(), obj.ID())
	}
	return nil
}

// Update partial object field, create new one if object does not exist, this function is significant slow than Set()
//
//	err = Update(ctx, Sample, map[string]interface{}{
//		"desc": "hi",
//	})
//
func (c *ClientFirestore) Update(ctx context.Context, obj db.Object, fields map[string]interface{}) error {
	if err := db.Check(ctx, obj, true); err != nil {
		return err
	}
	docRef := c.getDocRef(obj.Collection(), obj.ID())
	_, err := docRef.Set(ctx, fields, firestore.MergeAll)
	if err != nil {
		fieldStr := mapping.ToString(fields)
		return errors.Wrapf(err, "update field %v %v-%v"+fieldStr, obj.Collection(), obj.ID())
	}
	return nil
}

// Increment value on object field, return error if object does not exist
//
//	err := Increment(ctx,sample, "Value", 2)
//
func (c *ClientFirestore) Increment(ctx context.Context, obj db.Object, field string, value int) error {
	if err := db.Check(ctx, obj, true); err != nil {
		return err
	}
	docRef := c.getDocRef(obj.Collection(), obj.ID())
	_, err := docRef.Update(ctx, []firestore.Update{
		{Path: field, Value: firestore.Increment(value)},
	})
	if err != nil {
		return errors.Wrapf(err, "inc field "+strconv.Itoa(value)+" %v-%v", obj.Collection(), obj.ID())
	}
	return nil
}

// Delete object, no error if id not exist
//
//	Delete(ctx, sample)
//
func (c *ClientFirestore) Delete(ctx context.Context, obj db.Object) error {
	if err := db.Check(ctx, obj, true); err != nil {
		return err
	}
	docRef := c.objDeleteRef(obj)
	_, err := docRef.Delete(ctx)
	if err != nil {
		return errors.Wrapf(err, "delete %v-%v", obj.Collection(), obj.ID())
	}
	return nil
}

// DeleteCollection delete document in collection. delete max doc count. return true if no doc left in collection
//
//	cleared, err := DeleteCollection(ctx, 50, iter)
//
func (c *ClientFirestore) DeleteCollection(ctx context.Context, max int, iter *firestore.DocumentIterator) (bool, error) {
	numDeleted := 0
	// Iterate through the documents, adding a delete operation for each one to a WriteBatch.
	err := c.Batch(ctx, func(ctx context.Context, batch db.Batch) error {
		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return errors.Wrap(err, "iter next")
			}
			batch.DeleteRef(doc.Ref)
			numDeleted++
		}
		return nil
	})
	if err != nil {
		return false, errors.Wrapf(err, "commit batch")
	}

	if numDeleted < max {
		return true, nil
	}
	return false, nil
}

// Clear delete all document in collection. delete max doc count. return true if collection is cleared
//
//	cleared, err := Clear(ctx, &Sample{}, 50)
//
func (c *ClientFirestore) Clear(ctx context.Context, obj db.Object, max int) (bool, error) {
	if err := db.Check(ctx, obj, false); err != nil {
		return false, err
	}
	collectionRef := c.getCollectionRef(obj.Collection())
	iter := collectionRef.Limit(max).Documents(ctx)
	defer iter.Stop()
	cleared, err := c.DeleteCollection(ctx, max, iter)
	if err != nil {
		return false, errors.Wrap(err, "delete "+obj.Collection())
	}
	return cleared, nil
}

// Counter return counter, set numshards 100 times of concurrent usage. for example if you think concurrent use is 10/seconds then set numshards to 1000 to avoid too much retention error
// if DateHierarchyFull is set, counter will automatically generate year/month/day/hour hierarchy in utc timezone
//
//	sampleCounter,err = Counter("SampleCount", 100, DateHierarchyNone)
//
func (c *ClientFirestore) Counter(counterName string, numshards int, hierarchy db.DateHierarchy) db.Counter {
	if numshards <= 0 {
		numshards = 10
	}

	return &CounterFirestore{
		MetaFirestore: MetaFirestore{
			client:     c,
			collection: "Count",
			id:         counterName,
			numShards:  numshards,
		},
		keepDateHierarchy: hierarchy == db.DateHierarchyFull,
	}
}

// Coder return coder, set numshards 100 times of concurrent usage. for example if you think concurrent use is 10/seconds then set numshards to 1000 to avoid too much retention error
//
//	productCoder,err = Coder("coderName",100)
//
func (c *ClientFirestore) Coder(coderName string, numshards int) db.Coder {
	return &CoderFirestore{
		MetaFirestore: MetaFirestore{
			client:     c,
			collection: "Code",
			id:         coderName,
			numShards:  numshards,
		},
	}
}

// Serial return serial, create one if not exist, please be aware Serial can only generate 1 number per second, use serial with high frequency will cause too much retention error
//
//	productNo,err = Serial("serialName")
//
func (c *ClientFirestore) Serial(serialName string) db.Serial {
	return &SerialFirestore{
		MetaFirestore: MetaFirestore{
			client:     c,
			collection: "Serial",
			id:         serialName,
			numShards:  0,
		},
	}
}
