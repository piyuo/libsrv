package data

import (
	"context"
	"strconv"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

const GreetModelName = "Greet"

type Greet struct {
	StoredObject
	From        string
	Description string
	Value       int
}

//GreetFactory provide function to create instance
var GreetFactory = func() Object {
	return new(Greet)
}

func (g *Greet) ModelName() string {
	return GreetModelName
}

func TestGetWithNoID(t *testing.T) {
	Convey("get object with no id", t, func() {
		ctx := context.Background()
		db, err := firestoreGlobalDB(ctx)
		So(err, ShouldBeNil)
		defer db.Close()

		err = db.Get(ctx, &Greet{})
		So(err, ShouldNotBeNil)
	})
}

func TestGetWithNotExistID(t *testing.T) {
	Convey("get object with id not exist", t, func() {
		ctx := context.Background()
		db, err := firestoreGlobalDB(ctx)
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
	db, err := firestoreGlobalDB(ctx)
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
	Convey("delete greet from datastore", t, func() {
		So(err, ShouldBeNil)
	})
}

func TestDeleteByID(t *testing.T) {
	Convey("should delete object by id", t, func() {
		greet := Greet{
			From:        "me",
			Description: "hi",
		}
		ctx := context.Background()
		db, err := firestoreGlobalDB(ctx)
		defer db.Close()
		db.DeleteAll(ctx, GreetModelName, 9)
		db.Put(ctx, &greet)
		exist, err := db.ExistByID(ctx, GreetModelName, greet.ID())
		So(exist, ShouldBeTrue)
		So(err, ShouldBeNil)
		err = db.DeleteByID(ctx, GreetModelName, greet.ID())
		So(err, ShouldBeNil)
		exist, err = db.ExistByID(ctx, GreetModelName, greet.ID())
		So(exist, ShouldBeFalse)
		So(err, ShouldBeNil)
		db.DeleteAll(ctx, GreetModelName, 9)
	})
}

func TestGetPutDeleteWhenContextCanceled(t *testing.T) {
	Convey("should get error when context canceled", t, func() {
		greet := Greet{
			From:        "me",
			Description: "hi",
		}
		db, err := firestoreGlobalDB(context.Background())
		So(err, ShouldBeNil)

		dateline := time.Now().Add(time.Duration(1) * time.Millisecond)
		ctx, cancel := context.WithDeadline(context.Background(), dateline)
		defer cancel()
		time.Sleep(time.Duration(2) * time.Millisecond)

		err = db.Put(ctx, &greet)
		So(err, ShouldNotBeNil)
		err = db.Get(ctx, &greet)
		So(err, ShouldNotBeNil)
		err = db.GetAll(ctx, GreetFactory, 100, func(o Object) {})
		So(err, ShouldNotBeNil)

	})
}

func TestUpdate(t *testing.T) {
	greet := Greet{
		From:        "me",
		Description: "hello",
	}

	ctx := context.Background()
	db, _ := firestoreGlobalDB(ctx)
	defer db.Close()

	err := db.Put(ctx, &greet)
	Convey("put sample ", t, func() {
		So(err, ShouldBeNil)
	})

	err = db.Update(ctx, greet.ModelName(), map[string]interface{}{
		"Description": "helloworld",
	}, greet.ID())
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
	db, _ := firestoreGlobalDB(ctx)
	defer db.Close()
	db.DeleteAll(ctx, GreetModelName, 9)

	db.Put(ctx, &greet1)
	db.Put(ctx, &greet2)

	//test select
	qry := db.Select(ctx, func() Object {
		return new(Greet)
	})

	list, err := qry.Execute()
	Convey("test select count", t, func() {
		So(len(list), ShouldEqual, 2)
	})

	Convey("delete select document ", t, func() {
		err = db.Delete(ctx, list[0])
		So(err, ShouldBeNil)
		err = db.Delete(ctx, list[1])
		So(err, ShouldBeNil)
	})

	db.DeleteAll(ctx, GreetModelName, 9)
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
	db, _ := firestoreGlobalDB(ctx)
	defer db.Close()
	db.Put(ctx, &greet1)
	db.Put(ctx, &greet2)

	numDeleted, err := db.DeleteAll(ctx, greet1.ModelName(), 10)

	Convey("DeleteAll should not have error", t, func() {
		So(err, ShouldBeNil)
	})

	Convey("DeleteAll should return count", t, func() {
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
	db, _ := firestoreGlobalDB(ctx)
	defer db.Close()
	db.DeleteAll(ctx, greet1.ModelName(), 9)
	db.Put(ctx, &greet1)
	db.Put(ctx, &greet2)

	list := []*Greet{}
	err := db.GetAll(ctx, GreetFactory, 100, func(o Object) {
		list = append(list, o.(*Greet))
	})
	Convey("GetAll should not have error", t, func() {
		So(err, ShouldBeNil)
	})
	Convey("list should hold all object", t, func() {
		So(len(list), ShouldEqual, 2)
	})
	db.DeleteAll(ctx, greet1.ModelName(), 9)
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
	db, _ := firestoreGlobalDB(ctx)
	defer db.Close()
	db.DeleteAll(ctx, greet1.ModelName(), 9)
	db.Put(ctx, &greet1)
	db.Put(ctx, &greet2)

	list, err := db.ListAll(ctx, GreetFactory, 100)
	Convey("GetAll should not have error", t, func() {
		So(err, ShouldBeNil)
	})
	Convey("list should hold all object", t, func() {
		So(len(list), ShouldEqual, 2)
	})
	db.DeleteAll(ctx, greet1.ModelName(), 9)
}

func TestExist(t *testing.T) {
	Convey("Should check object exist", t, func() {
		greet := Greet{
			From:        "1",
			Description: "1",
		}
		ctx := context.Background()
		db, _ := firestoreGlobalDB(ctx)
		defer db.Close()
		db.DeleteAll(ctx, GreetModelName, 9)

		exist, err := db.Exist(ctx, GreetModelName, "From", "==", "1")
		So(exist, ShouldBeFalse)
		So(err, ShouldBeNil)

		db.Put(ctx, &greet)

		exist, err = db.Exist(ctx, GreetModelName, "From", "==", "1")
		So(exist, ShouldBeTrue)
		So(err, ShouldBeNil)

		db.DeleteAll(ctx, GreetModelName, 9)
	})
}

func TestExistByID(t *testing.T) {
	Convey("Should check object exist by id", t, func() {
		greet := Greet{
			From:        "1",
			Description: "1",
		}
		ctx := context.Background()
		db, _ := firestoreGlobalDB(ctx)
		defer db.Close()
		db.DeleteAll(ctx, GreetModelName, 9)

		exist, err := db.ExistByID(ctx, GreetModelName, "mockID")
		So(exist, ShouldBeFalse)
		So(err, ShouldBeNil)

		db.Put(ctx, &greet)
		exist, err = db.ExistByID(ctx, GreetModelName, greet.ID())
		So(exist, ShouldBeTrue)
		So(err, ShouldBeNil)

		db.DeleteAll(ctx, GreetModelName, 9)
	})
}

func TestCount(t *testing.T) {
	Convey("Should Count object", t, func() {
		greet := Greet{
			From:        "1",
			Description: "1",
		}
		ctx := context.Background()
		db, _ := firestoreGlobalDB(ctx)
		defer db.Close()
		db.DeleteAll(ctx, GreetModelName, 9)

		count, err := db.Count10(ctx, GreetModelName, "From", "==", "1")
		So(count, ShouldEqual, 0)
		So(err, ShouldBeNil)

		db.Put(ctx, &greet)

		count, err = db.Count10(ctx, GreetModelName, "From", "==", "1")
		So(count, ShouldEqual, 1)
		So(err, ShouldBeNil)

		db.DeleteAll(ctx, GreetModelName, 9)

		count, err = db.Count10(ctx, GreetModelName, "From", "==", "1")
		So(count, ShouldEqual, 0)
		So(err, ShouldBeNil)

	})
}

func TestIncrement(t *testing.T) {
	Convey("Should increment object field", t, func() {
		greet := Greet{
			Value: 0,
		}
		ctx := context.Background()
		db, _ := firestoreGlobalDB(ctx)
		defer db.Close()
		db.DeleteAll(ctx, GreetModelName, 9)
		db.Put(ctx, &greet)

		err := db.Get(ctx, &greet)
		So(greet.Value, ShouldEqual, 0)
		So(err, ShouldBeNil)

		err = db.Increment(ctx, GreetModelName, "Value", greet.ID(), 2)

		err = db.Get(ctx, &greet)
		So(greet.Value, ShouldEqual, 2)
		So(err, ShouldBeNil)
	})
}

func TestCounter(t *testing.T) {
	Convey("Should init, increment, count on counter", t, func() {
		ctx := context.Background()
		db, _ := firestoreGlobalDB(ctx)
		defer db.Close()
		db.DeleteAll(ctx, "Counter", 9)

		counter, err := db.Counter(ctx, "mockCounter", 10)
		So(counter, ShouldNotBeNil)
		So(err, ShouldBeNil)

		count, err := counter.Count(ctx)
		So(count, ShouldEqual, 0)
		So(err, ShouldBeNil)

		err = counter.Increment(ctx, 2)
		So(err, ShouldBeNil)

		count, err = counter.Count(ctx)
		So(count, ShouldEqual, 2)
		So(err, ShouldBeNil)

		err = db.DeleteCounter(ctx, "mockCounter")
		So(err, ShouldBeNil)
	})
	Convey("Should delete not exist counter", t, func() {
		ctx := context.Background()
		db, _ := firestoreGlobalDB(ctx)
		defer db.Close()
		db.DeleteAll(ctx, "Counter", 9)

		err := db.DeleteCounter(ctx, "notExistCounter")
		So(err, ShouldBeNil)
	})
}

func BenchmarkPutSpeed(b *testing.B) {
	greet := Greet{
		From:        "me",
		Description: "hello",
	}
	ctx := context.Background()
	db, _ := firestoreGlobalDB(ctx)
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
	db, _ := firestoreGlobalDB(ctx)
	defer db.Close()

	err := db.Put(ctx, &greet)
	if err != nil {
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		db.Update(ctx, greet.ModelName(), map[string]interface{}{
			"Description": "hello" + strconv.Itoa(i),
		}, greet.ID())
	}
}
