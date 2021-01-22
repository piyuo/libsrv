package data

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/pkg/errors"
	"google.golang.org/api/iterator"
)

// QueryFirestore implement google firestore
type QueryFirestore struct {
	BaseQuery

	// conn is current firestore connection
	//
	conn *ConnectionFirestore

	// query is firestore query
	//
	query firestore.Query

	// table name is  query target table name
	//
	tablename string
}

// Where set filter, if path == "ID" mean using document id in as filter
//
//	list, err := table.Query().Where("ID", "==", "sample1").Execute(ctx)
//	So((list[0].(*Sample)).Name, ShouldEqual, "sample1")
//
func (c *QueryFirestore) Where(path, op string, value interface{}) Query {
	if path == "ID" {
		path = firestore.DocumentID
		value = c.conn.client.Collection(c.tablename).Doc(value.(string))
	}
	c.query = c.query.Where(path, op, value)
	return c
}

// OrderBy set query order by asc
//
//	list, err = table.Query().OrderBy("Name").Execute(ctx)
//	So((list[0].(*Sample)).Name, ShouldEqual, "sample1")
//
func (c *QueryFirestore) OrderBy(path string) Query {
	c.query = c.query.OrderBy(path, firestore.Asc)
	return c
}

// OrderByDesc set query order by desc
//
//	list, err = table.Query().OrderByDesc("Name").Limit(1).Execute(ctx)
//	So((list[0].(*Sample)).Name, ShouldEqual, "sample2")
//
func (c *QueryFirestore) OrderByDesc(path string) Query {
	c.query = c.query.OrderBy(path, firestore.Desc)
	return c
}

// Limit set query limit
//
//	list, err = table.Query().OrderBy("Name").Limit(1).Execute(ctx)
//	So(len(list), ShouldEqual, 1)
//
func (c *QueryFirestore) Limit(n int) Query {
	c.limit = n
	c.query = c.query.Limit(n)
	return c
}

// StartAt implement Paginate on firestore, please be aware not use index but fieldValue to do the trick, see sample
//
//	list, err = table.Query().OrderBy("Name").StartAt("irvine city").Execute(ctx)
//	So(len(list), ShouldEqual, 1)
//	So((list[0].(*Sample)).Name, ShouldEqual, "irvine city")
//
func (c *QueryFirestore) StartAt(docSnapshotOrFieldValues ...interface{}) Query {
	c.query = c.query.StartAt(docSnapshotOrFieldValues...)
	return c
}

// StartAfter implement Paginate on firestore, please be aware not use index but fieldValue to do the trick, see sample
//
//	list, err = table.Query().OrderBy("Name").StartAfter("santa ana city").Execute(ctx)
//	So((list[0].(*Sample)).Name, ShouldEqual, "irvine city")
//
func (c *QueryFirestore) StartAfter(docSnapshotOrFieldValues ...interface{}) Query {
	c.query = c.query.StartAfter(docSnapshotOrFieldValues...)
	return c
}

// EndAt implement Paginate on firestore, please be aware not use index but fieldValue to do the trick, see sample
//
//	list, err = table.Query().OrderBy("Name").EndAt("irvine city").Execute(ctx)
//	So((list[0].(*Sample)).Name, ShouldEqual, "irvine city")
//
func (c *QueryFirestore) EndAt(docSnapshotOrFieldValues ...interface{}) Query {
	c.query = c.query.EndAt(docSnapshotOrFieldValues...)
	return c
}

// EndBefore implement Paginate on firestore, please be aware not use index but fieldValue to do the trick, see sample
//
//	list, err = table.Query().OrderBy("Name").EndBefore("irvine city").Execute(ctx)
//	So((list[0].(*Sample)).Name, ShouldEqual, "santa ana city")
//
func (c *QueryFirestore) EndBefore(docSnapshotOrFieldValues ...interface{}) Query {
	c.query = c.query.EndBefore(docSnapshotOrFieldValues...)
	return c
}

// GetFirstObject execute query return first object in result
//
//	obj, err := db.Select(ctx, GreetFactory).OrderBy("From").Limit(1).StartAt("b city").GetFirstObject(ctx)
//	greet := obj.(*Greet)
//	So(greet.From, ShouldEqual, "b city")
//
func (c *QueryFirestore) GetFirstObject(ctx context.Context) (Object, error) {
	list, err := c.Limit(1).Execute(ctx)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, nil
	}
	return list[0], nil
}

// GetFirstID execute query return first object id in result
//
//	id, err := db.Select(ctx, GreetFactory).OrderBy("From").Limit(1).StartAt("b city").GetFirstID(ctx)
//	So(id, ShouldEqual, "city1")
//
func (c *QueryFirestore) GetFirstID(ctx context.Context) (string, error) {
	list, err := c.Limit(1).Execute(ctx)
	if err != nil {
		return "", err
	}
	if len(list) == 0 {
		return "", nil
	}
	return list[0].GetID(), nil
}

// Execute query with default limit to 10 object, use Limit() to override default limit, return nil if anything wrong
//
//	list, err = table.Query().OrderByDesc("Name").Limit(1).Execute(ctx)
//	So(len(list), ShouldEqual, 1)
//	So((list[0].(*Sample)).Name, ShouldEqual, "sample2")
//
func (c *QueryFirestore) Execute(ctx context.Context) ([]Object, error) {
	if c.limit == 0 {
		c.Limit(limitQueryDefault)
	}
	result := []Object{}
	var iter *firestore.DocumentIterator
	if c.conn.tx != nil {
		iter = c.conn.tx.Documents(c.query)
	} else {
		iter = c.query.Documents(ctx)
	}
	defer iter.Stop()

	for {
		snapshot, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		object := c.factory()
		if object == nil {
			return nil, errors.New("failed to create object from factory")
		}

		err = snapshot.DataTo(object)
		if err != nil {
			return nil, err
		}
		object.SetRef(snapshot.Ref)
		object.SetID(snapshot.Ref.ID)
		result = append(result, object)
	}
	return result, nil
}

// GetIDs query with default limit to 20 object, use Limit() to override default limit, return nil if anything wrong
//
//	idList, err := db.Select(ctx, GreetFactory).OrderBy("From").Limit(1).StartAt("b city").GetIDs(ctx)
//	So(len(idList), ShouldEqual, 1)
//
func (c *QueryFirestore) GetIDs(ctx context.Context) ([]string, error) {
	if c.limit == 0 {
		c.Limit(limitQueryDefault)
	}
	var result []string

	var iter *firestore.DocumentIterator
	if c.conn.tx != nil {
		iter = c.conn.tx.Documents(c.query)
	} else {
		iter = c.query.Documents(ctx)
	}
	defer iter.Stop()

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

// Count execute query and return max 10 count
//
//	count, err := table.Query().Where("Name", "==", "sample1").Count(ctx)
//	So(count, ShouldEqual, 1)
//
func (c *QueryFirestore) Count(ctx context.Context) (int, error) {
	if c.limit == 0 {
		c.Limit(limitQueryDefault)
	}
	var iter *firestore.DocumentIterator
	if c.conn.tx != nil {
		iter = c.conn.tx.Documents(c.query)
	} else {
		iter = c.query.Documents(ctx)
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

// IsEmpty execute query and return false if object exist
//
//	isEmpty, err := table.Query().Where("Name", "==", "sample1").IsEmpty(ctx)
//	So(isEmpty, ShouldBeFalse)
//
func (c *QueryFirestore) IsEmpty(ctx context.Context) (bool, error) {
	c.Limit(1)
	var iter *firestore.DocumentIterator
	if c.conn.tx != nil {
		iter = c.conn.tx.Documents(c.query)
	} else {
		iter = c.query.Documents(ctx)
	}
	defer iter.Stop()

	_, err := iter.Next()
	if err == iterator.Done {
		return true, nil
	}
	if err != nil {
		return false, err
	}
	return false, nil
}

// IsExist execute query and return true if object exist
//
//	isExist, err := table.Query().Where("Name", "==", "sample1").IsExist(ctx)
//	So(isExist, ShouldBeFalse)
//
func (c *QueryFirestore) IsExist(ctx context.Context) (bool, error) {
	empty, err := c.IsEmpty(ctx)
	return !empty, err
}

// Clear keep delete all object in a query until ctx timeout or all object deleted. it delete 500 documents at a time, return total delete count
//
func (c *QueryFirestore) Clear(ctx context.Context) (int, error) {
	deleteCount := 0
	if c.conn.tx != nil {
		c.query.Limit(limitTransactionClear)
		iter := c.conn.tx.Documents(c.query)
		defer iter.Stop()
		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return deleteCount, errors.Wrap(err, "failed to iterator documents")
			}
			c.conn.tx.Delete(doc.Ref)
			deleteCount++
		}
		return deleteCount, nil
	}
	for {
		// keep delete until ctx timeout or all object deleted
		if ctx.Err() != nil {
			return deleteCount, ctx.Err()
		}
		numDeleted := 0
		c.query.Limit(limitClear)
		iter := c.query.Documents(ctx)
		defer iter.Stop()
		// Iterate through the documents, adding a delete operation for each one to a WriteBatch.
		batch := c.conn.client.Batch()
		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return deleteCount, errors.Wrap(err, "failed to iterator documents")
			}
			batch.Delete(doc.Ref)
			deleteCount++
			numDeleted++
		}
		if numDeleted > 0 {
			_, err := batch.Commit(ctx)
			if err != nil {
				return deleteCount, errors.Wrap(err, "failed to commit batch")
			}
		}
		if numDeleted < limitClear {
			break
		}
	}
	return deleteCount, nil
}
