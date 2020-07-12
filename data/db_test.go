package data

import (
	"context"

	. "github.com/smartystreets/goconvey/convey"
)

type SampleDB interface {
	DBRef
	SampleTable() *Table
	Counters() *SampleCounters
	Serial() *SampleSerial
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
	counter := &SampleCounters{
		Counters: Counters{
			Connection: db.Connection,
			TableName:  "sample-counter",
		},
	}
	return counter
}

func (db *SampleGlobalDB) Serial() *SampleSerial {
	serial := &SampleSerial{
		Serial: Serial{
			Connection: db.Connection,
			TableName:  "sample-serial",
		},
	}
	return serial
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
	counter := &SampleCounters{
		Counters: Counters{
			Connection: db.Connection,
			TableName:  "sample-counter",
		},
	}
	return counter
}

func (db *SampleRegionalDB) Serial() *SampleSerial {
	serial := &SampleSerial{
		Serial: Serial{
			Connection: db.Connection,
			TableName:  "sample-serial",
		},
	}
	return serial
}

// Sample
//
type Sample struct {
	Object `firestore:"-"`
	Name   string
	Value  int
}

// SampleSerial
//
type SampleSerial struct {
	Serial `firestore:"-"`
}

func (ss *SampleSerial) SampleID(ctx context.Context) (string, error) {
	return ss.Code(ctx, "sample-id")
}

// SampleCounter represent collection of counter
//
type SampleCounters struct {
	Counters `firestore:"-"`
}

// SampleTotal return sample total count
//
func (scs *SampleCounters) SampleTotal(ctx context.Context) CounterRef {
	return scs.Counter("sample-total", 4)
}

// DeleteSampleTotal return sample total count
//
func (scs *SampleCounters) DeleteSampleTotal(ctx context.Context) error {
	return scs.Delete(ctx, "sample-total")
}

func firestoreBeginTest() (*SampleGlobalDB, *SampleRegionalDB, *Table, *Table) {
	ctx := context.Background()
	dbG, err := NewSampleGlobalDB(ctx)
	So(err, ShouldBeNil)
	samplesG := dbG.SampleTable()
	So(samplesG, ShouldNotBeNil)
	err = samplesG.Clear(ctx)
	So(err, ShouldBeNil)

	dbR, err := NewSampleRegionalDB(ctx, "sample-namespace")
	samplesR := dbR.SampleTable()
	So(samplesR, ShouldNotBeNil)
	So(err, ShouldBeNil)
	samplesR.Clear(ctx)
	err = dbR.DeleteNamespace(ctx)
	So(err, ShouldBeNil)
	err = dbR.CreateNamespace(ctx)
	So(err, ShouldBeNil)
	return dbG, dbR, samplesG, samplesR
}

func firestoreEndTest(dbG *SampleGlobalDB, dbR *SampleRegionalDB, samplesG *Table, samplesR *Table) {
	ctx := context.Background()
	err := samplesG.Clear(ctx)
	So(err, ShouldBeNil)
	err = samplesR.Clear(ctx)
	So(err, ShouldBeNil)
	err = dbR.DeleteNamespace(ctx)
	So(err, ShouldBeNil)
}
