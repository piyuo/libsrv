package data

import (
	"context"
	"time"
)

type SampleDB interface {
	DB
	SampleTable() *Table
	Counters() *SampleCounters
	Serials() *SampleSerials
	Coders() *SampleCoders
}

// global connection
//
type SampleGlobalDB struct {
	BaseDB
}

func NewSampleGlobalDB(ctx context.Context) (*SampleGlobalDB, error) {
	conn, err := FirestoreGlobalConnection(ctx)
	if err != nil {
		return nil, err
	}
	db := &SampleGlobalDB{
		BaseDB: BaseDB{Connection: conn},
	}
	return db, nil
}

func (db *SampleGlobalDB) SampleTable() *Table {
	return &Table{
		Connection: db.Connection,
		TableName:  "Sample",
		Factory: func() Object {
			return &Sample{}
		},
	}
}

func (db *SampleGlobalDB) Counters() *SampleCounters {
	return &SampleCounters{
		Counters: Counters{
			Connection: db.Connection,
			TableName:  "SampleCount",
		},
	}
}

func (db *SampleGlobalDB) Serials() *SampleSerials {
	return &SampleSerials{
		Serials: Serials{
			Connection: db.Connection,
			TableName:  "SampleSerial",
		},
	}
}

func (db *SampleGlobalDB) Coders() *SampleCoders {
	return &SampleCoders{
		Coders: Coders{
			Connection: db.Connection,
			TableName:  "SampleCode",
		},
	}
}

// regional connection
//
type SampleRegionalDB struct {
	BaseDB
}

func NewSampleRegionalDB(ctx context.Context, databaseName string) (*SampleRegionalDB, error) {
	conn, err := FirestoreRegionalConnection(ctx, databaseName)
	if err != nil {
		return nil, err
	}
	db := &SampleRegionalDB{
		BaseDB: BaseDB{Connection: conn},
	}
	return db, nil
}

func (db *SampleRegionalDB) SampleTable() *Table {
	return &Table{
		Connection: db.Connection,
		TableName:  "Sample",
		Factory: func() Object {
			return &Sample{}
		},
	}
}

func (db *SampleRegionalDB) Counters() *SampleCounters {
	return &SampleCounters{
		Counters: Counters{
			Connection: db.Connection,
			TableName:  "SampleCount",
		},
	}
}

func (db *SampleRegionalDB) Serials() *SampleSerials {
	return &SampleSerials{
		Serials: Serials{
			Connection: db.Connection,
			TableName:  "SampleSerial",
		},
	}
}

func (db *SampleRegionalDB) Coders() *SampleCoders {
	return &SampleCoders{
		Coders: Coders{
			Connection: db.Connection,
			TableName:  "SampleCode",
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
	BaseObject
	Name    string
	Value   int
	Map     map[string]string
	Array   []string
	Numbers []int
	Obj     *PlainObject
}

// SampleCoders  represent collection of code
//
type SampleCoders struct {
	Coders `firestore:"-"`
}

// SampleCoder return sample code
//
func (ss *SampleCoders) SampleCoder() Coder {
	return ss.Coder("sample-code", 10)
}

// SampleCoder100 return sample code with 100 shards
//
func (ss *SampleCoders) SampleCoder1000() Coder {
	return ss.Coder("sample-code", 1000)
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

func (ss *SampleSerials) SampleSerial() Serial {
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

// SampleCounter return sample counter
//
func (scs *SampleCounters) SampleCounter() Counter {
	zone, offset := time.Now().Zone()
	return scs.Counter("Sample", 3, zone, offset)
}

// SampleCounter100 return sample counter with 100 shards
//
func (scs *SampleCounters) SampleCounter1000() Counter {
	zone, offset := time.Now().Zone()
	return scs.Counter("Sample", 1000, zone, offset)
}

// DeleteSampleCounter delete sample counter
//
func (scs *SampleCounters) DeleteSampleCounter(ctx context.Context) error {
	return scs.Delete(ctx, "Sample")
}

func createSampleDB() (*SampleGlobalDB, *SampleRegionalDB) {
	ctx := context.Background()
	dbG, _ := NewSampleGlobalDB(ctx)

	dbR, _ := NewSampleRegionalDB(ctx, "sample-namespace")
	dbR.DeleteNamespace(ctx)
	dbR.CreateNamespace(ctx)
	return dbG, dbR
}

func removeSampleDB(dbG *SampleGlobalDB, dbR *SampleRegionalDB) {
	ctx := context.Background()
	dbR.DeleteNamespace(ctx)
	dbG.Close()
	dbR.Close()
}

func createSampleTable(dbG *SampleGlobalDB, dbR *SampleRegionalDB) (*Table, *Table) {
	g := dbG.SampleTable()
	r := dbR.SampleTable()
	removeSampleTable(g, r)
	return g, r
}

func removeSampleTable(g *Table, r *Table) {
	ctx := context.Background()
	g.Clear(ctx)
	r.Clear(ctx)
}

func createSampleCounters(dbG *SampleGlobalDB, dbR *SampleRegionalDB) (*SampleCounters, *SampleCounters) {
	g := dbG.Counters()
	r := dbR.Counters()
	removeSampleCounters(g, r)
	return g, r
}

func removeSampleCounters(g *SampleCounters, r *SampleCounters) {
	ctx := context.Background()
	g.DeleteSampleCounter(ctx)
	r.DeleteSampleCounter(ctx)
}

func createSampleSerials(dbG *SampleGlobalDB, dbR *SampleRegionalDB) (*SampleSerials, *SampleSerials) {
	g := dbG.Serials()
	r := dbR.Serials()
	removeSampleSerials(g, r)
	return g, r
}

func removeSampleSerials(g *SampleSerials, r *SampleSerials) {
	ctx := context.Background()
	g.DeleteSampleSerial(ctx)
	r.DeleteSampleSerial(ctx)
}

func createSampleCoders(dbG *SampleGlobalDB, dbR *SampleRegionalDB) (*SampleCoders, *SampleCoders) {
	g := dbG.Coders()
	r := dbR.Coders()
	removeSampleCoders(g, r)
	return g, r
}

func removeSampleCoders(g *SampleCoders, r *SampleCoders) {
	ctx := context.Background()
	g.DeleteSampleCode(ctx)
	r.DeleteSampleCode(ctx)
}
