package data

import (
	"context"
	"strconv"

	"cloud.google.com/go/firestore"
	"github.com/piyuo/libsrv/gcp"
	"github.com/piyuo/libsrv/identifier"
	"github.com/pkg/errors"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// ConnectionFirestore implement firestore connection
//
type ConnectionFirestore struct {
	Connection

	// client is firestore native client, don't forget close it after use
	//
	client *firestore.Client

	//tx is curenet transacton, it is nil if not in transaction
	//
	tx *firestore.Transaction

	//batch is curenet batch, it is nil if not in batch
	//
	batch *firestore.WriteBatch
}

// FirestoreGlobalConnection create global firestore connection
//
//	ctx := context.Background()
//	conn, err := FirestoreGlobalConnection(ctx)
//	defer c.Close()
//
func FirestoreGlobalConnection(ctx context.Context) (Connection, error) {
	cred, err := gcp.GlobalCredential(ctx)
	if err != nil {
		return nil, err
	}
	return firestoreNewConnection(ctx, cred)
}

// FirestoreRegionalConnection create regional database instance
//
//	ctx := context.Background()
//	conn, err := FirestoreRegionalConnection(ctx)
//	defer c.Close()
//
func FirestoreRegionalConnection(ctx context.Context) (Connection, error) {
	cred, err := gcp.RegionalCredential(ctx)
	if err != nil {
		return nil, err
	}
	return firestoreNewConnection(ctx, cred)
}

// firestoreNewConnection create connection to firestore
//
//	cred, err := gcp.RegionalCredential(ctx)
//	return firestoreNewConnection(ctx, cred)
//
func firestoreNewConnection(ctx context.Context, cred *google.Credentials) (Connection, error) {
	client, err := firestore.NewClient(ctx, cred.ProjectID, option.WithCredentials(cred))
	if err != nil {
		return nil, err
	}

	conn := &ConnectionFirestore{
		client: client,
	}
	return conn, nil
}

// errorID return identifier for error
//
//	id := c.errorID("tablename", id)
//	So(id, ShouldEqual, "tablename-id")
//
func (c *ConnectionFirestore) errorID(tablename, name string) string {
	id := tablename
	if name != "" {
		id += "-" + name
	}
	return id
}

// snapshotToObject convert document snapshot to object
//
//	c.snapshotToObject(tablename, docRef, snapshot, object)
//
func (c *ConnectionFirestore) snapshotToObject(tablename string, docRef *firestore.DocumentRef, snapshot *firestore.DocumentSnapshot, object Object) error {
	if snapshot == nil {
		return errors.New("snapshot can not be nil: " + c.errorID(tablename, ""))
	}

	if err := snapshot.DataTo(object); err != nil {
		return errors.Wrap(err, "failed to convert document to object: "+c.errorID(tablename, object.GetID()))
	}
	object.SetRef(docRef)
	object.SetID(docRef.ID)
	return nil
}

// Close database connection
//
//	ctx := context.Background()
//	conn, err := FirestoreRegionalConnection(ctx)
//	defer c.Close()
//
func (c *ConnectionFirestore) Close() {
	c.tx = nil
	if c.client != nil {
		c.client.Close()
		c.client = nil
	}
}

// BatchBegin put connection into batch mode. Set/Update/Delete will hold operation until CommitBatch
//
//	err := c.BatchBegin()
//
func (c *ConnectionFirestore) BatchBegin() {
	c.batch = c.client.Batch()
}

// BatchCommit commit batch operation
//
//	err := c.BatchCommit(ctx)
//
func (c *ConnectionFirestore) BatchCommit(ctx context.Context) error {
	batch := c.batch
	c.batch = nil
	_, err := batch.Commit(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to commit batch")
	}
	return nil
}

// InBatch return true if connection is in batch mode
//
//	inBatch := c.InBatch()
//
func (c *ConnectionFirestore) InBatch() bool {
	return c.batch != nil
}

// Transaction start a transaction
//
//	err := c.Transaction(ctx, func(ctx context.Context) error {
//		return nil
//	})
//
func (c *ConnectionFirestore) Transaction(ctx context.Context, callback func(ctx context.Context) error) error {
	var stopTransaction = func() {
		c.tx = nil
	}

	return c.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		c.tx = tx
		defer stopTransaction()
		return callback(ctx)
	})
}

// InTransaction return true if connection is in transaction
//
//	inTx := c.InTransaction()
//
func (c *ConnectionFirestore) InTransaction() bool {
	return c.tx != nil
}

// getCollectionRef return collection reference in table
//
//	collectionRef, err := c.getCollectionRef(tablename)
//
func (c *ConnectionFirestore) getCollectionRef(tablename string) *firestore.CollectionRef {
	return c.client.Collection(tablename)
}

// getDocRef return document reference in table
//
//	docRef, err := c.getDocRef( tablename, id)
//
func (c *ConnectionFirestore) getDocRef(tablename, id string) *firestore.DocumentRef {
	return c.getCollectionRef(tablename).Doc(id)
}

// Get data object from table, return nil if object does not exist
//
//	factory := func() Object {
//		return &Sample{}
//	}
//	object, err := c.Get(ctx, "sample", id, factory)
//
func (c *ConnectionFirestore) Get(ctx context.Context, tablename, id string, factory func() Object) (Object, error) {
	if id == "" {
		return nil, nil
	}
	docRef := c.getDocRef(tablename, id)
	var err error
	var snapshot *firestore.DocumentSnapshot
	if c.tx != nil {
		snapshot, err = c.tx.Get(docRef)
	} else {
		snapshot, err = docRef.Get(ctx)
	}

	if snapshot != nil && !snapshot.Exists() {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to get document: "+c.errorID(tablename, id))
	}

	object := factory()
	if object == nil {
		return nil, errors.New("failed to create object from factory: " + c.errorID(tablename, id))
	}

	err = c.snapshotToObject(tablename, docRef, snapshot, object)
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
func (c *ConnectionFirestore) Set(ctx context.Context, tablename string, object Object) error {
	if object == nil {
		return errors.New("object can not be nil: " + c.errorID(tablename, ""))
	}
	var docRef *firestore.DocumentRef
	if object.GetRef() == nil { // this is new object
		if object.GetID() == "" {
			object.SetID(identifier.UUID())
		}
		docRef = c.getDocRef(tablename, object.GetID())
		object.SetRef(docRef)
	} else {
		docRef = object.GetRef().(*firestore.DocumentRef)
	}

	var err error
	if c.tx != nil {
		err = c.tx.Set(docRef, object)
	} else if c.batch != nil {
		c.batch.Set(docRef, object)
	} else {
		_, err = docRef.Set(ctx, object)
	}
	if err != nil {
		return errors.Wrap(err, "failed to set object: "+c.errorID(tablename, object.GetID()))
	}
	return nil
}

// Exist return true if object with id exist
//
//	return c.Exist(ctx, tablename, id)
//
func (c *ConnectionFirestore) Exist(ctx context.Context, tablename, id string) (bool, error) {
	if id == "" {
		return false, nil
	}
	var err error
	docRef := c.getDocRef(tablename, id)
	var snapshot *firestore.DocumentSnapshot
	if c.tx != nil {
		snapshot, err = c.tx.Get(docRef)
	} else {
		snapshot, err = docRef.Get(ctx)
	}
	if snapshot != nil && !snapshot.Exists() {
		return false, nil
	}
	if err != nil {
		return false, errors.Wrap(err, "failed to get document: "+c.errorID(tablename, id))
	}
	return true, nil
}

// All return max 10 object, if you need more! using query instead
//
//	return c.All(ctx, tablename, factory)
//
func (c *ConnectionFirestore) All(ctx context.Context, tablename string, factory func() Object) ([]Object, error) {
	collectionRef := c.getCollectionRef(tablename)
	list := []Object{}
	var iter *firestore.DocumentIterator
	if c.tx != nil {
		iter = c.tx.Documents(collectionRef.Query.Limit(limitQueryDefault))
	} else {
		iter = collectionRef.Limit(limitQueryDefault).Documents(ctx)
	}
	defer iter.Stop()

	for {
		snapshot, err := iter.Next()
		if err == iterator.Done {
			break
		} else if err != nil {
			return nil, errors.Wrap(err, "failed to iterator documents: "+c.errorID(tablename, ""))
		}
		object := factory()
		if object == nil {
			return nil, errors.New("failed to create object from factory: " + c.errorID(tablename, ""))
		}

		err = snapshot.DataTo(object)
		if err != nil {
			return nil, errors.Wrap(err, "failed to convert document to object: "+c.errorID(tablename, ""))
		}
		c.snapshotToObject(tablename, snapshot.Ref, snapshot, object)
		list = append(list, object)
	}
	return list, nil
}

// Select return object field from data store, return nil if object does not exist
//
//	return c.Select(ctx, "sample", "sample-id", "Name")
//
func (c *ConnectionFirestore) Select(ctx context.Context, tablename, id, field string) (interface{}, error) {
	docRef := c.getDocRef(tablename, id)
	var err error
	var snapshot *firestore.DocumentSnapshot
	if c.tx != nil {
		snapshot, err = c.tx.Get(docRef)
	} else {
		snapshot, err = docRef.Get(ctx)
	}

	if snapshot != nil && !snapshot.Exists() {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to get document: "+c.errorID(tablename, id))
	}
	value, err := snapshot.DataAt(field)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get value from document: "+c.errorID(tablename, id))
	}
	return value, nil
}

// Update partial object field, create new one if object does not exist,  this function is significant slow than Set()
//
//	err = c.Update(ctx, "sample", "sample-id", map[string]interface{}{
//		"Name": "helloworld",
//	})
//
func (c *ConnectionFirestore) Update(ctx context.Context, tablename, id string, fields map[string]interface{}) error {
	docRef := c.getDocRef(tablename, id)
	if c.tx != nil {
		err := c.tx.Set(docRef, fields, firestore.MergeAll)
		if err != nil {
			return errors.Wrap(err, "failed to update field in transaction: "+c.errorID(tablename, id))
		}
		return nil
	}

	if c.batch != nil {
		c.batch.Set(docRef, fields, firestore.MergeAll)
		return nil
	}

	_, err := docRef.Set(ctx, fields, firestore.MergeAll)
	if err != nil {
		return errors.Wrap(err, "failed to update field: "+c.errorID(tablename, id))
	}
	return nil
}

// Increment value on object field, return error if object does not exist
//
//	err := c.Increment(ctx,"sample", "sample-id", "Value", 1)
//
func (c *ConnectionFirestore) Increment(ctx context.Context, tablename, id, field string, value int) error {
	docRef := c.getDocRef(tablename, id)
	if c.tx != nil {
		err := c.tx.Update(docRef, []firestore.Update{
			{Path: field, Value: firestore.Increment(value)},
		})
		if err != nil {
			return errors.Wrap(err, "failed to increment "+field+" with "+strconv.Itoa(value)+" in transaction: "+c.errorID(tablename, id))
		}
		return nil
	}

	if c.batch != nil {
		c.batch.Update(docRef, []firestore.Update{
			{Path: field, Value: firestore.Increment(value)},
		})
		return nil
	}

	_, err := docRef.Update(ctx, []firestore.Update{
		{Path: field, Value: firestore.Increment(value)},
	})
	if err != nil {
		return errors.Wrap(err, "failed to increment "+field+" with "+strconv.Itoa(value)+": "+c.errorID(tablename, id))
	}
	return nil
}

// Delete object using table name and id, no error if id did not exist
//
//	c.Delete(ctx, "sample", "sample-id")
//
func (c *ConnectionFirestore) Delete(ctx context.Context, tablename, id string) error {
	docRef := c.getDocRef(tablename, id)
	if c.tx != nil {
		err := c.tx.Delete(docRef)
		if err != nil {
			return errors.Wrap(err, "failed to delete in transaction: "+c.errorID(tablename, id))
		}
		return nil
	}

	if c.batch != nil {
		c.batch.Delete(docRef)
		return nil
	}

	_, err := docRef.Delete(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to delete: "+c.errorID(tablename, id))
	}
	return nil
}

// DeleteBatch delete list of id use batch mode, no error if id not exist
//
//	c.DeleteBatch(ctx, dt.tablename, ids)
//
func (c *ConnectionFirestore) DeleteBatch(ctx context.Context, tablename string, ids []string) error {
	if len(ids) == 0 {
		return nil
	}

	c.BatchBegin()
	for _, id := range ids {
		c.Delete(ctx, tablename, id)
	}
	if err := c.BatchCommit(ctx); err != nil {
		return err
	}
	return nil
}

// DeleteObject delete object, no error if id did not exist
//
//	c.DeleteObject(ctx, "sample", object)
//
func (c *ConnectionFirestore) DeleteObject(ctx context.Context, tablename string, object Object) error {
	if object == nil || object.GetID() == "" {
		return nil
	}
	var docRef *firestore.DocumentRef
	if object.GetRef() == nil {
		docRef = c.getDocRef(tablename, object.GetID())
	} else {
		docRef = object.GetRef().(*firestore.DocumentRef)
	}

	if c.tx != nil {
		err := c.tx.Delete(docRef)
		if err != nil {
			return errors.Wrap(err, "failed to delete in transaction: "+c.errorID(tablename, object.GetID()))
		}
	} else if c.batch != nil {
		c.batch.Delete(docRef)

	} else {
		_, err := docRef.Delete(ctx)
		if err != nil {
			return errors.Wrap(err, "failed to delete: "+c.errorID(tablename, object.GetID()))
		}
	}
	object.SetRef(nil)
	object.SetID("")
	return nil
}

// Clear keep delete all object in table until ctx timeout or all object deleted. it delete 500 documents at a time, if in transaction only 10 documents can be delete
//
//	err := c.Clear(ctx, tablename)
//
func (c *ConnectionFirestore) Clear(ctx context.Context, tablename string) error {
	collectionRef := c.getCollectionRef(tablename)
	if c.tx != nil {
		iter := c.tx.Documents(collectionRef.Query.Limit(limitTransactionClear))
		defer iter.Stop()
		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return errors.Wrap(err, "failed to iterator documents: "+c.errorID(tablename, ""))
			}
			c.tx.Delete(doc.Ref)
		}
		return nil
	}
	for {
		// keep delete until ctx timeout or all object deleted
		if ctx.Err() != nil {
			return ctx.Err()
		}
		numDeleted := 0
		iter := collectionRef.Limit(limitClear).Documents(ctx)
		defer iter.Stop()
		// Iterate through the documents, adding a delete operation for each one to a WriteBatch.
		batch := c.client.Batch()
		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return errors.Wrap(err, "failed to iterator documents: "+c.errorID(tablename, ""))
			}
			batch.Delete(doc.Ref)
			numDeleted++
		}
		if numDeleted > 0 {
			_, err := batch.Commit(ctx)
			if err != nil {
				return errors.Wrap(err, "failed to commit batch: "+c.errorID(tablename, ""))
			}
		}
		if numDeleted < limitClear {
			break
		}
	}
	return nil
}

// Query create query
//
//	factory := func() Object {
//		return &Sample{}
//	}
//	c.Query(ctx, "sample", factory).Execute(ctx)
//
func (c *ConnectionFirestore) Query(tablename string, factory func() Object) Query {
	return &QueryFirestore{
		BaseQuery: BaseQuery{factory: factory},
		query:     c.getCollectionRef(tablename).Query,
		conn:      c,
		tablename: tablename,
	}
}
