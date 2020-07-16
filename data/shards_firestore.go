package data

import (
	"context"
	"strconv"

	"cloud.google.com/go/firestore"
	"github.com/pkg/errors"
	"google.golang.org/api/iterator"
)

// ShardsFirestore is root for counter and code
//
type ShardsFirestore struct {

	// conn is current firestore connection
	//
	conn *ConnectionFirestore

	// numShards is number of shards
	//
	numShards int

	// tableName  table name
	//
	tableName string `firestore:"-"`

	// id is document id in table
	//
	id string `firestore:"-"`

	// ensureShardsDocumentCanCreateDoc is used in ensureShardsDocumentRead return true mean ensureShardsDocumentWrite need create new shards document file
	//
	ensureShardsDocumentCanCreateDoc bool

	// ensureShardCanCreateShard is used in ensureShard return true mean ensureShardWrite need create new shard file
	//
	ensureShardCanCreateShard bool
}

func (c *ShardsFirestore) errorID() string {
	return c.conn.errorID(c.tableName, c.id)
}

// assert check ctx, table name, id are valid
//
func (c *ShardsFirestore) assert(ctx context.Context) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	if c.tableName == "" {
		return errors.New("table name can not be empty")
	}
	if c.id == "" {
		return errors.New("id can not be empty")
	}
	return nil
}

// getRef return docRef and shardsRef
//
func (c *ShardsFirestore) getRef() (*firestore.DocumentRef, *firestore.CollectionRef) {
	docRef := c.conn.getDocRef(c.tableName, c.id)
	shardsRef := docRef.Collection("shards")
	return docRef, shardsRef
}

// createShards create shards collection and document, this function need to be running in transaction
//
//	err = c.createShardsTX(tx)
//
func (c *ShardsFirestore) createShardsTX(tx *firestore.Transaction) error {
	if tx == nil {
		return errors.New("CreateShardsTX() need running in transaction")
	}

	docRef, shardsRef := c.getRef()
	err := tx.Set(docRef, &struct{}{}) //put empty struct
	if err != nil {
		return errors.Wrapf(err, "failed create shards doc in transaction: "+c.errorID())
	}

	for num := 0; num < c.numShards; num++ {
		shardRef := shardsRef.Doc(strconv.Itoa(num))
		if err := tx.Set(shardRef, map[string]interface{}{"N": 0}, firestore.MergeAll); err != nil {
			return errors.Wrapf(err, "failed create shared in transaction:%v", num)
		}
	}
	return nil
}

/*
// createShards create shards collection and document, it is safe to create shards as many time as you want, normally we re create shards when we need more shards
//
//	err = c.CreateShards(ctx)
//
func (c *ShardsFirestore) createShards(ctx context.Context) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	docRef, shardsRef := c.getRef()

	if c.conn.tx != nil {
		err := c.conn.tx.Set(docRef, &struct{}{}) //put empty struct
		if err != nil {
			return errors.Wrapf(err, "failed create shards doc in transaction: "+c.errorID())
		}

		for num := 0; num < c.numShards; num++ {
			shardRef := shardsRef.Doc(strconv.Itoa(num))
			if err := c.conn.tx.Set(shardRef, map[string]interface{}{"N": 0}, firestore.MergeAll); err != nil {
				return errors.Wrapf(err, "failed create shared in transaction:%v", num)
			}
		}
		return nil
	}

	_, err := docRef.Set(ctx, c)
	if err != nil {
		return errors.Wrapf(err, "failed create shards doc: "+c.errorID())
	}

	for num := 0; num < c.numShards; num++ {
		shardRef := shardsRef.Doc(strconv.Itoa(num))
		if err := c.createShard(ctx, shardRef, num); err != nil {
			return errors.Wrapf(err, "failed create shared:%v", num)
		}
	}
	return nil
}

// createShard check and create single shard if shard is not exist
//
//	err = c.CreateShards(ctx)
//
func (c *ShardsFirestore) createShard(ctx context.Context, shardRef *firestore.DocumentRef, num int) error {
	if c.conn.tx != nil {
		snapshot, err := c.conn.tx.Get(shardRef)
		if snapshot != nil && !snapshot.Exists() {
			if err := c.conn.tx.Set(shardRef, map[string]interface{}{"N": 0}, firestore.MergeAll); err != nil {
				return errors.Wrapf(err, "failed create shared in transaction:%v", num)
			}
			return nil
		}
		if err != nil {
			return errors.Wrapf(err, "failed get shared in transaction:%v", num)
		}
		// do nothing if shard already exist
		return nil
	}
	snapshot, err := shardRef.Get(ctx)
	if snapshot != nil && !snapshot.Exists() {
		_, err := shardRef.Set(ctx, map[string]interface{}{"N": 0}, firestore.MergeAll)
		if err != nil {
			return errors.Wrapf(err, "failed create shared:%v", num)
		}
		return nil
	}
	if err != nil {
		return errors.Wrapf(err, "failed get shared:%v", num)
	}
	// do nothing if shard already exist
	return nil
}
*/

// deleteDoc delete shards document
//
//	err = c.deleteDoc(ctx)
//
func (c *ShardsFirestore) deleteDoc(ctx context.Context) error {
	docRef, _ := c.getRef()
	if c.conn.tx != nil {
		if err := c.conn.tx.Delete(docRef); err != nil {
			return err
		}
		return nil
	}
	if _, err := docRef.Delete(ctx); err != nil {
		return err
	}
	return nil
}

// deleteShards delete counter
//
//	err = c.deleteShards(ctx)
//
func (c *ShardsFirestore) deleteShards(ctx context.Context) error {
	if c.conn.tx != nil {
		err := c.deleteShardsTx(ctx, c.conn.tx)
		if err != nil {
			return errors.Wrap(err, "failed to delete shards in transaction: "+c.errorID())
		}
		return nil
	}
	err := c.conn.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		return c.deleteShardsTx(ctx, tx)
	})
	if err != nil {
		return errors.Wrap(err, "failed to delete shards: "+c.errorID())
	}
	return nil
}

// deleteShardsTx will delete shards collection and document
//
//	err = c.deleteShardsTx(ctx, tx)
//
func (c *ShardsFirestore) deleteShardsTx(ctx context.Context, tx *firestore.Transaction) error {
	docRef, shardsRef := c.getRef()
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

// count is a debug function it return shards document and collection count
//
//	docCount,shardsCount,err = c.shardsInfo(ctx)
//
func (c *ShardsFirestore) shardsInfo(ctx context.Context) (int, int, error) {
	docRef, shardsRef := c.getRef()

	docCount := 0
	snapshot, err := docRef.Get(ctx)
	if snapshot != nil && !snapshot.Exists() {
		return 0, 0, nil
	}
	if err != nil {
		return 0, 0, errors.Wrap(err, "failed to get document: "+c.errorID())
	}
	docCount++
	shardsCount := 0
	shards := shardsRef.Documents(ctx)
	defer shards.Stop()
	for {
		_, err := shards.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return 0, 0, err
		}
		shardsCount++
	}
	return docCount, shardsCount, nil
}

// ensureShardsDocumentRX make sure shards document are created, this function perform read operation, because need avoid read after write in transaction
//
//	err = c.ensureShardsDocumentRX(ctx, tx)
//
func (c *ShardsFirestore) ensureShardsDocumentRX(tx *firestore.Transaction, docRef *firestore.DocumentRef) error {
	// do not use conn.tx , use tx instead , cause caller may create transaction
	c.ensureShardsDocumentCanCreateDoc = false
	snapshot, err := tx.Get(docRef)
	if snapshot != nil && !snapshot.Exists() {
		c.ensureShardsDocumentCanCreateDoc = true
		return nil
	}
	if err != nil {
		return errors.Wrap(err, "failed to get shards document in transaction: "+c.errorID())
	}
	return nil
}

// ensureShardsDocumentWX make sure shards document are created, this function perform write operation, because need avoid read after write in transaction
//
//	err = c.ensureShardsDocumentTx(ctx, tx)
//
func (c *ShardsFirestore) ensureShardsDocumentWX(tx *firestore.Transaction, docRef *firestore.DocumentRef) error {
	if !c.ensureShardsDocumentCanCreateDoc {
		//shards document already exist
		return nil
	}
	err := tx.Set(docRef, &struct{}{}) //put empty struct
	if err != nil {
		return errors.Wrap(err, "failed to create shards document in transaction: "+c.errorID())
	}
	return nil
}

// ensureShardsDocument make sure shards document are created
//
//	err = c.ensureShardsDocument(ctx, tx)
//
func (c *ShardsFirestore) ensureShardsDocument(ctx context.Context, docRef *firestore.DocumentRef) error {
	snapshot, err := docRef.Get(ctx)
	if snapshot != nil && !snapshot.Exists() {
		_, err := docRef.Set(ctx, &struct{}{}) //put empty struct
		if err != nil {
			return errors.Wrap(err, "failed to create shards document: "+c.errorID())
		}
		return nil

	}
	if err != nil {
		return errors.Wrap(err, "failed to get shards document: "+c.errorID())
	}
	return nil
}

// ensureShardRX make sure shard and document are created, this function perform read operation, because need avoid read after write in transaction
//
//	err := c.ensureShardRX(tx,dofRef,shardRef)
//
func (c *ShardsFirestore) ensureShardRX(tx *firestore.Transaction, docRef *firestore.DocumentRef, shardRef *firestore.DocumentRef) error {
	snapshot, err := c.conn.tx.Get(shardRef)
	c.ensureShardCanCreateShard = false
	if snapshot != nil && !snapshot.Exists() {

		if err := c.ensureShardsDocumentRX(tx, docRef); err != nil {
			return err
		}
		c.ensureShardCanCreateShard = true
		return nil
	}
	if err != nil {
		return errors.Wrap(err, "failed to get shard in transaction: "+c.errorID())
	}
	//shard already exist
	return nil
}

// ensureShardTx make sure shard and document are created, this function perform write operation, because need avoid read after write in transaction
//
//	err := c.ensureShardWX(tx,docRef,shardRef)
//
func (c *ShardsFirestore) ensureShardWX(tx *firestore.Transaction, docRef *firestore.DocumentRef, shardRef *firestore.DocumentRef) error {
	if !c.ensureShardCanCreateShard {
		//shard already exist
		return nil
	}

	err := tx.Set(shardRef, map[string]interface{}{"N": 0}, firestore.MergeAll)
	if err != nil {
		return errors.Wrap(err, "failed to create shard in transaction: "+c.errorID())
	}
	return nil
}

// ensureShard make sure shard and document are created
//
//	err = c.ensureShard(ctx, tx)
//
func (c *ShardsFirestore) ensureShard(ctx context.Context, docRef *firestore.DocumentRef, shardRef *firestore.DocumentRef) error {
	snapshot, err := shardRef.Get(ctx)
	if snapshot != nil && !snapshot.Exists() {
		if err := c.ensureShardsDocument(ctx, docRef); err != nil {
			return err
		}

		_, err = shardRef.Set(ctx, map[string]interface{}{"N": 0}, firestore.MergeAll)
		if err != nil {
			return errors.Wrap(err, "failed to create shard: "+c.errorID())
		}
		return nil
	}

	if err != nil {
		return errors.Wrap(err, "failed to get shard: "+c.errorID())
	}

	//shard already exist
	return nil
}
