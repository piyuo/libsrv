package data

import (
	"context"
	"strconv"
	"testing"

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

func TestGetNotFound(t *testing.T) {
	ctx := context.Background()
	db, _ := NewFirestoreDB(ctx)
	defer db.Close()
	err := db.Get(ctx, &Greet{})
	Convey("get not exist object", t, func() {
		So(err, ShouldNotBeNil)
	})
}

func TestPutGetDelete(t *testing.T) {
	greet := Greet{
		From:        "me",
		Description: "hello",
	}
	ctx := context.Background()
	db, _ := NewFirestoreDB(ctx)
	defer db.Close()

	//test put
	err := db.Put(ctx, &greet)
	Convey("greet should have id after db.put", t, func() {
		So(err, ShouldBeNil)
	})

	objID := greet.ID()
	Convey("greet ID should be set", t, func() {
		So(objID, ShouldNotBeEmpty)
	})

	//test get
	greet2 := Greet{}
	greet2.SetID(objID)
	_ = db.Get(ctx, &greet2)
	Convey("object load from datastore should equal to insert", t, func() {
		So(greet2.From, ShouldEqual, greet.From)
	})

	//test delete
	err = db.Delete(ctx, &greet)
	Convey("delete greet from datastore'", t, func() {
		So(err, ShouldBeNil)
	})

}

func TestUpdate(t *testing.T) {
	greet := Greet{
		From:        "me",
		Description: "hello",
	}

	ctx := context.Background()
	db, _ := NewFirestoreDB(ctx)
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
	db, _ := NewFirestoreDB(ctx)
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
	db, _ := NewFirestoreDB(ctx)
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
	db, _ := NewFirestoreDB(ctx)
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
	db, _ := NewFirestoreDB(ctx)
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
	db, _ := NewFirestoreDB(ctx)
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
	db, _ := NewFirestoreDB(ctx)
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
