package data

import (
	"context"
	"strconv"
	"time"

	"cloud.google.com/go/firestore"
	gcp "github.com/piyuo/libsrv/gcp"
	identifier "github.com/piyuo/libsrv/identifier"
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

	// nsRef point to a namespace in database
	//
	nsRef *firestore.DocumentRef

	//tx is curenet transacton, it is nil if not in transaction
	//
	tx *firestore.Transaction
}

// Namespace separate data into different namespace in database, a database can have multiple namespace
//
type Namespace struct {
	BaseObject `firestore:"-"`
}

// FirestoreGlobalConnection create global firestore connection
//
//	ctx := context.Background()
//	conn, err := FirestoreGlobalConnection(ctx)
//	defer conn.Close()
//
func FirestoreGlobalConnection(ctx context.Context) (Connection, error) {
	cred, err := gcp.GlobalCredential(ctx)
	if err != nil {
		return nil, err
	}
	return firestoreNewConnection(ctx, cred, "")
}

// FirestoreRegionalConnection create regional database instance, regional database use namespace to sepearate data
//
//	ctx := context.Background()
//	conn, err := FirestoreRegionalConnection(ctx, "sample-namespace")
//	defer conn.Close()
//
func FirestoreRegionalConnection(ctx context.Context, namespace string) (Connection, error) {
	cred, err := gcp.RegionalCredential(ctx)
	if err != nil {
		return nil, err
	}
	return firestoreNewConnection(ctx, cred, namespace)
}

// firestoreNewConnection create connection to firestore
//
//	cred, err := gcp.RegionalCredential(ctx)
//	return firestoreNewConnection(ctx, cred, namespace)
//
func firestoreNewConnection(ctx context.Context, cred *google.Credentials, namespace string) (Connection, error) {
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

// CreateNamespace create namespace, overwrite if namespace exist
//
//	conn, err := conn.CreateNamespace(ctx)
//
func (conn *ConnectionFirestore) CreateNamespace(ctx context.Context) error {
	if conn.nsRef == nil {
		return errors.New("no namespace can be create")
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
//	err := conn.DeleteNamespace(ctx)
//
func (conn *ConnectionFirestore) DeleteNamespace(ctx context.Context) error {
	if conn.nsRef == nil {
		return errors.New("no namespace can be delete")
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

// errorID return identifier for error, identifier text is from namespace,table and object name
//
//	id := conn.errorID("tablename", "")
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
//	conn.snapshotToObject(tablename, docRef, snapshot, object)
//
func (conn *ConnectionFirestore) snapshotToObject(tablename string, docRef *firestore.DocumentRef, snapshot *firestore.DocumentSnapshot, object Object) error {
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
//	ctx := context.Background()
//	conn, err := FirestoreRegionalConnection(ctx, "sample-namespace")
//	defer conn.Close()
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
	var stopTransaction = func() {
		conn.tx = nil
	}

	return conn.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		conn.tx = tx
		defer stopTransaction()
		return callback(ctx)
	})
}

// IsInTransaction return true if connection is in transaction
//
//	inTx := conn.IsInTransaction()
//
func (conn *ConnectionFirestore) IsInTransaction() bool {
	return conn.tx != nil
}

// getCollectionRef return collection reference in table
//
//	collectionRef, err := conn.getCollectionRef(tablename)
//
func (conn *ConnectionFirestore) getCollectionRef(tablename string) *firestore.CollectionRef {
	if conn.nsRef != nil {
		return conn.nsRef.Collection(tablename)
	}
	return conn.client.Collection(tablename)
}

// getDocRef return document reference in table
//
//	docRef, err := conn.getDocRef( tablename, id)
//
func (conn *ConnectionFirestore) getDocRef(tablename, id string) *firestore.DocumentRef {
	return conn.getCollectionRef(tablename).Doc(id)
}

// Get data object from table, return nil if object does not exist
//
//	factory := func() Object {
//		return &Sample{}
//	}
//	object, err := conn.Get(ctx, "sample", id, factory)
//
func (conn *ConnectionFirestore) Get(ctx context.Context, tablename, id string, factory func() Object) (Object, error) {
	if id == "" {
		return nil, nil
	}
	docRef := conn.getDocRef(tablename, id)
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

	object := factory()
	if object == nil {
		return nil, errors.New("failed to create object from factory: " + conn.errorID(tablename, id))
	}

	err = conn.snapshotToObject(tablename, docRef, snapshot, object)
	if err != nil {
		return nil, err
	}
	return object, nil
}

// Set object into table, If the document does not exist, it will be created. If the document does exist, its contents will be overwritten with the newly provided data, if object does not have id, it will created using UUID
//
//	if err := conn.Set(ctx, tablename, object); err != nil {
//		return err
//	}
//
func (conn *ConnectionFirestore) Set(ctx context.Context, tablename string, object Object) error {
	if object == nil {
		return errors.New("object can not be nil: " + conn.errorID(tablename, ""))
	}
	var docRef *firestore.DocumentRef
	if object.GetRef() == nil { // this is new object
		if object.GetID() == "" {
			object.SetID(identifier.UUID())
		}
		docRef = conn.getDocRef(tablename, object.GetID())
		object.SetRef(docRef)
	} else {
		docRef = object.GetRef().(*firestore.DocumentRef)
	}

	var err error
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
	docRef := conn.getDocRef(tablename, id)
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

// All return max 10 object, if you need more! using query instead
//
//	return conn.All(ctx, tablename, factory)
//
func (conn *ConnectionFirestore) All(ctx context.Context, tablename string, factory func() Object) ([]Object, error) {
	collectionRef := conn.getCollectionRef(tablename)
	list := []Object{}
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
		if object == nil {
			return nil, errors.New("failed to create object from factory: " + conn.errorID(tablename, ""))
		}

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
//	return conn.Select(ctx, "sample", "sample-id", "Name")
//
func (conn *ConnectionFirestore) Select(ctx context.Context, tablename, id, field string) (interface{}, error) {
	docRef := conn.getDocRef(tablename, id)
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
//	err = conn.Update(ctx, "sample", "sample-id", map[string]interface{}{
//		"Name": "helloworld",
//	})
//
func (conn *ConnectionFirestore) Update(ctx context.Context, tablename, id string, fields map[string]interface{}) error {
	docRef := conn.getDocRef(tablename, id)
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
//	err := conn.Increment(ctx,"sample", "sample-id", "Value", 1)
//
func (conn *ConnectionFirestore) Increment(ctx context.Context, tablename, id, field string, value int) error {
	docRef := conn.getDocRef(tablename, id)
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

// Delete object using table name and id, no error if id did not exist
//
//	conn.Delete(ctx, "sample", "sample-id")
//
func (conn *ConnectionFirestore) Delete(ctx context.Context, tablename, id string) error {
	docRef := conn.getDocRef(tablename, id)
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

// DeleteObject delete object, no error if id did not exist
//
//	conn.DeleteObject(ctx, "sample", object)
//
func (conn *ConnectionFirestore) DeleteObject(ctx context.Context, tablename string, object Object) error {
	if object == nil || object.GetID() == "" {
		return nil
	}
	var docRef *firestore.DocumentRef
	if object.GetRef() == nil {
		docRef = conn.getDocRef(tablename, object.GetID())
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
	object.SetCreateTime(time.Time{})
	object.SetReadTime(time.Time{})
	object.SetUpdateTime(time.Time{})
	return nil
}

// Clear delete all object table, 500 documents at a time, if in transaction only 10 documents can be delete
//
//	err := conn.Clear(ctx, tablename)
//
func (conn *ConnectionFirestore) Clear(ctx context.Context, tablename string) error {
	collectionRef := conn.getCollectionRef(tablename)
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
		return nil
	}
	for {
		//keep delete until ctx timeout or all object deleted
		if ctx.Err() != nil {
			return ctx.Err()
		}
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
//	conn.Query(ctx, "sample", factory).Execute(ctx)
//
func (conn *ConnectionFirestore) Query(tablename string, factory func() Object) Query {
	return &QueryFirestore{
		BaseQuery: BaseQuery{factory: factory},
		query:     conn.getCollectionRef(tablename).Query,
		conn:      conn,
	}
}
