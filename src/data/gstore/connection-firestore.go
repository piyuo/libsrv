package gstore

import (
	"context"
	"fmt"
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

// snapshotToObject2 convert document snapshot to object
//
//	c.snapshotToObject(tablename, docRef, snapshot, object)
//
func snapshotToObject2(tablename string, docRef *firestore.DocumentRef, snapshot *firestore.DocumentSnapshot, object data.Object) error {
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

func snapshotToObject(obj data.Object, id string, docRef *firestore.DocumentRef, snapshot *firestore.DocumentSnapshot, err error) (data.Object, error) {
	if snapshot != nil && !snapshot.Exists() {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Wrapf(err, "get snapshot from table %v,%v", obj.TableName(), id)
	}

	if err := snapshot.DataTo(obj); err != nil {
		return nil, errors.Wrapf(err, "make snapshot to object %v,%v", obj.TableName(), id)
	}
	obj.SetRef(docRef)
	obj.SetID(docRef.ID)
	return obj, nil
}

// Get data object from table, return nil if object does not exist
//
//	object, err := Get(ctx, &Sample{}, "id")
//
func (c *ConnectionFirestore) Get(ctx context.Context, obj data.Object, id string) (data.Object, error) {
	if obj == nil {
		return nil, errors.New(fmt.Sprintf("obj must not nil %v", id))
	}
	if id == "" {
		return nil, errors.New(fmt.Sprintf("id must not empty %v", obj.TableName()))
	}
	docRef := c.getDocRef(obj.TableName(), id)
	snapshot, err := docRef.Get(ctx)
	return snapshotToObject(obj, id, docRef, snapshot, err)
}

// Set object into table, If the document does not exist, it will be created. If the document does exist, its contents will be overwritten with the newly provided data, if object does not have id, it will created using UUID
//
//	 err := Set(ctx, object)
//
func (c *ConnectionFirestore) Set(ctx context.Context, obj data.Object) error {
	if obj == nil {
		return errors.New("Set() obj must not nil")
	}
	var docRef *firestore.DocumentRef
	if obj.Ref() == nil { // new object
		if obj.ID() == "" {
			obj.SetID(identifier.UUID())
		}
		docRef = c.getDocRef(obj.TableName(), obj.ID())
		obj.SetRef(docRef)
	} else { // object already exist
		docRef = obj.Ref().(*firestore.DocumentRef)
	}

	_, err := docRef.Set(ctx, obj)
	if err != nil {
		return errors.Wrapf(err, "Set(%v,%v)", errorID(obj.TableName(), obj.ID()))
	}
	return nil
}

// Exists return true if object with id exist
//
//	found,err := Exists(ctx, &Sample{}, "id")
//
func (c *ConnectionFirestore) Exists(ctx context.Context, obj data.Object, id string) (bool, error) {
	if obj == nil {
		return false, errors.New("Exists() obj must not nil")
	}
	if id == "" {
		return false, errors.New(fmt.Sprintf("Exists(%v) id must not empty", obj.TableName()))
	}
	docRef := c.getDocRef(obj.TableName(), id)
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
//	list,err := All(ctx, &Sample{})
//
func (c *ConnectionFirestore) All(ctx context.Context, obj data.Object) ([]data.Object, error) {
	collectionRef := c.getCollectionRef(obj.TableName())
	list := []data.Object{}
	iter := collectionRef.Limit(data.LimitQueryDefault).Documents(ctx)
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
	_, err := docRef.Delete(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to delete: "+errorID(tablename, id))
	}
	return nil
}

// DeleteObject delete object, no error if id did not exist
//
//	c.DeleteObject(ctx, "sample", object)
//
func (c *ConnectionFirestore) DeleteObject(ctx context.Context, tablename string, object data.Object) error {
	if object == nil || object.ID() == "" {
		return nil
	}
	var docRef *firestore.DocumentRef
	if object.Ref() == nil {
		docRef = c.getDocRef(tablename, object.ID())
	} else {
		docRef = object.Ref().(*firestore.DocumentRef)
	}

	_, err := docRef.Delete(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to delete: "+errorID(tablename, object.ID()))
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

// CreateTransaction create transaction
//
func (c *ConnectionFirestore) CreateTransaction() data.Transaction {
	return &TransactionFirestore{
		conn: c,
	}
}

// CreateBatch create batch
//
func (c *ConnectionFirestore) CreateBatch() data.Batch {
	return &BatchFirestore{
		conn: c,
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
