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

	numberCallRX bool

	numberCanCreate bool

	numberCanIncrement bool
}

// NumberRX return sequence number, number is unique and serial, please be aware serial can only generate one sequence per second, use it with high frequency will cause error and  must used it in transaction with NumberWX()
//
//	err = db.Transaction(ctx, func(ctx context.Context) error {
//		num, err:= serial.NumberRX()
//		So(err, ShouldBeNil)
//		So(num, ShouldEqual,1)
//		err := serial.NumberWX()
//	})
//
func (c *SerialFirestore) NumberRX() (int64, error) {
	if c.conn.tx == nil {
		return 0, errors.New("this function must run in transaction")
	}

	c.numberCallRX = true
	c.numberCanCreate = false
	c.numberCanIncrement = false

	docRef, _ := c.getRef()
	snapshot, err := c.conn.tx.Get(docRef)
	if snapshot != nil && !snapshot.Exists() {
		c.numberCanCreate = true
		return 1, nil
	}

	if err != nil {
		return 0, errors.Wrap(err, "failed to get serial: "+c.errorID())
	}

	idRef, err := snapshot.DataAt("N")
	if err != nil {
		return 0, errors.Wrap(err, "failed to get value from serial: "+c.errorID())
	}
	c.numberCanIncrement = true
	id := idRef.(int64)
	return id + 1, nil
}

// NumberWX commit NumberRX
//
//	err = db.Transaction(ctx, func(ctx context.Context) error {
//		num, err:= serial.NumberRX()
//		So(err, ShouldBeNil)
//		So(num, ShouldEqual,1)
//		err := serial.NumberWX()
//	})
//
func (c *SerialFirestore) NumberWX() error {
	if c.conn.tx == nil {
		return errors.New("This function must run in transaction")
	}
	if c.numberCallRX == false {
		return errors.New("WX() function need call NumberRX() first")
	}

	docRef, _ := c.getRef()
	if c.numberCanCreate {
		err := c.conn.tx.Set(docRef, map[string]interface{}{"N": 1}, firestore.MergeAll)
		if err != nil {
			return errors.Wrap(err, "failed to create serial: "+c.errorID())
		}
	}

	if c.numberCanIncrement {
		err := c.conn.tx.Update(docRef, []firestore.Update{
			{Path: "N", Value: firestore.Increment(1)},
		})
		if err != nil {
			return errors.Wrap(err, "failed to increment serial: "+c.errorID())
		}
	}
	c.numberCallRX = false
	c.numberCanCreate = false
	c.numberCanIncrement = false
	return nil
}

// Reset reset sequence number
//
//	err = db.Transaction(ctx, func(ctx context.Context) error {
//		err:= serial.Reset(ctx)
//	})
//
func (c *SerialFirestore) Reset(ctx context.Context) error {
	if err := c.assert(ctx); err != nil {
		return err
	}
	return c.deleteDoc(ctx)
}