package data

import (
	"context"
)

type SampleDB interface {
	DBRef
	SampleTable() *Table
	Counters() *SampleCounters
	Serials() *SampleSerials
	Codes() *SampleCodes
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
	table := &Table{
		Connection: db.Connection,
		TableName:  "sample",
		Factory: func() ObjectRef {
			return &Sample{}
		},
	}
	return table
}

func (db *SampleGlobalDB) Counters() *SampleCounters {
	counters := &SampleCounters{
		Counters: Counters{
			Connection: db.Connection,
			TableName:  "sample-counter",
		},
	}
	return counters
}

func (db *SampleGlobalDB) Serials() *SampleSerials {
	serials := &SampleSerials{
		Serials: Serials{
			Connection: db.Connection,
			TableName:  "sample-serial",
		},
	}
	return serials
}

func (db *SampleGlobalDB) Codes() *SampleCodes {
	codes := &SampleCodes{
		Codes: Codes{
			Connection: db.Connection,
			TableName:  "sample-code",
		},
	}
	return codes
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
	table := &Table{
		Connection: db.Connection,
		TableName:  "sample",
		Factory: func() ObjectRef {
			return &Sample{}
		},
	}
	return table
}

func (db *SampleRegionalDB) Counters() *SampleCounters {
	counters := &SampleCounters{
		Counters: Counters{
			Connection: db.Connection,
			TableName:  "sample-counter",
		},
	}
	return counters
}

func (db *SampleRegionalDB) Serials() *SampleSerials {
	serials := &SampleSerials{
		Serials: Serials{
			Connection: db.Connection,
			TableName:  "sample-serial",
		},
	}
	return serials
}

func (db *SampleRegionalDB) Codes() *SampleCodes {
	codes := &SampleCodes{
		Codes: Codes{
			Connection: db.Connection,
			TableName:  "sample-code",
		},
	}
	return codes
}

// Sample
//
type Sample struct {
	Object `firestore:"-"`
	Name   string
	Value  int
}

// SampleCodes  represent collection of code
//
type SampleCodes struct {
	Codes `firestore:"-"`
}

func (ss *SampleCodes) SampleCode() CodeRef {
	return ss.Code("sample-code", 10)
}

// DeleteSampleSerial delete sample serial
//
func (ss *SampleCodes) DeleteSampleCode(ctx context.Context) error {
	return ss.Delete(ctx, "sample-code")
}

// SampleSerials  represent collection of serial
//
type SampleSerials struct {
	Serials `firestore:"-"`
}

func (ss *SampleSerials) SampleSerial() SerialRef {
	return ss.Serial("sample-serial")
}

// DeleteSampleSerial delete sample serial
//
func (ss *SampleSerials) DeleteSampleSerial(ctx context.Context) error {
	return ss.Delete(ctx, "sample-serial")
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
	return scs.Delete(ctx, "sample-total")
}
