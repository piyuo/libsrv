package data

import (
	"context"
	"strconv"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

type Greet struct {
	object
	From        string
	Description string
}

//GreetFactory provide function to create instance
var GreetFactory = func() Object {
	return new(Greet)
}

func (g *Greet) Class() string {
	return "Greet"
}

func TestGetWithNoID(t *testing.T) {
	Convey("get object with no id", t, func() {
		ctx := context.Background()
		db, err := firestoreNewDB(ctx)
		So(err, ShouldBeNil)
		defer db.Close()

		err = db.Get(ctx, &Greet{})
		So(err, ShouldNotBeNil)
	})
}

func TestGetWithNotExistID(t *testing.T) {
	Convey("get object with id not exist", t, func() {
		ctx := context.Background()
		db, err := firestoreNewDB(ctx)
		So(err, ShouldBeNil)
		defer db.Close()

		greet := &Greet{}
		greet.SetID("notexist")
		err = db.Get(ctx, greet)
		So(err, ShouldEqual, ErrObjectNotFound)
	})
}

func TestPutGetDelete(t *testing.T) {
	greet := Greet{
		From:        "me",
		Description: "hi",
	}
	ctx := context.Background()
	db, err := firestoreNewDB(ctx)
	Convey("should get db without error", t, func() {
		So(err, ShouldBeNil)
	})

	defer db.Close()

	err = db.Put(ctx, &greet)
	Convey("greet should have id after put", t, func() {
		So(err, ShouldBeNil)
	})

	objID := greet.ID()
	Convey("greet ID should set", t, func() {
		So(objID, ShouldNotBeEmpty)
	})

	Convey("object load from datastore should equal to insert", t, func() {
		greet2 := Greet{}
		greet2.SetID(objID)
		err = db.Get(ctx, &greet2)
		So(err, ShouldBeNil)
		So(greet2.From, ShouldEqual, greet.From)
	})

	//test delete
	err = db.Delete(ctx, &greet)
	Convey("delete greet from datastore'", t, func() {
		So(err, ShouldBeNil)
	})
}

func TestGetPutDeleteWhenContextCanceled(t *testing.T) {
	Convey("should get error when context canceled", t, func() {
		greet := Greet{
			From:        "me",
			Description: "hi",
		}
		db, err := firestoreNewDB(context.Background())
		So(err, ShouldBeNil)

		dateline := time.Now().Add(time.Duration(1) * time.Millisecond)
		ctx, cancel := context.WithDeadline(context.Background(), dateline)
		defer cancel()
		time.Sleep(time.Duration(2) * time.Millisecond)

		err = db.Put(ctx, &greet)
		So(err, ShouldNotBeNil)
		err = db.Get(ctx, &greet)
		So(err, ShouldNotBeNil)
		err = db.GetByClass(ctx, "Greet", &greet)
		So(err, ShouldNotBeNil)
		err = db.GetAll(ctx, GreetFactory, func(o Object) {}, 100)
		So(err, ShouldNotBeNil)

	})
}

func TestUpdate(t *testing.T) {
	greet := Greet{
		From:        "me",
		Description: "hello",
	}

	ctx := context.Background()
	db, _ := firestoreNewDB(ctx)
	defer db.Close()

	err := db.Put(ctx, &greet)
	Convey("put sample ", t, func() {
		So(err, ShouldBeNil)
	})

	err = db.Update(ctx, greet.Class(), greet.ID(), map[string]interface{}{
		"Description": "helloworld",
	})
	Convey("update sample description", t, func() {
		So(err, ShouldBeNil)
	})

	Convey("sample description should be updated", t, func() {
		greet2 := Greet{}
		greet2.SetID(greet.ID())
		err = db.Get(ctx, &greet2)
		So(err, ShouldBeNil)
		So(greet2.Description, ShouldEqual, "helloworld")
	})

	_ = db.Delete(ctx, &greet)
}

func TestSelectDelete(t *testing.T) {
	greet1 := Greet{
		From:        "1",
		Description: "1",
	}
	greet2 := Greet{
		From:        "2",
		Description: "2",
	}
	ctx := context.Background()
	db, _ := firestoreNewDB(ctx)
	defer db.Close()

	db.Put(ctx, &greet1)
	db.Put(ctx, &greet2)

	//test select
	qry := db.Select(ctx, func() Object {
		return new(Greet)
	})

	var i int
	qry.Run(func(o Object) {
		i++
		err := db.Delete(ctx, o)
		Convey("delete select document ", t, func() {
			So(err, ShouldBeNil)
		})
	})
	Convey("test select count '", t, func() {
		So(i, ShouldBeGreaterThanOrEqualTo, 2)
	})
}

func TestDeleteAll(t *testing.T) {
	greet1 := Greet{
		From:        "1",
		Description: "1",
	}
	greet2 := Greet{
		From:        "2",
		Description: "2",
	}
	ctx := context.Background()
	db, _ := firestoreNewDB(ctx)
	defer db.Close()
	db.Put(ctx, &greet1)
	db.Put(ctx, &greet2)

	numDeleted, err := db.DeleteAll(ctx, greet1.Class(), 10)

	Convey("DeleteAll should not have error '", t, func() {
		So(err, ShouldBeNil)
	})

	Convey("DeleteAll should return count '", t, func() {
		So(numDeleted, ShouldBeGreaterThanOrEqualTo, 2)
	})
}

func TestGetAll(t *testing.T) {
	greet1 := Greet{
		From:        "1",
		Description: "1",
	}
	greet2 := Greet{
		From:        "2",
		Description: "2",
	}
	ctx := context.Background()
	db, _ := firestoreNewDB(ctx)
	defer db.Close()
	db.DeleteAll(ctx, greet1.Class(), 9)
	db.Put(ctx, &greet1)
	db.Put(ctx, &greet2)

	list := []*Greet{}
	err := db.GetAll(ctx, GreetFactory, func(o Object) {
		list = append(list, o.(*Greet))
	}, 100)
	Convey("GetAll should not have error '", t, func() {
		So(err, ShouldBeNil)
	})
	Convey("list should hold all object", t, func() {
		So(len(list), ShouldEqual, 2)
	})
	db.DeleteAll(ctx, greet1.Class(), 9)
}

func TestListAll(t *testing.T) {
	greet1 := Greet{
		From:        "1",
		Description: "1",
	}
	greet2 := Greet{
		From:        "2",
		Description: "2",
	}
	ctx := context.Background()
	db, _ := firestoreNewDB(ctx)
	defer db.Close()
	db.DeleteAll(ctx, greet1.Class(), 9)
	db.Put(ctx, &greet1)
	db.Put(ctx, &greet2)

	list, err := db.ListAll(ctx, GreetFactory, 100)
	Convey("GetAll should not have error '", t, func() {
		So(err, ShouldBeNil)
	})
	Convey("list should hold all object", t, func() {
		So(len(list), ShouldEqual, 2)
	})
	db.DeleteAll(ctx, greet1.Class(), 9)
}

func BenchmarkPutSpeed(b *testing.B) {
	greet := Greet{
		From:        "me",
		Description: "hello",
	}
	ctx := context.Background()
	db, _ := firestoreNewDB(ctx)
	defer db.Close()

	err := db.Put(ctx, &greet)
	if err != nil {
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		greet.Description = "hello" + strconv.Itoa(i)
		db.Put(ctx, &greet)
	}
}

func BenchmarkUpdateSpeed(b *testing.B) {
	greet := Greet{
		From:        "me",
		Description: "hello",
	}
	ctx := context.Background()
	db, _ := firestoreNewDB(ctx)
	defer db.Close()

	err := db.Put(ctx, &greet)
	if err != nil {
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		db.Update(ctx, greet.Class(), greet.ID(), map[string]interface{}{
			"Description": "hello" + strconv.Itoa(i),
		})
	}
}