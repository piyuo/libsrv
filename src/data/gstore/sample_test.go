package gstore

import (
	"context"

	"github.com/piyuo/libsrv/src/data"
)

type SampleDB interface {
	data.DB
	SampleTable() *data.Table
	Counters() *SampleCounters
	Serials() *SampleSerials
	Coders() *SampleCoders
}

// global connection
//
type SampleGlobalDB struct {
	data.BaseDB
}

func NewSampleGlobalDB(ctx context.Context) (*SampleGlobalDB, error) {
	conn, err := ConnectGlobalFirestore(ctx)
	if err != nil {
		return nil, err
	}
	db := &SampleGlobalDB{
		BaseDB: data.BaseDB{Connection: conn},
	}
	return db, nil
}

func (db *SampleGlobalDB) SampleTable() *data.Table {
	return &data.Table{
		Connection: db.Connection,
		TableName:  "Sample",
		Factory: func() data.Object {
			return &Sample{}
		},
	}
}

func (db *SampleGlobalDB) Counters() *SampleCounters {
	return &SampleCounters{
		Counters: data.Counters{
			Connection: db.Connection,
			TableName:  "Count",
		},
	}
}

func (db *SampleGlobalDB) Serials() *SampleSerials {
	return &SampleSerials{
		Serials: data.Serials{
			Connection: db.Connection,
			TableName:  "Serial",
		},
	}
}

func (db *SampleGlobalDB) Coders() *SampleCoders {
	return &SampleCoders{
		Coders: data.Coders{
			Connection: db.Connection,
			TableName:  "Code",
		},
	}
}

// regional connection
//
type SampleRegionalDB struct {
	data.BaseDB
}

func NewSampleRegionalDB(ctx context.Context) (*SampleRegionalDB, error) {
	conn, err := ConnectRegionalFirestore(ctx)
	if err != nil {
		return nil, err
	}
	db := &SampleRegionalDB{
		BaseDB: data.BaseDB{Connection: conn},
	}
	return db, nil
}

func (db *SampleRegionalDB) SampleTable() *data.Table {
	return &data.Table{
		Connection: db.Connection,
		TableName:  "Sample",
		Factory: func() data.Object {
			return &Sample{}
		},
	}
}

func (db *SampleRegionalDB) Counters() *SampleCounters {
	return &SampleCounters{
		Counters: data.Counters{
			Connection: db.Connection,
			TableName:  "Count",
		},
	}
}

func (db *SampleRegionalDB) Serials() *SampleSerials {
	return &SampleSerials{
		Serials: data.Serials{
			Connection: db.Connection,
			TableName:  "Serial",
		},
	}
}

func (db *SampleRegionalDB) Coders() *SampleCoders {
	return &SampleCoders{
		Coders: data.Coders{
			Connection: db.Connection,
			TableName:  "Code",
		},
	}
}

type PlainObject struct {
	ID   string
	Name string
}

// Sample
//
type Sample struct {
	data.DomainObject
	Name    string
	Value   int
	Map     map[string]string
	Array   []string
	Numbers []int
	Obj     *PlainObject
}

func (c *Sample) Factory() data.Object {
	return &Sample{}
}

func (c *Sample) TableName() string {
	return "Sample"
}

// SampleCoders  represent collection of code
//
type SampleCoders struct {
	data.Coders `firestore:"-"`
}

// SampleCoder return sample code
//
func (ss *SampleCoders) SampleCoder() data.Coder {
	return ss.Connection.CreateCoder(ss.TableName, "SampleCode", 10)
}

// SampleCoder100 return sample code with 100 shards
//
func (ss *SampleCoders) SampleCoder1000() data.Coder {
	return ss.Connection.CreateCoder(ss.TableName, "SampleCode", 1000)
}

// SampleSerials  represent collection of serial
//
type SampleSerials struct {
	data.Serials `firestore:"-"`
}

func (ss *SampleSerials) SampleSerial() data.Serial {
	return ss.Connection.CreateSerial(ss.TableName, "SampleSerial")
}

// SampleCounters represent collection of counter
//
type SampleCounters struct {
	data.Counters `firestore:"-"`
}

// SampleCounter return sample counter
//
func (scs *SampleCounters) SampleCounter() data.Counter {
	return scs.Connection.CreateCounter(scs.TableName, "SampleCount", 3, data.DateHierarchyNone)
}

// SampleCounter100 return sample counter with 100 shards
//
func (scs *SampleCounters) SampleCounter1000() data.Counter {
	return scs.Connection.CreateCounter(scs.TableName, "SampleCount", 1000, data.DateHierarchyNone)
}
