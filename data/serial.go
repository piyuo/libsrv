package data

import (
	"context"

	"cloud.google.com/go/firestore"
	util "github.com/piyuo/libsrv/util"
	"github.com/pkg/errors"
)

// Serial is collections of serial in document database
//
type Serial struct {
	conn      ConnectionRef
	TableName string
}

// Code16 encode uint16 number into string, please be aware serial can only generate one number per second
//
//	code, err := serial.Code16(ctx, "sample-id")
//	So(code, ShouldBeEmpty)
//
func (s *Serial) Code16(ctx context.Context, name string) (string, error) {
	if ctx.Err() != nil {
		return "", ctx.Err()
	}

	number, err := s.Number(ctx, name)
	if err != nil {
		return "", err
	}
	return util.SerialID16(uint16(number)), nil
}

// Code32 encode uint32 number into string, please be aware serial can only generate one number per second
//
//	code, err := serial.Code32(ctx, "sample-id")
//	So(code, ShouldBeEmpty)
//
func (s *Serial) Code32(ctx context.Context, name string) (string, error) {
	if ctx.Err() != nil {
		return "", ctx.Err()
	}
	number, err := s.Number(ctx, name)
	if err != nil {
		return "", err
	}
	return util.SerialID32(uint32(number)), nil
}

// Code64 encode int64 serial number to string, please be aware serial can only generate one number per second
//
//	code, err := serial.Code64(ctx, "sample-id")
//	So(code, ShouldBeEmpty)
//
func (s *Serial) Code64(ctx context.Context, name string) (string, error) {
	if ctx.Err() != nil {
		return "", ctx.Err()
	}
	number, err := s.Number(ctx, name)
	if err != nil {
		return "", err
	}
	return util.SerialID64(uint64(number)), nil
}

// Number create unique serial number, please be aware serial can only generate one number per second and use with transation to ensure unique
//
//	num, err := serial.Number(ctx, "sample-id")
//	So(num, ShouldEqual, 1)
//
func (s *Serial) Number(ctx context.Context, name string) (int64, error) {
	if ctx.Err() != nil {
		return 0, ctx.Err()
	}
	if s.TableName == "" {
		return 0, errors.New("table name can not be empty: " + name)
	}

	fConn := s.conn.(*ConnectionFirestore)
	if fConn.tx != nil {
		return s.getNumberInTx(ctx, fConn.tx, name)
	}

	var id int64
	var err error
	err = fConn.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		id, err = s.getNumberInTx(ctx, tx, name)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return 0, errors.Wrap(err, "failed to commit serial transaction: "+fConn.errorID(s.TableName, name))
	}
	return id, nil

}

// getNumberInTx generate number in transaction
//
//	num, err := s.getNumberInTx(ctx, "sample-id")
//	So(num, ShouldEqual, 1)
//
func (s *Serial) getNumberInTx(ctx context.Context, tx *firestore.Transaction, name string) (int64, error) {
	fConn := s.conn.(*ConnectionFirestore)
	docRef := fConn.getDocRef(s.TableName, name)

	snapshot, err := tx.Get(docRef)
	if snapshot != nil && !snapshot.Exists() {
		err := tx.Set(docRef, map[string]interface{}{
			"S": 1,
		}, firestore.MergeAll)
		if err != nil {
			return 0, errors.Wrap(err, "failed to init serial: "+fConn.errorID(s.TableName, name))
		}
		return 1, nil
	}

	if err != nil {
		return 0, errors.Wrap(err, "failed to get serial: "+fConn.errorID(s.TableName, name))
	}
	idRef, err := snapshot.DataAt("S")
	if err != nil {
		return 0, errors.Wrap(err, "failed to get value from serial: "+fConn.errorID(s.TableName, name))
	}
	id := idRef.(int64)
	err = tx.Update(docRef, []firestore.Update{
		{Path: "S", Value: firestore.Increment(1)},
	})
	if err != nil {
		return 0, errors.Wrap(err, "failed to increment serial: "+fConn.errorID(s.TableName, name))
	}
	return id + 1, nil
}

// Delete serial
//
//	err = serial.Delete(ctx, "sample-id")
//
func (s *Serial) Delete(ctx context.Context, name string) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	if s.TableName == "" {
		return errors.New("serial table name can not be empty: " + name)
	}

	return s.conn.Delete(ctx, s.TableName, name)
}
