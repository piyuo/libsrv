package data

import (
	"context"
	"strconv"
	"testing"

	gcp "github.com/piyuo/libsrv/gcp"
	util "github.com/piyuo/libsrv/util"
	. "github.com/smartystreets/goconvey/convey"
)

func TestFirestoreNewDB(t *testing.T) {
	Convey("should create db", t, func() {
		ctx := context.Background()
		cred, err := gcp.GlobalCredential(ctx)
		So(err, ShouldBeNil)
		db, err := firestoreNewConnection(ctx, cred, "")
		defer db.Close()
		So(err, ShouldBeNil)
		So(db, ShouldNotBeNil)
	})
}

func TestFirestoreGlobalDB(t *testing.T) {
	Convey("should create global db", t, func() {
		ctx := context.Background()
		conn, err := FirestoreGlobalConnection(ctx)
		defer conn.Close()
		So(err, ShouldBeNil)
		So(conn, ShouldNotBeNil)

		// error id should use root if not having namespace
		firestoreConn := conn.(*ConnectionFirestore)
		id := firestoreConn.errorID("tablename", "")
		So(id, ShouldEqual, "tablename{root}")
		id = firestoreConn.errorID("tablename", "id")
		So(id, ShouldEqual, "tablename{root}-id")

		// global has no namespace to be create or delete
		err = firestoreConn.CreateNamespace(ctx)
		So(err, ShouldNotBeNil)
		err = firestoreConn.DeleteNamespace(ctx)
		So(err, ShouldNotBeNil)

	})
}

func TestFirestoreRegionalDB(t *testing.T) {
	Convey("should create regional db", t, func() {
		ctx := context.Background()
		db, err := FirestoreRegionalConnection(ctx, "sample-namespace")
		defer db.Close()
		So(err, ShouldBeNil)
		So(db, ShouldNotBeNil)

		firestoreDB := db.(*ConnectionFirestore)
		id := firestoreDB.errorID("tablename", "")
		So(id, ShouldEqual, "tablename{sample-namespace}")
		id = firestoreDB.errorID("tablename", "id")
		So(id, ShouldEqual, "tablename{sample-namespace}-id")

		err = firestoreDB.snapshotToObject("tableName", nil, nil, nil)
		So(err, ShouldNotBeNil)

	})
}

func TestNameSpace(t *testing.T) {
	Convey("should create name space", t, func() {
		ctx := context.Background()
		conn, err := FirestoreRegionalConnection(ctx, "sample-namespace")
		defer conn.Close()
		So(err, ShouldBeNil)
		So(conn, ShouldNotBeNil)

		err = conn.CreateNamespace(ctx)
		So(err, ShouldBeNil)
		err = conn.DeleteNamespace(ctx)
		So(err, ShouldBeNil)

		ctxCanceled := util.CanceledCtx()
		err = conn.CreateNamespace(ctxCanceled)
		So(err, ShouldNotBeNil)
		err = conn.DeleteNamespace(ctxCanceled)
		So(err, ShouldNotBeNil)
	})
}

func TestConnection(t *testing.T) {
	Convey("test genreal operation on connection", t, func() {
		ctx := context.Background()
		dbG, dbR := createSampleDB()
		defer removeSampleDB(dbG, dbR)
		samplesG, samplesR := createSampleTable(dbG, dbR)
		defer removeSampleTable(samplesG, samplesR)

		testGroup(ctx, samplesG)
		testGroup(ctx, samplesR)

	})
}

func testGroup(ctx context.Context, table *Table) {
	testID(ctx, table)
	testSetGetExistDelete(ctx, table)
	testSelectUpdateIncrementDelete(ctx, table)
	testListQueryFindCountClear(ctx, table)
	testDelete(ctx, table)
	testConnectionContextCanceled(table)
	testSearchCountIsEmpty(ctx, table)
}

func testID(ctx context.Context, table *Table) {
	sample := &Sample{
		Name:  "sample",
		Value: 1,
	}
	So(sample.ID, ShouldBeEmpty)

	o, err := table.Get(ctx, "")
	So(err, ShouldBeNil)
	So(o, ShouldBeNil)

	// auto id
	err = table.Set(ctx, sample)
	So(err, ShouldBeNil)
	So(sample.ID, ShouldNotBeEmpty)

	sample2, err := table.Get(ctx, sample.ID)
	So(err, ShouldBeNil)
	So(sample2, ShouldNotBeNil)
	So(sample.Name, ShouldEqual, sample2.(*Sample).Name)
	So(sample2.GetCreateTime().IsZero(), ShouldBeFalse)
	So(sample2.GetUpdateTime().IsZero(), ShouldBeFalse)
	So(sample2.GetReadTime().IsZero(), ShouldBeFalse)

	// factory has no object return must error
	bakFactory := table.Factory
	table.Factory = func() Object {
		return nil
	}
	sampleX, err := table.Get(ctx, sample.ID)
	So(err, ShouldNotBeNil)
	So(sampleX, ShouldBeNil)
	table.Factory = bakFactory

	// set sample again
	sample.Name = "modified"
	err = table.Set(ctx, sample)
	So(err, ShouldBeNil)

	m, err := table.Get(ctx, sample.ID)
	sampleM := m.(*Sample)
	So(err, ShouldBeNil)
	So(sampleM, ShouldNotBeNil)
	So(sampleM.Name, ShouldEqual, "modified")

	// set nil object
	err = table.Set(ctx, nil)
	So(err, ShouldNotBeNil)

	err = table.DeleteObject(ctx, sample2)
	So(err, ShouldBeNil)

	// manual id
	sample = &Sample{
		Name:  "sample",
		Value: 1,
	}
	sample.ID = "sample-id"
	err = table.Set(ctx, sample)
	So(err, ShouldBeNil)
	So(sample.ID, ShouldEqual, "sample-id")

	sample3, err := table.Get(ctx, "sample-id")
	So(err, ShouldBeNil)
	So(sample3, ShouldNotBeNil)
	So(sample.Name, ShouldEqual, sample3.(*Sample).Name)

	err = table.DeleteObject(ctx, sample3)
	So(err, ShouldBeNil)

}

func testSetGetExistDelete(ctx context.Context, table *Table) {
	sample := &Sample{
		Name:  "sample",
		Value: 1,
	}

	err := table.Set(ctx, sample)
	So(err, ShouldBeNil)
	sampleID := sample.ID
	sample2, err := table.Get(ctx, sampleID)
	So(err, ShouldBeNil)
	So(sample2, ShouldNotBeNil)
	So(sample.Name, ShouldEqual, sample2.(*Sample).Name)

	exist, err := table.Exist(ctx, sampleID)
	So(err, ShouldBeNil)
	So(exist, ShouldBeTrue)

	exist, err = table.Exist(ctx, "")
	So(err, ShouldBeNil)
	So(exist, ShouldBeFalse)

	err = table.Delete(ctx, sampleID)
	So(err, ShouldBeNil)

	exist, err = table.Exist(ctx, sampleID)
	So(err, ShouldBeNil)
	So(exist, ShouldBeFalse)

	sample3, err := table.Get(ctx, sampleID)
	So(err, ShouldBeNil)
	So(sample3, ShouldBeNil)

	err = table.Clear(ctx)
	So(err, ShouldBeNil)
}

func testSelectUpdateIncrementDelete(ctx context.Context, table *Table) {
	sample := &Sample{
		Name:  "sample",
		Value: 6,
	}
	err := table.Set(ctx, sample)
	So(err, ShouldBeNil)

	value, err := table.Select(ctx, "NotExistID", "Value")
	So(err, ShouldBeNil)
	So(value, ShouldBeNil)

	value, err = table.Select(ctx, sample.ID, "Value")
	So(err, ShouldBeNil)
	So(value, ShouldEqual, 6)

	err = table.Update(ctx, "NotExistID", map[string]interface{}{
		"Name":  "sample2",
		"Value": 2,
	})
	So(err, ShouldBeNil)

	err = table.Delete(ctx, "NotExistID")
	So(err, ShouldBeNil)

	err = table.Update(ctx, sample.ID, map[string]interface{}{
		"Name":  "sample2",
		"Value": 2,
	})
	So(err, ShouldBeNil)

	name, err := table.Select(ctx, sample.ID, "Name")
	So(err, ShouldBeNil)
	So(name, ShouldEqual, "sample2")

	value, err = table.Select(ctx, sample.ID, "Value")
	So(err, ShouldBeNil)
	So(value, ShouldEqual, 2)

	err = table.Increment(ctx, "NotExistID", "Value", 3)
	So(err, ShouldNotBeNil)

	err = table.Delete(ctx, "NotExistID")
	So(err, ShouldBeNil)

	err = table.Increment(ctx, sample.ID, "Value", 3)
	So(err, ShouldBeNil)

	value, err = table.Select(ctx, sample.ID, "Value")
	So(err, ShouldBeNil)
	So(value, ShouldEqual, 5)

	err = table.DeleteObject(ctx, sample)
	So(err, ShouldBeNil)

}

func testListQueryFindCountClear(ctx context.Context, table *Table) {
	sample1 := &Sample{
		Name:  "sample1",
		Value: 1,
	}
	sample2 := &Sample{
		Name:  "sample2",
		Value: 2,
	}
	err := table.Set(ctx, sample1)
	So(err, ShouldBeNil)
	err = table.Set(ctx, sample2)
	So(err, ShouldBeNil)

	list, err := table.All(ctx)
	So(err, ShouldBeNil)
	So(len(list), ShouldEqual, 2)
	So(list[0].(*Sample).Name, ShouldStartWith, "sample")
	So(list[1].(*Sample).Name, ShouldStartWith, "sample")

	// factory has no object return must error
	bakFactory := table.Factory
	table.Factory = func() Object {
		return nil
	}
	listX, err := table.All(ctx)
	So(err, ShouldNotBeNil)
	So(listX, ShouldBeNil)
	table.Factory = bakFactory

	obj, err := table.Find(ctx, "Name", "==", "sample1")
	So(err, ShouldBeNil)
	So((obj.(*Sample)).Name, ShouldEqual, "sample1")

	list, err = table.Query().OrderBy("Name").Execute(ctx)
	So(err, ShouldBeNil)
	So(len(list), ShouldEqual, 2)
	So(list[0].(*Sample).Name, ShouldEqual, sample1.Name)
	So(list[1].(*Sample).Name, ShouldEqual, sample2.Name)

	obj, err = table.Find(ctx, "Value", "==", 2)
	So(err, ShouldBeNil)
	So((obj.(*Sample)).Name, ShouldEqual, "sample2")

	err = table.Clear(ctx)
	So(err, ShouldBeNil)

	obj, err = table.Find(ctx, "Value", "==", 2)
	So(err, ShouldBeNil)
	So(obj, ShouldBeNil)
}

func testSearchCountIsEmpty(ctx context.Context, table *Table) {
	sample := &Sample{
		Name:  "sample",
		Value: 0,
	}
	err := table.Set(ctx, sample)
	So(err, ShouldBeNil)

	objects, err := table.List(ctx, "Name", "==", "sample")
	So(err, ShouldBeNil)
	So(len(objects), ShouldEqual, 1)

	count, err := table.Count(ctx)
	So(err, ShouldBeNil)
	So(count, ShouldEqual, 1)

	empty, err := table.IsEmpty(ctx)
	So(err, ShouldBeNil)
	So(empty, ShouldEqual, false)

	err = table.DeleteObject(ctx, sample)
	So(err, ShouldBeNil)
}

func testDelete(ctx context.Context, table *Table) {
	sample := &Sample{
		Name:  "sample",
		Value: 0,
	}
	err := table.DeleteObject(ctx, sample)
	So(err, ShouldBeNil)
	err = table.Delete(ctx, "NotExistID")
	So(err, ShouldBeNil)
	err = table.Delete(ctx, "NotExistID")
	So(err, ShouldBeNil)

	err = table.Set(ctx, sample)
	So(err, ShouldBeNil)
	exist, err := table.Exist(ctx, sample.ID)
	So(err, ShouldBeNil)
	So(exist, ShouldBeTrue)

	sample2 := &Sample{}
	sample2.ID = sample.ID
	err = table.DeleteObject(ctx, sample2)
	So(err, ShouldBeNil)
	exist, err = table.Exist(ctx, sample.ID)
	So(err, ShouldBeNil)
	So(exist, ShouldBeFalse)

	//delete batch
	//delete empty batch
	ids := []string{}
	err = table.DeleteBatch(ctx, ids)
	So(err, ShouldBeNil)

	err = table.Set(ctx, sample)
	So(err, ShouldBeNil)
	exist, err = table.Exist(ctx, sample.ID)
	So(err, ShouldBeNil)
	So(exist, ShouldBeTrue)

	ids = []string{sample.ID}
	err = table.DeleteBatch(ctx, ids)
	So(err, ShouldBeNil)
	exist, err = table.Exist(ctx, sample.ID)
	So(err, ShouldBeNil)
	So(exist, ShouldBeFalse)

}

func testConnectionContextCanceled(table *Table) {
	ctx := util.CanceledCtx()
	sample := &Sample{}

	err := table.Set(ctx, sample)
	So(err, ShouldNotBeNil)
	_, err = table.Get(ctx, "notexist")
	So(err, ShouldNotBeNil)
	err = table.Delete(ctx, "notexist")
	So(err, ShouldNotBeNil)
	err = table.DeleteObject(ctx, sample)
	So(err, ShouldNotBeNil)
	err = table.DeleteBatch(ctx, []string{})
	So(err, ShouldNotBeNil)
	_, err = table.All(ctx)
	So(err, ShouldNotBeNil)
	_, err = table.Exist(ctx, "notexist")
	So(err, ShouldNotBeNil)
	_, err = table.Select(ctx, "notexist", "Value")
	So(err, ShouldNotBeNil)
	err = table.Update(ctx, "notexist", map[string]interface{}{
		"Name":  "Sample2",
		"Value": "2",
	})
	So(err, ShouldNotBeNil)
	err = table.Clear(ctx)
	So(err, ShouldNotBeNil)
	_, err = table.Query().Execute(ctx)
	So(err, ShouldNotBeNil)
	_, err = table.Find(ctx, "Value", "==", "2")
	So(err, ShouldNotBeNil)
	_, err = table.Count(ctx)
	So(err, ShouldNotBeNil)
	err = table.Increment(ctx, "notexist", "Value", 2)
	So(err, ShouldNotBeNil)
	_, err = table.List(ctx, "Name", "==", "1")
	So(err, ShouldNotBeNil)
	_, err = table.SortList(ctx, "Name", "==", "1", "", ASC)
	So(err, ShouldNotBeNil)
	_, err = table.All(ctx)
	So(err, ShouldNotBeNil)
	_, err = table.IsEmpty(ctx)
	So(err, ShouldNotBeNil)
	err = table.Clear(ctx)
	So(err, ShouldNotBeNil)
}

func BenchmarkPutSpeed(b *testing.B) {
	ctx := context.Background()
	dbG, err := NewSampleGlobalDB(ctx)
	defer dbG.Close()
	So(err, ShouldBeNil)
	table := dbG.SampleTable()
	So(table, ShouldBeNil)

	dbR, err := NewSampleRegionalDB(ctx, "sample-namespace")
	defer dbR.Close()
	samplesR := dbR.SampleTable()
	So(samplesR, ShouldBeNil)
	So(err, ShouldBeNil)

	sample := &Sample{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sample.Name = "hello" + strconv.Itoa(i)
		err = table.Set(ctx, sample)
		if err != nil {
			return
		}
	}
	table.DeleteObject(ctx, sample)
}

func BenchmarkUpdateSpeed(b *testing.B) {
	ctx := context.Background()
	dbG, err := NewSampleGlobalDB(ctx)
	defer dbG.Close()
	So(err, ShouldBeNil)
	table := dbG.SampleTable()
	So(table, ShouldBeNil)

	dbR, err := NewSampleRegionalDB(ctx, "sample-namespace")
	defer dbR.Close()
	samplesR := dbR.SampleTable()
	So(samplesR, ShouldBeNil)
	So(err, ShouldBeNil)

	sample := &Sample{}
	err = table.Set(ctx, sample)
	if err != nil {
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		table.Update(ctx, sample.ID, map[string]interface{}{
			"Name": "hello" + strconv.Itoa(i),
		})
	}
	table.DeleteObject(ctx, sample)
}
