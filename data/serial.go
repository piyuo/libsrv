package data

import (
	"context"

	util "github.com/piyuo/libsrv/util"
	"github.com/pkg/errors"
)

// Serial is collections of serial in document database
//
type Serial struct {
	Connection ConnectionRef
	TableName  string
}

// Number table save all serial
//
type Number struct {
	Object `firestore:"-"`
	S      uint32
}

// newNumber create number object
//
func newNumber() ObjectRef {
	return &Number{}
}

// Code encode serial number to string, please be aware serial can only generate one number per second and use with transation to ensure unique
//
//	id,err = db.Serial(ctx,"", "myID")
//
func (s *Serial) Code(ctx context.Context, name string) (string, error) {
	number, err := s.Number(ctx, name)
	if err != nil {
		return "", err
	}
	return util.SerialID32(number), nil
}

// Number create unique serial number, please be aware serial can only generate one number per second and use with transation to ensure unique
//
//	id,err = db.Serial(ctx,"", "myID")
//
func (s *Serial) Number(ctx context.Context, name string) (uint32, error) {
	if ctx.Err() != nil {
		return 0, ctx.Err()
	}
	if s.TableName == "" {
		return 0, errors.New("table name can not be empty: " + name)
	}

	num, err := s.Connection.Get(ctx, s.TableName, name, newNumber)
	if err != nil {
		return 0, errors.Wrap(err, "failed to get serial: "+name)
	}

	var number *Number
	if num == nil {
		number = &Number{
			S: 1,
		}
		number.ID = name
	} else {
		number = num.(*Number)
		number.S++
	}

	err = s.Connection.Set(ctx, s.TableName, number)
	if err != nil {
		return 0, errors.Wrap(err, "failed to set serial: "+name)
	}
	return number.S, nil
}

// Delete serial
//
//	counter,err = db.GetCounter(ctx, "myCounter")
//
func (s *Serial) Delete(ctx context.Context, name string) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	if s.TableName == "" {
		return errors.New("serial table name can not be empty: " + name)
	}

	return s.Connection.Delete(ctx, s.TableName, name)
}
