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
	Connection
	// nsRef point to a namespace in database
	//
	nsRef  *firestore.DocumentRef
	client *firestore.Client
	tx     *firestore.Transaction
}

// Namespace separate data into different namespace in database, a database can have multiple namespace
//
type Namespace struct {
	DocObject `firestore:"-"`
}

// FirestoreGlobalConnection create global firestore connection
//
//	ctx := context.Background()
//	conn, err := FirestoreGlobalConnection(ctx, "")
//	defer conn.Close()
//
func FirestoreGlobalConnection(ctx context.Context, namespace string) (Connection, error) {
	cred, err := gcp.GlobalCredential(ctx)
	if err != nil {
		return nil, err
	}
	return firestoreNewConnection(ctx, cred, namespace)
}

// FirestoreRegionalConnection create regional database instance
//
//	ctx := context.Background()
//	conn, err := FirestoreRegionalConnection(ctx, "sample-database")
//	defer conn.Close()
//
func FirestoreRegionalConnection(ctx context.Context, namespace string) (Connection, error) {
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
func firestoreNewConnection(ctx context.Context, cred *google.Credentials, namespace string) (Connection, error) {
	client, err := firestore.NewClient(ctx, cred.ProjectID, option.WithCredentials(cred))
	if err != nil {
		return nil, err
	}

	db := &ConnectionFirestore{
		client: client,
	}
	if namespace != "" {
		db.nsRef = client.Collection("namespace").Doc(namespace)

	}
	return db, nil
}

// CreateNamespace create namespace, create new one if not exist
//
//	dbRef, err := conn.CreateNamespace(ctx)
//
func (db *ConnectionFirestore) CreateNamespace(ctx context.Context) error {
	if db.nsRef != nil {
		_, err := db.nsRef.Set(ctx, &Namespace{})
		if err != nil {
			return errors.Wrap(err, "failed to create namespace: "+db.nsRef.ID)
		}
	}
	return nil
}

// DeleteNamespace delete namespace
//
//	err := db.DeleteNamespace(ctx)
//
func (db *ConnectionFirestore) DeleteNamespace(ctx context.Context) error {
	if db.nsRef != nil {
		_, err := db.nsRef.Delete(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

// errorID return error identifier from database name,table name and object name
//
//	id := firestoreDB.errorID("tablename", "")
//	So(id, ShouldEqual, "tablename{sample-database}")
//
func (db *ConnectionFirestore) errorID(tablename, name string) string {
	id := "{root}"
	if db.nsRef != nil {
		id = "{" + db.nsRef.ID + "}"
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
func (db *ConnectionFirestore) snapshotToObject(tablename string, docRef *firestore.DocumentRef, snapshot *firestore.DocumentSnapshot, object Object) error {
	if snapshot == nil {
		return errors.New("snapshot can not be nil: " + db.errorID(tablename, ""))
	}

	if err := snapshot.DataTo(object); err != nil {
		return errors.Wrap(err, "failed to convert document to object: "+db.errorID(tablename, object.ID()))
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
func (db *ConnectionFirestore) Close() {
	db.tx = nil
	if db.client != nil {
		db.client.Close()
		db.client = nil
	}
}

// Transaction start a transaction
//
//	err := conn.Transaction(ctx, func(ctx context.Context) error {
//		return nil
//	})
//
func (db *ConnectionFirestore) Transaction(ctx context.Context, callback func(ctx context.Context) error) error {
	return db.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		db.tx = tx
		defer db.stopTransaction()
		return callback(ctx)
	})
}

// removeTransaction stop running transaction on connection
//
//	defer db.stopTransaction()
//
func (db *ConnectionFirestore) stopTransaction() {
	db.tx = nil
}

// getCollectionRef return collection reference in current database
//
//	collectionRef, err := db.getCollectionRef(ctx, tablename)
//
func (db *ConnectionFirestore) getCollectionRef(ctx context.Context, tablename string) *firestore.CollectionRef {
	if db.nsRef != nil {
		return db.nsRef.Collection(tablename)
	}
	return db.client.Collection(tablename)
}

// getDocRef return document reference in current database
//
//	docRef, err := db.getDocRef(ctx, tablename, id)
//
func (db *ConnectionFirestore) getDocRef(ctx context.Context, tablename, id string) *firestore.DocumentRef {
	return db.getCollectionRef(ctx, tablename).Doc(id)
}

// Get data object from data store, return nil if object does not exist
//
//	object, err := conn.Get(ctx, tablename, id, factory)
//
func (db *ConnectionFirestore) Get(ctx context.Context, tablename, id string, factory func() Object) (Object, error) {
	if id == "" {
		return nil, nil
	}

	docRef := db.getDocRef(ctx, tablename, id)

	var err error
	var docSnapshot *firestore.DocumentSnapshot
	if db.tx != nil {
		docSnapshot, err = db.tx.Get(docRef)
	} else {
		docSnapshot, err = docRef.Get(ctx)
	}

	if docSnapshot != nil && !docSnapshot.Exists() {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to get document: "+db.errorID(tablename, id))
	}

	object := factory()
	if object == nil {
		return nil, errors.Wrap(err, "failed to create object from factory: "+db.errorID(tablename, id))
	}
	err = db.snapshotToObject(tablename, docRef, docSnapshot, object)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert document to object: "+db.errorID(tablename, id))
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
func (db *ConnectionFirestore) Set(ctx context.Context, tablename string, object Object) error {
	if object == nil {
		return errors.New("object can not be nil: " + db.errorID(tablename, ""))
	}
	var err error
	var docRef *firestore.DocumentRef
	if object.Ref() == nil { // this is new object
		if object.ID() == "" {
			object.SetID(util.UUID())
		}
		docRef = db.getDocRef(ctx, tablename, object.ID())
		if err != nil {
			return err
		}
		object.SetRef(docRef)
	} else {
		docRef = object.Ref().(*firestore.DocumentRef)
	}

	if db.tx != nil {
		err = db.tx.Set(docRef, object)
	} else {
		_, err = docRef.Set(ctx, object)
	}
	if err != nil {
		return errors.Wrap(err, "failed to set object: "+db.errorID(tablename, object.ID()))
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
func (db *ConnectionFirestore) Exist(ctx context.Context, tablename, id string) (bool, error) {
	if id == "" {
		return false, nil
	}
	var err error
	docRef := db.getDocRef(ctx, tablename, id)
	var snapshot *firestore.DocumentSnapshot
	if db.tx != nil {
		snapshot, err = db.tx.Get(docRef)
	} else {
		snapshot, err = docRef.Get(ctx)
	}
	if snapshot != nil && !snapshot.Exists() {
		return false, nil
	}
	if err != nil {
		return false, errors.Wrap(err, "failed to get document: "+db.errorID(tablename, id))
	}
	return true, nil
}

// List return max 10 object, if you need more! using query instead
//
//	return conn.List(ctx, tablename, factory)
//
func (db *ConnectionFirestore) List(ctx context.Context, tablename string, factory func() Object) ([]Object, error) {
	collectionRef := db.getCollectionRef(ctx, tablename)
	list := []Object{}
	var iter *firestore.DocumentIterator
	if db.tx != nil {
		iter = db.tx.Documents(collectionRef.Query.Limit(10))
	} else {
		iter = collectionRef.Limit(10).Documents(ctx)
	}
	defer iter.Stop()

	for {
		snapshot, err := iter.Next()
		if err == iterator.Done {
			break
		} else if err != nil {
			return nil, errors.Wrap(err, "failed to iterator documents: "+db.errorID(tablename, ""))
		}
		object := factory()
		err = snapshot.DataTo(object)
		if err != nil {
			return nil, errors.Wrap(err, "failed to convert document to object: "+db.errorID(tablename, ""))
		}
		db.snapshotToObject(tablename, snapshot.Ref, snapshot, object)
		list = append(list, object)
	}
	return list, nil
}

// Select return object field from data store, return nil if object does not exist
//
//	return conn.Select(ctx, tablename, id, field)
//
func (db *ConnectionFirestore) Select(ctx context.Context, tablename, id, field string) (interface{}, error) {
	docRef := db.getDocRef(ctx, tablename, id)
	var err error
	var snapshot *firestore.DocumentSnapshot
	if db.tx != nil {
		snapshot, err = db.tx.Get(docRef)
	} else {
		snapshot, err = docRef.Get(ctx)
	}

	if snapshot != nil && !snapshot.Exists() {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to get document: "+db.errorID(tablename, id))
	}
	value, err := snapshot.DataAt(field)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get value from document: "+db.errorID(tablename, id))
	}
	return value, nil
}

// Update partial object field, create new one if object does not exist,  this function is significant slow than Set()
//
//	err = conn.Update(ctx, tablename, greet.ID(), map[string]interface{}{
//		"Description": "helloworld",
//	})
//
func (db *ConnectionFirestore) Update(ctx context.Context, tablename, id string, fields map[string]interface{}) error {
	docRef := db.getDocRef(ctx, tablename, id)
	if db.tx != nil {
		err := db.tx.Set(docRef, fields, firestore.MergeAll)
		if err != nil {
			return errors.Wrap(err, "failed to update field in transaction: "+db.errorID(tablename, id))
		}
		return nil
	}
	_, err := docRef.Set(ctx, fields, firestore.MergeAll)
	if err != nil {
		return errors.Wrap(err, "failed to update field: "+db.errorID(tablename, id))
	}
	return nil
}

// Increment value on object field, return error if object does not exist
//
//	err := conn.Increment(ctx,"", GreetModelName, greet.ID(), "Value", 2)
//
func (db *ConnectionFirestore) Increment(ctx context.Context, tablename, id, field string, value int) error {
	docRef := db.getDocRef(ctx, tablename, id)
	if db.tx != nil {
		err := db.tx.Update(docRef, []firestore.Update{
			{Path: field, Value: firestore.Increment(value)},
		})
		if err != nil {
			return errors.Wrap(err, "failed to increment "+field+" with "+strconv.Itoa(value)+" in transaction: "+db.errorID(tablename, id))
		}
		return nil
	}
	_, err := docRef.Update(ctx, []firestore.Update{
		{Path: field, Value: firestore.Increment(value)},
	})
	if err != nil {
		return errors.Wrap(err, "failed to increment "+field+" with "+strconv.Itoa(value)+": "+db.errorID(tablename, id))
	}
	return nil
}

// Delete object using table name and id, no error if id not exist
//
//	conn.Delete(ctx, tablename, id)
//
func (db *ConnectionFirestore) Delete(ctx context.Context, tablename, id string) error {
	docRef := db.getDocRef(ctx, tablename, id)
	if db.tx != nil {
		err := db.tx.Delete(docRef)
		if err != nil {
			return errors.Wrap(err, "failed to delete in transaction: "+db.errorID(tablename, id))
		}
		return nil
	}
	_, err := docRef.Delete(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to delete: "+db.errorID(tablename, id))
	}
	return nil
}

// DeleteObject delete object, no error if id not exist
//
//	conn.DeleteObject(ctx, dt.tablename, object)
//
func (db *ConnectionFirestore) DeleteObject(ctx context.Context, tablename string, obj Object) error {
	if obj.ID() == "" {
		return nil
	}
	var docRef *firestore.DocumentRef
	if obj.Ref() == nil {
		docRef = db.getDocRef(ctx, tablename, obj.ID())
	} else {
		docRef = obj.Ref().(*firestore.DocumentRef)
	}

	if db.tx != nil {
		err := db.tx.Delete(docRef)
		if err != nil {
			return errors.Wrap(err, "failed to delete in transaction: "+db.errorID(tablename, obj.ID()))
		}
	} else {
		_, err := docRef.Delete(ctx)
		if err != nil {
			return errors.Wrap(err, "failed to delete: "+db.errorID(tablename, obj.ID()))
		}
	}
	obj.SetRef(nil)
	obj.SetID("")
	return nil
}

// Clear delete all object in specific time, 500 documents at a time, return false if still has object need to be delete
//	if in transaction , only 10 documents can be delete
//
//	err := conn.Clear(ctx, tablename)
//
func (db *ConnectionFirestore) Clear(ctx context.Context, tablename string) error {
	collectionRef := db.getCollectionRef(ctx, tablename)
	for {
		if db.tx != nil {
			iter := db.tx.Documents(collectionRef.Query.Limit(10))
			defer iter.Stop()
			for {
				doc, err := iter.Next()
				if err == iterator.Done {
					break
				}
				if err != nil {
					return errors.Wrap(err, "failed to iterator documents: "+db.errorID(tablename, ""))
				}
				db.tx.Delete(doc.Ref)
			}
			break
		} else {
			numDeleted := 0
			iter := collectionRef.Limit(500).Documents(ctx)
			defer iter.Stop()
			// Iterate through the documents, adding a delete operation for each one to a WriteBatch.
			batch := db.client.Batch()
			for {
				doc, err := iter.Next()
				if err == iterator.Done {
					break
				}
				if err != nil {
					return errors.Wrap(err, "failed to iterator documents: "+db.errorID(tablename, ""))
				}
				batch.Delete(doc.Ref)
				numDeleted++
			}
			if numDeleted > 0 {
				_, err := batch.Commit(ctx)
				if err != nil {
					return errors.Wrap(err, "failed to commit batch: "+db.errorID(tablename, ""))
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
func (db *ConnectionFirestore) Query(ctx context.Context, tablename string, factory func() Object) Query {
	return &QueryFirestore{
		DocQuery: DocQuery{factory: factory},
		query:    db.getCollectionRef(ctx, tablename).Query,
		tx:       db.tx,
	}
}

// Counter return counter from data store, create one if not exist
//
//	counter,err = conn.Counter(ctx, tablename, countername, numshards)
//
func (db *ConnectionFirestore) Counter(ctx context.Context, tablename, countername string, numShards int) (Counter, error) {
	docRef := db.getDocRef(ctx, tablename, countername)
	counter := &CounterFirestore{
		nsRef:       db.nsRef,
		tablename:   tablename,
		countername: countername,
		client:      db.client,
		docRef:      docRef,
		tx:          db.tx,
	}

	snapshot, err := docRef.Get(ctx)
	if snapshot != nil && !snapshot.Exists() {

		if db.tx != nil {
			err := counter.create(ctx, db.tx, docRef, counter, numShards)
			if err != nil {
				return nil, errors.Wrap(err, "failed to set counter in current transaction: "+db.errorID(tablename, countername))
			}
		} else {
			err := db.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
				return counter.create(ctx, tx, docRef, counter, numShards)
			})
			if err != nil {
				return nil, errors.Wrap(err, "failed to set counter in new transaction: "+db.errorID(tablename, countername))
			}
		}
		counter.SetCreateTime(time.Now())
		counter.SetUpdateTime(time.Now())
		counter.SetReadTime(time.Now())
		return counter, nil
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to get counter: "+db.errorID(tablename, countername))
	}
	if err := snapshot.DataTo(counter); err != nil {
		return nil, errors.Wrap(err, "failed to convert document to counter: "+db.errorID(tablename, countername))
	}
	counter.SetCreateTime(snapshot.CreateTime)
	counter.SetUpdateTime(snapshot.UpdateTime)
	counter.SetReadTime(snapshot.ReadTime)
	return counter, nil
}
