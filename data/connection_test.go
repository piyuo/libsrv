package data

import (
	"context"

	. "github.com/smartystreets/goconvey/convey"
)

// global connection
//
type SampleGlobalDB struct {
	DocDB
}

func NewSampleGlobalDB(ctx context.Context) (*SampleGlobalDB, error) {
	conn, err := FirestoreGlobalConnection(ctx, "")
	if err != nil {
		return nil, err
	}
	db := &SampleGlobalDB{}
	db.SetConnection(conn)
	return db, nil
}

func (db *SampleGlobalDB) SampleTable() Table {
	table := &DocTable{}
	table.SetConnection(db.Connection())
	table.SetTableName("sample")
	factory := func() Object {
		return &Sample{}
	}
	table.SetFactory(factory)
	return table
}

func (db *SampleGlobalDB) Counter() *SampleCounters {
	counter := &SampleCounters{}
	counter.SetConnection(db.Connection())
	counter.SetTableName("sample-counter")
	return counter
}

func (db *SampleGlobalDB) Serial() *SampleSerial {
	serial := &SampleSerial{}
	serial.SetConnection(db.Connection())
	serial.SetTableName("sample-serial")
	return serial
}

// regional connection
//
type SampleRegionalDB struct {
	DocDB
}

func NewSampleRegionalDB(ctx context.Context, databaseName string) (*SampleRegionalDB, error) {
	conn, err := FirestoreRegionalConnection(ctx, databaseName)
	if err != nil {
		return nil, err
	}
	db := &SampleRegionalDB{}
	db.SetConnection(conn)
	return db, nil
}

func (db *SampleRegionalDB) SampleTable() Table {
	table := &DocTable{}
	table.SetConnection(db.Connection())
	table.SetTableName("sample")
	factory := func() Object {
		return &Sample{}
	}
	table.SetFactory(factory)
	return table
}

func (db *SampleRegionalDB) Counter() *SampleCounters {
	counter := &SampleCounters{}
	counter.SetConnection(db.Connection())
	counter.SetTableName("sample-counter")
	return counter
}

func (db *SampleRegionalDB) Serial() *SampleSerial {
	serial := &SampleSerial{}
	serial.SetConnection(db.Connection())
	serial.SetTableName("sample-serial")
	return serial
}

// Sample
//
type Sample struct {
	DocObject `firestore:"-"`
	Name      string
	Value     int
}

// SampleSerial
//
type SampleSerial struct {
	DocSerial `firestore:"-"`
}

func (ss *SampleSerial) SampleID(ctx context.Context) (string, error) {
	return ss.Code(ctx, "sample-id")
}

// SampleCounter represent collection of counter
//
type SampleCounters struct {
	DocCounters `firestore:"-"`
}

// SampleTotal return sample total count
//
func (scs *SampleCounters) SampleTotal(ctx context.Context) (Counter, error) {
	return scs.Counter(ctx, "sample-total", 4)
}

func firestoreBeginTest() (*SampleGlobalDB, *SampleRegionalDB, Table, Table) {
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

func firestoreEndTest(dbG *SampleGlobalDB, dbR *SampleRegionalDB, samplesG Table, samplesR Table) {
	ctx := context.Background()
	err := samplesG.Clear(ctx)
	So(err, ShouldBeNil)
	err = samplesR.Clear(ctx)
	So(err, ShouldBeNil)
	err = dbR.DeleteNamespace(ctx)
	So(err, ShouldBeNil)
}
