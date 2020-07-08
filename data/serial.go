package data

import (
	"context"

	util "github.com/piyuo/libsrv/util"
	"github.com/pkg/errors"
)

// Serial is collections of serial in document database
//
type Serial struct {
	conn      Connection
	tablename string
}

// Number table save all serial
//
type Number struct {
	DocObject `firestore:"-"`
	S         uint32
}

// newNumber create number object
//
func newNumber() Object {
	return &Number{}
}

// SetConnection set connection for table
//
//	table.SetConnection(conn)
//
func (s *Serial) SetConnection(conn Connection) {
	s.conn = conn
}

// SetTableName set table name
//
//	table.SetTableName("sample")
//
func (s *Serial) SetTableName(tablename string) {
	s.tablename = tablename
}

// TableName return table name
//
//	table.TableName()
//
func (s *Serial) TableName() string {
	return s.tablename
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
	if s.tablename == "" {
		return 0, errors.New("table name can not be empty: " + name)
	}

	num, err := s.conn.Get(ctx, s.tablename, name, newNumber)
	if err != nil {
		return 0, errors.Wrap(err, "failed to get serial: "+name)
	}

	var number *Number
	if num == nil {
		number = &Number{
			S: 1,
		}
		number.SetID(name)
	} else {
		number = num.(*Number)
		number.S++
	}

	err = s.conn.Set(ctx, s.tablename, number)
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
	if s.tablename == "" {
		return errors.New("serial table name can not be empty: " + name)
	}

	return s.conn.Delete(ctx, s.tablename, name)
}
