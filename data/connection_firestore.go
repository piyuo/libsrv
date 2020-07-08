package data

import (
	"context"
	"strconv"
	"time"

	"cloud.google.com/go/firestore"
	gcp "github.com/piyuo/libsrv/secure/gcp"
	util "github.com/piyuo/libsrv/util"
	"github.com/pkg/errors"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// ConnectionFirestore implement firestore connection
//
type ConnectionFirestore struct {
	ConnectionRef
	// nsRef point to a namespace in database
	//
	nsRef  *firestore.DocumentRef
	client *firestore.Client
	tx     *firestore.Transaction
}

// Namespace separate data into different namespace in database, a database can have multiple namespace
//
type Namespace struct {
	Object `firestore:"-"`
}

// FirestoreGlobalConnection create global firestore connection
//
//	ctx := context.Background()
//	conn, err := FirestoreGlobalConnection(ctx, "")
//	defer conn.Close()
//
func FirestoreGlobalConnection(ctx context.Context, namespace string) (ConnectionRef, error) {
	cred, err := gcp.GlobalCredential(ctx)
	if err != nil {
		return nil, err
	}
	return firestoreNewConnection(ctx, cred, namespace)
}

// FirestoreRegionalConnection create regional database instance
//
//	ctx := context.Background()
//	conn, err := FirestoreRegionalConnection(ctx, "sample-namespace")
//	defer conn.Close()
//
func FirestoreRegionalConnection(ctx context.Context, namespace string) (ConnectionRef, error) {
	cred, err := gcp.CurrentRegionalCredential(ctx)
	if err != nil {
		return nil, err
	}
	return firestoreNewConnection(ctx, cred, namespace)
}

// firestoreNewConnection create connection to firestore
//
//	cred, err := gcp.CurrentRegionalCredential(ctx)
//	return firestoreNewConnection(ctx, cred, namespace)
//
func firestoreNewConnection(ctx context.Context, cred *google.Credentials, namespace string) (ConnectionRef, error) {
	client, err := firestore.NewClient(ctx, cred.ProjectID, option.WithCredentials(cred))
	if err != nil {
		return nil, err
	}

	conn := &ConnectionFirestore{
		client: client,
	}
	if namespace != "" {
		conn.nsRef = client.Collection("namespace").Doc(namespace)

	}
	return conn, nil
}

// CreateNamespace create namespace, create new one if not exist
//
//	dbRef, err := conn.CreateNamespace(ctx)
//
func (conn *ConnectionFirestore) CreateNamespace(ctx context.Context) error {
	if conn.nsRef == nil {
		return nil
	}

	var err error
	if conn.tx != nil {
		err = conn.tx.Set(conn.nsRef, &Namespace{})
	} else {
		_, err = conn.nsRef.Set(ctx, &Namespace{})
	}
	if err != nil {
		return errors.Wrap(err, "failed to create namespace: "+conn.nsRef.ID)
	}
	return nil
}

// DeleteNamespace delete namespace
//
//	err := db.DeleteNamespace(ctx)
//
func (conn *ConnectionFirestore) DeleteNamespace(ctx context.Context) error {
	if conn.nsRef == nil {
		return nil
	}

	var err error
	if conn.tx != nil {
		err = conn.tx.Delete(conn.nsRef)
	} else {
		_, err = conn.nsRef.Delete(ctx)
	}
	if err != nil {
		return errors.Wrap(err, "failed to delete namespace: "+conn.nsRef.ID)
	}
	return nil
}

// errorID return error identifier from database name,table name and object name
//
//	id := firestoreDB.errorID("tablename", "")
//	So(id, ShouldEqual, "tablename{sample-namespace}")
//
func (conn *ConnectionFirestore) errorID(tablename, name string) string {
	id := "{root}"
	if conn.nsRef != nil {
		id = "{" + conn.nsRef.ID + "}"
	}
	id = tablename + id
	if name != "" {
		id += "-" + name
	}
	return id
}

// snapshotToObject convert document snapshot to object
//
//	db.snapshotToObject(tablename, docRef, docSnapshot, object)
//
func (conn *ConnectionFirestore) snapshotToObject(tablename string, docRef *firestore.DocumentRef, snapshot *firestore.DocumentSnapshot, object ObjectRef) error {
	if snapshot == nil {
		return errors.New("snapshot can not be nil: " + conn.errorID(tablename, ""))
	}

	if err := snapshot.DataTo(object); err != nil {
		return errors.Wrap(err, "failed to convert document to object: "+conn.errorID(tablename, object.GetID()))
	}
	object.SetRef(docRef)
	object.SetID(docRef.ID)
	object.SetCreateTime(snapshot.CreateTime)
	object.SetUpdateTime(snapshot.UpdateTime)
	object.SetReadTime(snapshot.ReadTime)
	return nil
}

// Close database connection
//
//	conn.Close()
//
func (conn *ConnectionFirestore) Close() {
	conn.tx = nil
	conn.nsRef = nil
	if conn.client != nil {
		conn.client.Close()
		conn.client = nil
	}
}

// Transaction start a transaction
//
//	err := conn.Transaction(ctx, func(ctx context.Context) error {
//		return nil
//	})
//
func (conn *ConnectionFirestore) Transaction(ctx context.Context, callback func(ctx context.Context) error) error {
	return conn.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		conn.tx = tx
		defer conn.stopTransaction()
		return callback(ctx)
	})
}

// removeTransaction stop running transaction on connection
//
//	defer db.stopTransaction()
//
func (conn *ConnectionFirestore) stopTransaction() {
	conn.tx = nil
}

// getCollectionRef return collection reference in current database
//
//	collectionRef, err := db.getCollectionRef(ctx, tablename)
//
func (conn *ConnectionFirestore) getCollectionRef(ctx context.Context, tablename string) *firestore.CollectionRef {
	if conn.nsRef != nil {
		return conn.nsRef.Collection(tablename)
	}
	return conn.client.Collection(tablename)
}

// getDocRef return document reference in current database
//
//	docRef, err := db.getDocRef(ctx, tablename, id)
//
func (conn *ConnectionFirestore) getDocRef(ctx context.Context, tablename, id string) *firestore.DocumentRef {
	return conn.getCollectionRef(ctx, tablename).Doc(id)
}

// Get data object from data store, return nil if object does not exist
//
//	object, err := conn.Get(ctx, tablename, id, factory)
//
func (conn *ConnectionFirestore) Get(ctx context.Context, tablename, id string, factory func() ObjectRef) (ObjectRef, error) {
	if id == "" {
		return nil, nil
	}
	docRef := conn.getDocRef(ctx, tablename, id)
	var err error
	var docSnapshot *firestore.DocumentSnapshot
	if conn.tx != nil {
		docSnapshot, err = conn.tx.Get(docRef)
	} else {
		docSnapshot, err = docRef.Get(ctx)
	}

	if docSnapshot != nil && !docSnapshot.Exists() {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to get document: "+conn.errorID(tablename, id))
	}

	object := factory()
	if object == nil {
		return nil, errors.Wrap(err, "failed to create object from factory: "+conn.errorID(tablename, id))
	}
	err = conn.snapshotToObject(tablename, docRef, docSnapshot, object)
	if err != nil {
		return nil, err
	}
	return object, nil
}

// Set object into data store, If the document does not exist, it will be created. If the document does exist, its contents will be overwritten with the newly provided data,
//
// if object does not have id, it will created using UUID
//
//	if err := conn.Set(ctx, tablename, object); err != nil {
//		return err
//	}
//
func (conn *ConnectionFirestore) Set(ctx context.Context, tablename string, object ObjectRef) error {
	if object == nil {
		return errors.New("object can not be nil: " + conn.errorID(tablename, ""))
	}
	var err error
	var docRef *firestore.DocumentRef
	if object.GetRef() == nil { // this is new object
		if object.GetID() == "" {
			object.SetID(util.UUID())
		}
		docRef = conn.getDocRef(ctx, tablename, object.GetID())
		if err != nil {
			return err
		}
		object.SetRef(docRef)
	} else {
		docRef = object.GetRef().(*firestore.DocumentRef)
	}

	if conn.tx != nil {
		err = conn.tx.Set(docRef, object)
	} else {
		_, err = docRef.Set(ctx, object)
	}
	if err != nil {
		return errors.Wrap(err, "failed to set object: "+conn.errorID(tablename, object.GetID()))
	}
	object.SetCreateTime(time.Now())
	object.SetUpdateTime(time.Now())
	object.SetReadTime(time.Now())
	return nil
}

// Exist return true if object with id exist
//
//	return conn.Exist(ctx, tablename, id)
//
func (conn *ConnectionFirestore) Exist(ctx context.Context, tablename, id string) (bool, error) {
	if id == "" {
		return false, nil
	}
	var err error
	docRef := conn.getDocRef(ctx, tablename, id)
	var snapshot *firestore.DocumentSnapshot
	if conn.tx != nil {
		snapshot, err = conn.tx.Get(docRef)
	} else {
		snapshot, err = docRef.Get(ctx)
	}
	if snapshot != nil && !snapshot.Exists() {
		return false, nil
	}
	if err != nil {
		return false, errors.Wrap(err, "failed to get document: "+conn.errorID(tablename, id))
	}
	return true, nil
}

// List return max 10 object, if you need more! using query instead
//
//	return conn.List(ctx, tablename, factory)
//
func (conn *ConnectionFirestore) List(ctx context.Context, tablename string, factory func() ObjectRef) ([]ObjectRef, error) {
	collectionRef := conn.getCollectionRef(ctx, tablename)
	list := []ObjectRef{}
	var iter *firestore.DocumentIterator
	if conn.tx != nil {
		iter = conn.tx.Documents(collectionRef.Query.Limit(limitQueryDefault))
	} else {
		iter = collectionRef.Limit(limitQueryDefault).Documents(ctx)
	}
	defer iter.Stop()

	for {
		snapshot, err := iter.Next()
		if err == iterator.Done {
			break
		} else if err != nil {
			return nil, errors.Wrap(err, "failed to iterator documents: "+conn.errorID(tablename, ""))
		}
		object := factory()
		err = snapshot.DataTo(object)
		if err != nil {
			return nil, errors.Wrap(err, "failed to convert document to object: "+conn.errorID(tablename, ""))
		}
		conn.snapshotToObject(tablename, snapshot.Ref, snapshot, object)
		list = append(list, object)
	}
	return list, nil
}

// Select return object field from data store, return nil if object does not exist
//
//	return conn.Select(ctx, tablename, id, field)
//
func (conn *ConnectionFirestore) Select(ctx context.Context, tablename, id, field string) (interface{}, error) {
	docRef := conn.getDocRef(ctx, tablename, id)
	var err error
	var snapshot *firestore.DocumentSnapshot
	if conn.tx != nil {
		snapshot, err = conn.tx.Get(docRef)
	} else {
		snapshot, err = docRef.Get(ctx)
	}

	if snapshot != nil && !snapshot.Exists() {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to get document: "+conn.errorID(tablename, id))
	}
	value, err := snapshot.DataAt(field)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get value from document: "+conn.errorID(tablename, id))
	}
	return value, nil
}

// Update partial object field, create new one if object does not exist,  this function is significant slow than Set()
//
//	err = conn.Update(ctx, tablename, greet.ID, map[string]interface{}{
//		"Description": "helloworld",
//	})
//
func (conn *ConnectionFirestore) Update(ctx context.Context, tablename, id string, fields map[string]interface{}) error {
	docRef := conn.getDocRef(ctx, tablename, id)
	if conn.tx != nil {
		err := conn.tx.Set(docRef, fields, firestore.MergeAll)
		if err != nil {
			return errors.Wrap(err, "failed to update field in transaction: "+conn.errorID(tablename, id))
		}
		return nil
	}
	_, err := docRef.Set(ctx, fields, firestore.MergeAll)
	if err != nil {
		return errors.Wrap(err, "failed to update field: "+conn.errorID(tablename, id))
	}
	return nil
}

// Increment value on object field, return error if object does not exist
//
//	err := conn.Increment(ctx,"", GreetModelName, greet.ID, "Value", 2)
//
func (conn *ConnectionFirestore) Increment(ctx context.Context, tablename, id, field string, value int) error {
	docRef := conn.getDocRef(ctx, tablename, id)
	if conn.tx != nil {
		err := conn.tx.Update(docRef, []firestore.Update{
			{Path: field, Value: firestore.Increment(value)},
		})
		if err != nil {
			return errors.Wrap(err, "failed to increment "+field+" with "+strconv.Itoa(value)+" in transaction: "+conn.errorID(tablename, id))
		}
		return nil
	}
	_, err := docRef.Update(ctx, []firestore.Update{
		{Path: field, Value: firestore.Increment(value)},
	})
	if err != nil {
		return errors.Wrap(err, "failed to increment "+field+" with "+strconv.Itoa(value)+": "+conn.errorID(tablename, id))
	}
	return nil
}

// Delete object using table name and id, no error if id not exist
//
//	conn.Delete(ctx, tablename, id)
//
func (conn *ConnectionFirestore) Delete(ctx context.Context, tablename, id string) error {
	docRef := conn.getDocRef(ctx, tablename, id)
	if conn.tx != nil {
		err := conn.tx.Delete(docRef)
		if err != nil {
			return errors.Wrap(err, "failed to delete in transaction: "+conn.errorID(tablename, id))
		}
		return nil
	}
	_, err := docRef.Delete(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to delete: "+conn.errorID(tablename, id))
	}
	return nil
}

// DeleteObject delete object, no error if id not exist
//
//	conn.DeleteObject(ctx, dt.tablename, object)
//
func (conn *ConnectionFirestore) DeleteObject(ctx context.Context, tablename string, object ObjectRef) error {
	if object == nil || object.GetID() == "" {
		return nil
	}
	var docRef *firestore.DocumentRef
	if object.GetRef() == nil {
		docRef = conn.getDocRef(ctx, tablename, object.GetID())
	} else {
		docRef = object.GetRef().(*firestore.DocumentRef)
	}

	if conn.tx != nil {
		err := conn.tx.Delete(docRef)
		if err != nil {
			return errors.Wrap(err, "failed to delete in transaction: "+conn.errorID(tablename, object.GetID()))
		}
	} else {
		_, err := docRef.Delete(ctx)
		if err != nil {
			return errors.Wrap(err, "failed to delete: "+conn.errorID(tablename, object.GetID()))
		}
	}
	object.SetRef(nil)
	object.SetID("")
	return nil
}

// Clear delete all object in specific time, 500 documents at a time, return false if still has object need to be delete
//	if in transaction , only 10 documents can be delete
//
//	err := conn.Clear(ctx, tablename)
//
func (conn *ConnectionFirestore) Clear(ctx context.Context, tablename string) error {
	collectionRef := conn.getCollectionRef(ctx, tablename)
	for {
		if conn.tx != nil {
			iter := conn.tx.Documents(collectionRef.Query.Limit(limitTransactionClear))
			defer iter.Stop()
			for {
				doc, err := iter.Next()
				if err == iterator.Done {
					break
				}
				if err != nil {
					return errors.Wrap(err, "failed to iterator documents: "+conn.errorID(tablename, ""))
				}
				conn.tx.Delete(doc.Ref)
			}
			break
		} else {
			numDeleted := 0
			iter := collectionRef.Limit(limitClear).Documents(ctx)
			defer iter.Stop()
			// Iterate through the documents, adding a delete operation for each one to a WriteBatch.
			batch := conn.client.Batch()
			for {
				doc, err := iter.Next()
				if err == iterator.Done {
					break
				}
				if err != nil {
					return errors.Wrap(err, "failed to iterator documents: "+conn.errorID(tablename, ""))
				}
				batch.Delete(doc.Ref)
				numDeleted++
			}
			if numDeleted > 0 {
				_, err := batch.Commit(ctx)
				if err != nil {
					return errors.Wrap(err, "failed to commit batch: "+conn.errorID(tablename, ""))
				}
			}
			if numDeleted < 500 {
				break
			}
		}
	}
	return nil
}

// Query create query
//
//	conn.Query(ctx, tablename, factory)
//
func (conn *ConnectionFirestore) Query(ctx context.Context, tablename string, factory func() ObjectRef) QueryRef {
	return &QueryFirestore{
		Query: Query{factory: factory},
		query: conn.getCollectionRef(ctx, tablename).Query,
		tx:    conn.tx,
	}
}

// Counter return counter from data store, create one if not exist
//
//	counter,err = conn.Counter(ctx, tablename, countername, numshards)
//
func (conn *ConnectionFirestore) Counter(ctx context.Context, tablename, countername string, numShards int) (CounterRef, error) {
	docRef := conn.getDocRef(ctx, tablename, countername)
	shardsRef := docRef.Collection("shards")
	if conn.tx != nil {
		counter, err := conn.ensureCounterInTx(ctx, tablename, countername, conn.tx, docRef, shardsRef, numShards)
		if err != nil {
			return nil, errors.Wrap(err, "failed to ensure counter in current transaction: "+conn.errorID(tablename, countername))
		}
		return counter, nil
	}
	var counter *CounterFirestore
	err := conn.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		var err error
		counter, err = conn.ensureCounterInTx(ctx, tablename, countername, tx, docRef, shardsRef, numShards)
		if err != nil {
			return errors.Wrap(err, "failed to ensure counter: "+conn.errorID(tablename, countername))
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to run transaction ensure counter: "+conn.errorID(tablename, countername))
	}
	counter.tx = nil
	return counter, nil
}

// createCounterInTx create counter in transaction with a given number of shards as subcollection of specified document.
//
//	err := db.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
//		return counter.create(ctx, db.tx, docRef, counter, numShards)
//	})
//
func (conn *ConnectionFirestore) ensureCounterInTx(ctx context.Context, tablename, countername string, tx *firestore.Transaction, docRef *firestore.DocumentRef, shardsRef *firestore.CollectionRef, numShards int) (*CounterFirestore, error) {
	namespace := ""
	if conn.nsRef != nil {
		namespace = conn.nsRef.ID
	}
	counter := &CounterFirestore{
		shardsRef:   shardsRef,
		client:      conn.client,
		tx:          conn.tx,
		NameSpace:   namespace,
		TableName:   tablename,
		CounterName: countername,
	}

	snapshot, err := tx.Get(docRef)
	if snapshot != nil && !snapshot.Exists() {
		counter.N = numShards
		err = tx.Set(docRef, counter)
		if err != nil {
			return nil, err
		}
		// Initialize each shard with count=0
		for num := 0; num < numShards; num++ {
			shard := &Shard{C: 0}
			sharedRef := shardsRef.Doc(strconv.Itoa(num))
			err = tx.Set(sharedRef, shard)
			if err != nil {
				return nil, errors.Wrapf(err, "failed init counter shared:%v ", num)
			}
		}
		counter.CreateTime = time.Now()
		counter.UpdateTime = time.Now()
		counter.ReadTime = time.Now()
		return counter, nil
	}
	if err != nil {
		return nil, err
	}
	if err := snapshot.DataTo(counter); err != nil {
		return nil, err
	}
	counter.CreateTime = snapshot.CreateTime
	counter.UpdateTime = snapshot.UpdateTime
	counter.ReadTime = snapshot.ReadTime
	return counter, nil
}

// DeleteCounter delete counter
//
//	err = conn.DeleteCounter(ctx, tablename, countername)
//
func (conn *ConnectionFirestore) DeleteCounter(ctx context.Context, tablename, countername string) error {
	docRef := conn.getDocRef(ctx, tablename, countername)
	shardsRef := docRef.Collection("shards")
	if conn.tx != nil {
		err := conn.deleteCounterInTx(ctx, conn.tx, docRef, shardsRef)
		if err != nil {
			return errors.Wrap(err, "failed to set counter in current transaction: "+conn.errorID(tablename, countername))
		}
		return nil
	}
	err := conn.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		return conn.deleteCounterInTx(ctx, tx, docRef, shardsRef)
	})
	if err != nil {
		return errors.Wrap(err, "failed to set counter in new transaction: "+conn.errorID(tablename, countername))
	}
	return nil
}

// Counter return counter from data store, create one if not exist
//
//	err = conn.DeleteCounter(ctx, tablename, countername)
//
func (conn *ConnectionFirestore) deleteCounterInTx(ctx context.Context, tx *firestore.Transaction, docRef *firestore.DocumentRef, shardsRef *firestore.CollectionRef) error {
	shards := tx.Documents(shardsRef)
	defer shards.Stop()
	for {
		shard, err := shards.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}

		if err = tx.Delete(shard.Ref); err != nil {
			return err
		}

	}
	if err := tx.Delete(docRef); err != nil {
		return err
	}
	return nil
}
