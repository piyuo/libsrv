package data

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/pkg/errors"
)

// SerialFirestore generial serial from firestore
//
type SerialFirestore struct {
	SerialRef `firestore:"-"`

	ShardsFirestore `firestore:"-"`
}

// Number return code number, number is unique but not serial
//
//	code, err := serial.Get(ctx)
//	So(code, ShouldBeEmpty)
//
func (c *SerialFirestore) Number(ctx context.Context) (int64, error) {
	if err := c.assert(ctx); err != nil {
		return 0, err
	}

	if c.conn.tx != nil {
		return c.getTx(ctx, c.conn.tx)
	}

	var id int64
	var err error
	err = c.conn.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		id, err = c.getTx(ctx, tx)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return 0, errors.Wrap(err, "failed to commit serial transaction: "+c.errorID())
	}
	return id, nil
}

// getTx generate code number in transaction, return number is not serial
//
//	num, err := s.getTx(ctx)
//	So(num, ShouldEqual, 1)
//
func (c *SerialFirestore) getTx(ctx context.Context, tx *firestore.Transaction) (int64, error) {
	docRef, _ := c.getRef()
	snapshot, err := tx.Get(docRef)
	if snapshot != nil && !snapshot.Exists() {
		err := tx.Set(docRef, map[string]interface{}{"N": 1}, firestore.MergeAll)
		if err != nil {
			return 0, errors.Wrap(err, "failed to create serial: "+c.errorID())
		}
		return 1, nil
	}

	if err != nil {
		return 0, errors.Wrap(err, "failed to get serial: "+c.errorID())
	}
	idRef, err := snapshot.DataAt("N")
	if err != nil {
		return 0, errors.Wrap(err, "failed to get value from serial: "+c.errorID())
	}
	id := idRef.(int64)
	err = tx.Update(docRef, []firestore.Update{
		{Path: "S", Value: firestore.Increment(1)},
	})
	if err != nil {
		return 0, errors.Wrap(err, "failed to increment serial: "+c.errorID())
	}
	return id + 1, nil
}
