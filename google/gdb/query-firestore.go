package gdb

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/piyuo/libsrv/db"
	"github.com/piyuo/libsrv/log"
	"github.com/pkg/errors"
	"google.golang.org/api/iterator"
)

// QueryFirestore implement google firestore
type QueryFirestore struct {
	db.BaseQuery

	// conn is current firestore connection
	//
	client *ClientFirestore

	// query is firestore query
	//
	query firestore.Query
}

// Where set filter, if path == "ID" mean using document id in as filter
//
//	list, err := Query(&Sample{}).Where("ID", "==", "sample1").Return(ctx)
//
func (c *QueryFirestore) Where(path, op string, value interface{}) db.Query {
	if c.QueryObject != nil && path == "ID" {
		path = firestore.DocumentID
		value = c.client.getDocRef(c.QueryObject.Collection(), value.(string))
	}
	c.query = c.query.Where(path, op, value)
	return c
}

// OrderBy set query order by asc
//
//	list, err = Query(&Sample{}).OrderBy("Name").Return(ctx)
//
func (c *QueryFirestore) OrderBy(path string) db.Query {
	c.query = c.query.OrderBy(path, firestore.Asc)
	return c
}

// OrderByDesc set query order by desc
//
//	list, err = Query(&Sample{}).OrderByDesc("Name").Limit(1).Return(ctx)
//
func (c *QueryFirestore) OrderByDesc(path string) db.Query {
	c.query = c.query.OrderBy(path, firestore.Desc)
	return c
}

// Limit set query limit
//
//	list, err = Query(&Sample{}).OrderBy("Name").Limit(1).Return(ctx)
//
func (c *QueryFirestore) Limit(n int) db.Query {
	c.query = c.query.Limit(n)
	return c
}

// StartAt implement Paginate on firestore, please be aware not use index but fieldValue to do the trick, see sample
//
//	list, err = Query(&Sample{}).OrderBy("Name").StartAt("irvine city").Return(ctx)
//
func (c *QueryFirestore) StartAt(docSnapshotOrFieldValues ...interface{}) db.Query {
	c.query = c.query.StartAt(docSnapshotOrFieldValues...)
	return c
}

// StartAfter implement Paginate on firestore, please be aware not use index but fieldValue to do the trick, see sample
//
//	list, err = Query(&Sample{}).OrderBy("Name").StartAfter("santa ana city").Return(ctx)
//
func (c *QueryFirestore) StartAfter(docSnapshotOrFieldValues ...interface{}) db.Query {
	c.query = c.query.StartAfter(docSnapshotOrFieldValues...)
	return c
}

// EndAt implement Paginate on firestore, please be aware not use index but fieldValue to do the trick, see sample
//
//	list, err = Query(&Sample{}).OrderBy("Name").EndAt("irvine city").Return(ctx)
//
func (c *QueryFirestore) EndAt(docSnapshotOrFieldValues ...interface{}) db.Query {
	c.query = c.query.EndAt(docSnapshotOrFieldValues...)
	return c
}

// EndBefore implement Paginate on firestore, please be aware not use index but fieldValue to do the trick, see sample
//
//	list, err = Query(&Sample{}).OrderBy("Name").EndBefore("irvine city").Return(ctx)
//
func (c *QueryFirestore) EndBefore(docSnapshotOrFieldValues ...interface{}) db.Query {
	c.query = c.query.EndBefore(docSnapshotOrFieldValues...)
	return c
}

func (c *QueryFirestore) returnIter(ctx context.Context) (*firestore.DocumentIterator, error) {
	if err := db.AssertObject(ctx, c.QueryObject, false); err != nil {
		return nil, err
	}

	var iter *firestore.DocumentIterator
	if c.QueryTransaction != nil {
		trans := c.QueryTransaction.(*TransactionFirestore)
		iter = trans.tx.Documents(c.query)
	} else {
		iter = c.query.Documents(ctx)
	}
	return iter, nil
}

// Return query result with default limit to 20 object, use Limit() to override default limit, return nil if anything wrong
//
//	list, err = Query(&Sample{}).OrderByDesc("Name").Limit(1).Return(ctx)
//
func (c *QueryFirestore) Return(ctx context.Context) ([]db.Object, error) {
	iter, err := c.returnIter(ctx)
	if err != nil {
		return nil, err
	}
	defer iter.Stop()
	return iterObjects(c.QueryObject, iter)
}

// ReturnID only return object id with default limit to 20 object, use Limit() to override default limit, return nil if anything wrong
//
//	idList, err := Query(&Sample{}).OrderBy("From").Limit(1).StartAt("b city").ReturnID(ctx)
//
func (c *QueryFirestore) ReturnID(ctx context.Context) ([]string, error) {
	iter, err := c.returnIter(ctx)
	if err != nil {
		return nil, err
	}
	defer iter.Stop()
	result := []string{}
	for {
		snapshot, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		result = append(result, snapshot.Ref.ID)
	}
	return result, nil
}

// ReturnCount return object count with default limit to 20 object, use Limit() to override default limit
//
//	count, err := Query(&Sample{}).Where("Name", "==", "sample1").ReturnCount(ctx)
//
func (c *QueryFirestore) ReturnCount(ctx context.Context) (int, error) {
	iter, err := c.returnIter(ctx)
	if err != nil {
		return 0, err
	}
	defer iter.Stop()

	count := 0
	for {
		_, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return 0, err
		}
		count++
	}
	return count, nil
}

// ReturnEmpty return true if no object exist
//
//	isEmpty, err := Query(&Sample{}).Where("Name", "==", "sample1").ReturnEmpty(ctx)
//
func (c *QueryFirestore) ReturnEmpty(ctx context.Context) (bool, error) {
	c.Limit(1)
	iter, err := c.returnIter(ctx)
	if err != nil {
		return false, err
	}
	defer iter.Stop()

	_, err = iter.Next()
	if err == iterator.Done {
		return true, nil
	}
	if err != nil {
		return false, err
	}
	return false, nil
}

// ReturnExists return true if object exist
//
//	isExists, err := Query(&Sample{}).Where("Name", "==", "sample1").ReturnExists(ctx)
//
func (c *QueryFirestore) ReturnExists(ctx context.Context) (bool, error) {
	empty, err := c.ReturnEmpty(ctx)
	return !empty, err
}

// ReturnFirst return first object from query
//
//	obj, err := Query(&Sample{}).OrderBy("From").Limit(1).StartAt("b city").ReturnFirst(ctx)
//	greet := obj.(*Greet)
//
func (c *QueryFirestore) ReturnFirst(ctx context.Context) (db.Object, error) {
	if err := db.AssertObject(ctx, c.QueryObject, false); err != nil {
		return nil, err
	}
	list, err := c.Limit(1).Return(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "return first")
	}
	if len(list) == 0 {
		return nil, nil
	}
	return list[0], nil
}

// ReturnFirstID return first object id from query
//
//	id, err := Query(&Sample{}).OrderBy("From").Limit(1).StartAt("b city").ReturnFirstID(ctx)
//
func (c *QueryFirestore) ReturnFirstID(ctx context.Context) (string, error) {
	if err := db.AssertObject(ctx, c.QueryObject, false); err != nil {
		return "", err
	}
	list, err := c.Limit(1).Return(ctx)
	if err != nil {
		return "", errors.Wrap(err, "return first id")
	}
	if len(list) == 0 {
		return "", nil
	}
	return list[0].ID(), nil
}

// Delete delete all document return from query. delete max doc count. return is done,delete count, error
//
//	done, count, err := client.Query(&Sample{}).Where("Name", "==", name).Delete(ctx, 100)
//
func (c *QueryFirestore) Delete(ctx context.Context, max int) (bool, int, error) {
	if ctx.Err() != nil {
		return false, 0, ctx.Err()
	}
	if c.QueryTransaction != nil {
		return false, 0, errors.New("delete query is not support in transaction, use tx.Delete() instead")
		/*
			trans := c.QueryTransaction.(*TransactionFirestore)
			iter := trans.tx.Documents(c.query)
			defer iter.Stop()
			complete, numDeleted, err := trans.DeleteCollection(ctx, c.QueryObject, max, iter)
			if err != nil {
				return false, numDeleted, errors.Wrapf(err, "tx delete %v", c.QueryObject.Collection())
			}
			return complete, numDeleted, nil
		*/
	}

	c.Limit(max)
	iter := c.query.Documents(ctx)
	defer iter.Stop()
	complete, numDeletd, err := c.client.deleteByIterator(ctx, max, iter)
	if err != nil {
		return false, numDeletd, errors.Wrap(err, "delete "+c.QueryObject.Collection())
	}
	return complete, numDeletd, nil
}

// Cleanup delete 25 document a time, max 60,000 object. return true if no object left in collection
//
//	done,  err := client.Query(&Sample{}).Where("Name", "==", name).Cleanup(ctx)
//
func (c *QueryFirestore) Cleanup(ctx context.Context) (bool, error) {
	if ctx.Err() != nil {
		return false, ctx.Err()
	}
	if c.QueryTransaction != nil {
		return false, errors.New("cleanup query is not support in transaction")
	}

	numDeleted := 0
	complete := false
	for i := 0; i < 2400; i++ {
		done, count, err := c.Delete(ctx, 25)
		if err != nil {
			return false, errors.Wrap(err, "clean 25 doc")
		}
		numDeleted += count
		if done {
			complete = true
			break
		}
	}
	if numDeleted > 0 {
		log.Info(ctx, "cleanup %s, deleted %v "+c.QueryObject.Collection(), numDeleted)
	}
	return complete, nil
}
