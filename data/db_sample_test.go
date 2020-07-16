package data

import (
	"context"
)

type SampleDB interface {
	DBRef
	SampleTable() *Table
	Counters() *SampleCounters
	Serials() *SampleSerials
	Coders() *SampleCoders
}

// global connection
//
type SampleGlobalDB struct {
	DB
}

func NewSampleGlobalDB(ctx context.Context) (*SampleGlobalDB, error) {
	conn, err := FirestoreGlobalConnection(ctx)
	if err != nil {
		return nil, err
	}
	db := &SampleGlobalDB{
		DB: DB{Connection: conn},
	}
	return db, nil
}

func (db *SampleGlobalDB) SampleTable() *Table {
	return &Table{
		Connection: db.Connection,
		TableName:  "sample",
		Factory: func() ObjectRef {
			return &Sample{}
		},
	}
}

func (db *SampleGlobalDB) Counters() *SampleCounters {
	return &SampleCounters{
		Counters: Counters{
			Connection: db.Connection,
			TableName:  "sample-count",
		},
	}
}

func (db *SampleGlobalDB) Serials() *SampleSerials {
	return &SampleSerials{
		Serials: Serials{
			Connection: db.Connection,
			TableName:  "sample-serial",
		},
	}
}

func (db *SampleGlobalDB) Coders() *SampleCoders {
	return &SampleCoders{
		Coders: Coders{
			Connection: db.Connection,
			TableName:  "sample-code",
		},
	}
}

// regional connection
//
type SampleRegionalDB struct {
	DB
}

func NewSampleRegionalDB(ctx context.Context, databaseName string) (*SampleRegionalDB, error) {
	conn, err := FirestoreRegionalConnection(ctx, databaseName)
	if err != nil {
		return nil, err
	}
	db := &SampleRegionalDB{
		DB: DB{Connection: conn},
	}
	return db, nil
}

func (db *SampleRegionalDB) SampleTable() *Table {
	return &Table{
		Connection: db.Connection,
		TableName:  "sample",
		Factory: func() ObjectRef {
			return &Sample{}
		},
	}
}

func (db *SampleRegionalDB) Counters() *SampleCounters {
	return &SampleCounters{
		Counters: Counters{
			Connection: db.Connection,
			TableName:  "sample-count",
		},
	}
}

func (db *SampleRegionalDB) Serials() *SampleSerials {
	return &SampleSerials{
		Serials: Serials{
			Connection: db.Connection,
			TableName:  "sample-serial",
		},
	}
}

func (db *SampleRegionalDB) Coders() *SampleCoders {
	return &SampleCoders{
		Coders: Coders{
			Connection: db.Connection,
			TableName:  "sample-code",
		},
	}
}

// Sample
//
type Sample struct {
	Object `firestore:"-"`
	Name   string
	Value  int
}

// SampleCoders  represent collection of code
//
type SampleCoders struct {
	Coders `firestore:"-"`
}

// SampleCoder return sample code
//
func (ss *SampleCoders) SampleCoder() CoderRef {
	return ss.Coder("sample-code", 10)
}

// DeleteSampleSerial delete sample serial
//
func (ss *SampleCoders) DeleteSampleCode(ctx context.Context) error {
	return ss.Delete(ctx, "sample-code")
}

// SampleSerials  represent collection of serial
//
type SampleSerials struct {
	Serials `firestore:"-"`
}

func (ss *SampleSerials) SampleSerial() SerialRef {
	return ss.Serial("sample-no")
}

// DeleteSampleSerial delete sample serial
//
func (ss *SampleSerials) DeleteSampleSerial(ctx context.Context) error {
	return ss.Delete(ctx, "sample-no")
}

// SampleCounters represent collection of counter
//
type SampleCounters struct {
	Counters `firestore:"-"`
}

// SampleTotal return sample total count
//
func (scs *SampleCounters) SampleCounter() CounterRef {
	return scs.Counter("sample-counter", 4)
}

// DeleteSampleCounter delete sample counter
//
func (scs *SampleCounters) DeleteSampleCounter(ctx context.Context) error {
	return scs.Delete(ctx, "sample-counter")
}
