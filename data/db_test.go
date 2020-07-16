package data

import (
	"context"
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
