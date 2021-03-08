package gstore

import (
	"context"
	"strconv"

	"cloud.google.com/go/firestore"
	"github.com/piyuo/libsrv/src/data"
	"github.com/piyuo/libsrv/src/gaccount"
	"github.com/piyuo/libsrv/src/identifier"
	"github.com/pkg/errors"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// firestoreGlobalConnection keep global connection to reuse in the future
//
var firestoreGlobalConnection data.Connection

// firestoreRegionalConnection keep regional connection to reuse in the future
//
var firestoreRegionalConnection data.Connection

// ConnectionFirestore implement firestore connection
//
type ConnectionFirestore struct {
	data.Connection

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

// ConnectGlobalFirestore create global firestore connection
//
//	conn, err := ConnectGlobalFirestore(ctx)
//	defer c.Close()
//
func ConnectGlobalFirestore(ctx context.Context) (data.Connection, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	if firestoreGlobalConnection == nil {
		cred, err := gaccount.GlobalCredential(ctx)
		if err != nil {
			return nil, err
		}
		firestoreGlobalConnection, err = firestoreCreateConnection(cred)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create global database connection")
		}
	}
	return firestoreGlobalConnection, nil
}

// ConnectRegionalFirestore create regional database instance
//
//	conn, err := ConnectRegionalFirestore(ctx)
//	defer conn.Close()
//
func ConnectRegionalFirestore(ctx context.Context) (data.Connection, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	if firestoreRegionalConnection == nil {
		cred, err := gaccount.RegionalCredential(ctx)
		if err != nil {
			return nil, err
		}
		firestoreRegionalConnection, err = firestoreCreateConnection(cred)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create regional database connection")
		}
	}
	return firestoreRegionalConnection, nil
}

// firestoreCreateConnection create connection to firestore
//
//	cred, err := gaccount.RegionalCredential(ctx)
//	return firestoreCreateConnection(cred)
//
func firestoreCreateConnection(cred *google.Credentials) (data.Connection, error) {
	//use context.Background() cause client will be reuse
	client, err := firestore.NewClient(context.Background(), cred.ProjectID, option.WithCredentials(cred))
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
//	id := errorID("tablename", id)
//	So(id, ShouldEqual, "tablename-id")
//
func errorID(tablename, name string) string {
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
func (c *ConnectionFirestore) snapshotToObject(tablename string, docRef *firestore.DocumentRef, snapshot *firestore.DocumentSnapshot, object data.Object) error {
	if snapshot == nil {
		return errors.New("snapshot can not be nil: " + errorID(tablename, ""))
	}

	if err := snapshot.DataTo(object); err != nil {
		return errors.Wrap(err, "failed to convert document to object: "+errorID(tablename, object.GetID()))
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
	//c.tx = nil
	//	if c.client != nil {
	//		c.client.Close()
	//		c.client = nil
	//	}
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
func (c *ConnectionFirestore) Get(ctx context.Context, tablename, id string, factory func() data.Object) (data.Object, error) {
	if id == "" {
		return nil, nil
	}
	docRef := c.getDocRef(tablename, id)
	snapshot, err := docRef.Get(ctx)

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
func (c *ConnectionFirestore) Set(ctx context.Context, tablename string, object data.Object) error {
	if object == nil {
		return errors.New("object can not be nil: " + errorID(tablename, ""))
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
	if c.batch != nil {
		c.batch.Set(docRef, object)
	} else {
		_, err = docRef.Set(ctx, object)
	}
	if err != nil {
		return errors.Wrap(err, "failed to set object: "+errorID(tablename, object.GetID()))
	}
	return nil
}

// IsExists return true if object with id exist
//
//	return c.IsExists(ctx, tablename, id)
//
func (c *ConnectionFirestore) IsExists(ctx context.Context, tablename, id string) (bool, error) {
	if id == "" {
		return false, nil
	}
	docRef := c.getDocRef(tablename, id)
	snapshot, err := docRef.Get(ctx)

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
func (c *ConnectionFirestore) All(ctx context.Context, tablename string, factory func() data.Object) ([]data.Object, error) {
	collectionRef := c.getCollectionRef(tablename)
	list := []data.Object{}
	iter := collectionRef.Limit(data.LimitQueryDefault).Documents(ctx)
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
	snapshot, err := docRef.Get(ctx)

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
func (c *ConnectionFirestore) Update(ctx context.Context, tablename, id string, fields map[string]interface{}) error {
	docRef := c.getDocRef(tablename, id)

	if c.batch != nil {
		c.batch.Set(docRef, fields, firestore.MergeAll)
		return nil
	}

	_, err := docRef.Set(ctx, fields, firestore.MergeAll)
	if err != nil {
		return errors.Wrap(err, "failed to update field: "+errorID(tablename, id))
	}
	return nil
}

// Increment value on object field, return error if object does not exist
//
//	err := c.Increment(ctx,"sample", "sample-id", "Value", 1)
//
func (c *ConnectionFirestore) Increment(ctx context.Context, tablename, id, field string, value int) error {
	docRef := c.getDocRef(tablename, id)

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
		return errors.Wrap(err, "failed to increment "+field+" with "+strconv.Itoa(value)+": "+errorID(tablename, id))
	}
	return nil
}

// Delete object using table name and id, no error if id did not exist
//
//	c.Delete(ctx, "sample", "sample-id")
//
func (c *ConnectionFirestore) Delete(ctx context.Context, tablename, id string) error {
	docRef := c.getDocRef(tablename, id)
	if c.batch != nil {
		c.batch.Delete(docRef)
		return nil
	}

	_, err := docRef.Delete(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to delete: "+errorID(tablename, id))
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
func (c *ConnectionFirestore) DeleteObject(ctx context.Context, tablename string, object data.Object) error {
	if object == nil || object.GetID() == "" {
		return nil
	}
	var docRef *firestore.DocumentRef
	if object.GetRef() == nil {
		docRef = c.getDocRef(tablename, object.GetID())
	} else {
		docRef = object.GetRef().(*firestore.DocumentRef)
	}

	if c.batch != nil {
		c.batch.Delete(docRef)

	} else {
		_, err := docRef.Delete(ctx)
		if err != nil {
			return errors.Wrap(err, "failed to delete: "+errorID(tablename, object.GetID()))
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
	for {
		// keep delete until ctx timeout or all object deleted
		if ctx.Err() != nil {
			return ctx.Err()
		}
		numDeleted := 0
		iter := collectionRef.Limit(data.LimitClear).Documents(ctx)
		defer iter.Stop()
		// Iterate through the documents, adding a delete operation for each one to a WriteBatch.
		batch := c.client.Batch()
		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return errors.Wrap(err, "failed to iterator documents: "+errorID(tablename, ""))
			}
			batch.Delete(doc.Ref)
			numDeleted++
		}
		if numDeleted > 0 {
			_, err := batch.Commit(ctx)
			if err != nil {
				return errors.Wrap(err, "failed to commit batch: "+errorID(tablename, ""))
			}
		}
		if numDeleted < data.LimitClear {
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
func (c *ConnectionFirestore) Query(tablename string, factory func() data.Object) data.Query {
	return &QueryFirestore{
		BaseQuery: data.BaseQuery{Factory: factory},
		query:     c.getCollectionRef(tablename).Query,
		conn:      c,
		tablename: tablename,
	}
}

// CreateCoder return coder from database, set numshards 100 times of concurrent usage. for example if you think concurrent use is 10/seconds then set numshards to 1000 to avoid too much retention error
//
//	productCoder,err = conn.CreateCoder("tableName","coderName",100)
//
func (c *ConnectionFirestore) CreateCoder(tableName, coderName string, numshards int) data.Coder {
	return &CoderFirestore{
		MetaFirestore: MetaFirestore{
			conn:      c.Connection.(*ConnectionFirestore),
			tableName: tableName,
			id:        coderName,
			numShards: numshards,
		},
	}
}

// CreateCounter return counter from database, create one if not exist, set numshards 100 times of concurrent usage. for example if you think concurrent use is 10/seconds then set numshards to 1000 to avoid too much retention error
// if keepDateHierarchy is true, counter will automatically generate year/month/day/hour hierarchy in utc timezone
//
//	orderCountCounter,err = conn.Counter("tableName","coderName",100,true)
//
func (c *ConnectionFirestore) CreateCounter(tableName, counterName string, numshards int, hierarchy data.DateHierarchy) data.Counter {
	if numshards <= 0 {
		numshards = 10
	}

	return &CounterFirestore{
		MetaFirestore: MetaFirestore{
			conn:      c.Connection.(*ConnectionFirestore),
			tableName: tableName,
			id:        counterName,
			numShards: numshards,
		},
		keepDateHierarchy: hierarchy == data.DateHierarchyFull,
	}
}

// CreateSerial return serial from database, create one if not exist, please be aware Serial can only generate 1 number per second, use serial with high frequency will cause too much retention error
//
//	productNo,err = conn.CreateSerial("tableName","serialName")
//
func (c *ConnectionFirestore) CreateSerial(tableName, serialName string) data.Serial {
	return &SerialFirestore{
		MetaFirestore: MetaFirestore{
			conn:      c.Connection.(*ConnectionFirestore),
			tableName: tableName,
			id:        serialName,
			numShards: 0,
		},
	}
}
