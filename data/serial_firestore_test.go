package data

import (
	"context"
	"fmt"
	"sync"
	"testing"

	util "github.com/piyuo/libsrv/util"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSerial(t *testing.T) {
	Convey("check serial function", t, func() {
		ctx := context.Background()
		dbG, dbR := createSampleDB()
		defer removeSampleDB(dbG, dbR)
		serialsG, serialsR := createSampleSerials(dbG, dbR)
		defer removeSampleSerials(serialsG, serialsR)

		serialTest(ctx, serialsG)
		serialTest(ctx, serialsR)
		serialInTransactionTest(ctx, dbG, serialsG)
		serialInTransactionTest(ctx, dbR, serialsR)
		serialContextCanceled(serialsG)
		serialContextCanceled(serialsR)
	})
}

func serialTest(ctx context.Context, serials *SampleSerials) {
	serial := serials.SampleSerial()
	err := serials.DeleteSampleSerial(ctx)
	So(err, ShouldBeNil)

	num, err := serial.Number(ctx)
	So(err, ShouldBeNil)
	So(num, ShouldEqual, 1)

	num, err = serial.Number(ctx)
	So(err, ShouldBeNil)
	So(num, ShouldEqual, 2)

	num, err = serial.Number(ctx)
	So(err, ShouldBeNil)
	So(num, ShouldEqual, 3)

	err = serials.DeleteSampleSerial(ctx)
	So(err, ShouldBeNil)
}

func serialInTransactionTest(ctx context.Context, db SampleDB, serials *SampleSerials) {
	serial := serials.SampleSerial()
	err := serials.DeleteSampleSerial(ctx)
	So(err, ShouldBeNil)

	err = db.Transaction(ctx, func(ctx context.Context) error {
		num, err := serial.Number(ctx)
		So(err, ShouldBeNil)
		So(num, ShouldEqual, 1)
		return nil
	})
	So(err, ShouldBeNil)

	err = db.Transaction(ctx, func(ctx context.Context) error {
		num, err := serial.Number(ctx)
		So(err, ShouldBeNil)
		So(num, ShouldEqual, 2)
		// this line will cause read after write error, only one serial can generated in transaction
		num, err = serial.Number(ctx)
		So(err, ShouldNotBeNil)
		So(num, ShouldEqual, 0)
		return err
	})
	So(err, ShouldNotBeNil)

	err = db.Transaction(ctx, func(ctx context.Context) error {
		err := serials.DeleteSampleSerial(ctx)
		So(err, ShouldBeNil)
		return nil
	})
	So(err, ShouldBeNil)
}

func serialContextCanceled(serials *SampleSerials) {
	ctx := util.CanceledCtx()

	serial := serials.SampleSerial()
	num, err := serial.Number(ctx)
	So(err, ShouldNotBeNil)
	So(num, ShouldEqual, 0)

	err = serials.DeleteSampleSerial(ctx)
	So(err, ShouldNotBeNil)
}

func TestConcurrentSerial(t *testing.T) {
	ctx := context.Background()
	dbG, dbR := createSampleDB()
	defer removeSampleDB(dbG, dbR)
	serialsG, serialsR := createSampleSerials(dbG, dbR)
	defer removeSampleSerials(serialsG, serialsR)

	serialG := serialsG.SampleSerial()

	var concurrent = 3
	var wg sync.WaitGroup
	wg.Add(concurrent)
	createSerial := func() {
		db, _ := NewSampleGlobalDB(ctx)
		defer db.Close()
		serials := db.Serials()

		for i := 0; i < 5; i++ {
			serial := serials.SampleSerial()
			num, err := serial.Number(ctx)
			if err != nil {
				t.Errorf("err should be nil, got %v", err)
			}

			if err != nil {
				t.Errorf("err should be nil, got %v", err)
			}
			fmt.Printf("num:%v\n", num)
		}
		wg.Done()
	}

	//create go routing to do counting
	for i := 0; i < concurrent; i++ {
		go createSerial()
	}
	wg.Wait()
	num, err := serialG.Number(ctx)
	if err != nil {
		t.Errorf("err should be nil, got %v", err)
		return
	}
	if num != 16 {
		t.Errorf("serial = %d; want 16", num)
	}
}
