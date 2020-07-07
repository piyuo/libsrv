package data

import (
	"context"

	util "github.com/piyuo/libsrv/util"
	"github.com/pkg/errors"
)

// Serial represent serial table
//
type Serial interface {
	SetConnection(conn Connection)
	SetTableName(tablename string)
	TableName()
	Number(ctx context.Context, name string) (uint32, error)
	Code(ctx context.Context, name string) (string, error)
	Delete(ctx context.Context, name string)
}

// DocSerial is collections of serial in document database
//
type DocSerial struct {
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
func (ds *DocSerial) SetConnection(conn Connection) {
	ds.conn = conn
}

// SetTableName set table name
//
//	table.SetTableName("sample")
//
func (ds *DocSerial) SetTableName(tablename string) {
	ds.tablename = tablename
}

// TableName return table name
//
//	table.TableName()
//
func (ds *DocSerial) TableName() string {
	return ds.tablename
}

// Code encode serial number to string, please be aware serial can only generate one number per second and use with transation to ensure unique
//
//	id,err = db.Serial(ctx,"", "myID")
//
func (ds *DocSerial) Code(ctx context.Context, name string) (string, error) {
	number, err := ds.Number(ctx, name)
	if err != nil {
		return "", err
	}
	return util.SerialID32(number), nil
}

// Number create unique serial number, please be aware serial can only generate one number per second and use with transation to ensure unique
//
//	id,err = db.Serial(ctx,"", "myID")
//
func (ds *DocSerial) Number(ctx context.Context, name string) (uint32, error) {
	if ctx.Err() != nil {
		return 0, ctx.Err()
	}
	if ds.tablename == "" {
		return 0, errors.New("table name can not be empty: " + name)
	}

	num, err := ds.conn.Get(ctx, ds.tablename, name, newNumber)
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

	err = ds.conn.Set(ctx, ds.tablename, number)
	if err != nil {
		return 0, errors.Wrap(err, "failed to set serial: "+name)
	}
	return number.S, nil
}

// Delete serial
//
//	counter,err = db.GetCounter(ctx, "myCounter")
//
func (ds *DocSerial) Delete(ctx context.Context, name string) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	if ds.tablename == "" {
		return errors.New("please implement TableName() on serial: " + name)
	}

	return ds.conn.Delete(ctx, ds.tablename, name)
}
