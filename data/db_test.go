package data

import (
	"context"

	. "github.com/smartystreets/goconvey/convey"
)

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
	err := g.Clear(ctx)
	So(err, ShouldBeNil)
	err = r.Clear(ctx)
	So(err, ShouldBeNil)
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

func createSampleCodes(dbG *SampleGlobalDB, dbR *SampleRegionalDB) (*SampleCodes, *SampleCodes) {
	g := dbG.Codes()
	r := dbR.Codes()
	removeSampleCodes(g, r)
	return g, r
}

func removeSampleCodes(g *SampleCodes, r *SampleCodes) {
	ctx := context.Background()
	g.DeleteSampleCode(ctx)
	r.DeleteSampleCode(ctx)
}
