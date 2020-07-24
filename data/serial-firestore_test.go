package data

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	util "github.com/piyuo/libsrv/util"
	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSerial(t *testing.T) {
	Convey("check serial function", t, func() {
		ctx := context.Background()
		dbG, dbR := createSampleDB()
		defer removeSampleDB(dbG, dbR)
		serialsG, serialsR := createSampleSerials(dbG, dbR)
		defer removeSampleSerials(serialsG, serialsR)

		serialMustUseWithInTransacton(ctx, serialsG)
		serialMustUseWithInTransacton(ctx, serialsR)

		serialInTransactionTest(ctx, dbG, serialsG)
		serialInTransactionTest(ctx, dbR, serialsR)

		testSerialInCanceledCtx(ctx, dbG, serialsG)
		testSerialInCanceledCtx(ctx, dbR, serialsR)

	})
}

func testSerialInCanceledCtx(ctx context.Context, db SampleDB, serials *SampleSerials) {
	serial := serials.SampleSerial()
	So(serial, ShouldNotBeNil)

	ctxCanceled := util.CanceledCtx()
	err := serial.Reset(ctxCanceled)
	So(err, ShouldNotBeNil)
}

func serialMustUseWithInTransacton(ctx context.Context, serials *SampleSerials) {
	serial := serials.SampleSerial()

	num, err := serial.NumberRX()
	So(err, ShouldNotBeNil)
	So(num, ShouldEqual, 0)
	err = serial.NumberWX()
	So(err, ShouldNotBeNil)
}

func serialInTransactionTest(ctx context.Context, db SampleDB, serials *SampleSerials) {
	serial := serials.SampleSerial()
	err := serials.DeleteSampleSerial(ctx)
	So(err, ShouldBeNil)
	defer serials.DeleteSampleSerial(ctx)

	err = db.Transaction(ctx, func(ctx context.Context) error {
		num, err := serial.NumberRX()
		So(err, ShouldBeNil)
		So(num, ShouldEqual, 1)
		return serial.NumberWX()
	})
	So(err, ShouldBeNil)

	err = db.Transaction(ctx, func(ctx context.Context) error {
		num, err := serial.NumberRX()
		So(err, ShouldBeNil)
		So(num, ShouldEqual, 2)
		err = serial.NumberWX()
		So(err, ShouldBeNil)
		return errors.New("make fail transaction")
	})
	So(err, ShouldNotBeNil)

	err = db.Transaction(ctx, func(ctx context.Context) error {
		num, err := serial.NumberRX()
		So(err, ShouldBeNil)
		So(num, ShouldEqual, 2)
		return serial.NumberWX()
	})
	So(err, ShouldBeNil)

	err = db.Transaction(ctx, func(ctx context.Context) error {
		num, err := serial.NumberRX()
		So(err, ShouldBeNil)
		So(num, ShouldEqual, 3)
		return serial.NumberWX()
	})

	// reset serial
	err = serial.Reset(ctx)
	So(err, ShouldBeNil)

	err = db.Transaction(ctx, func(ctx context.Context) error {
		num, err := serial.NumberRX()
		So(err, ShouldBeNil)
		So(num, ShouldEqual, 1)
		return serial.NumberWX()
	})

	// reset in transaction
	err = db.Transaction(ctx, func(ctx context.Context) error {
		return serial.Reset(ctx)
	})

	err = db.Transaction(ctx, func(ctx context.Context) error {
		num, err := serial.NumberRX()
		So(err, ShouldBeNil)
		So(num, ShouldEqual, 1)
		return serial.NumberWX()
	})

	So(err, ShouldBeNil)
}

func TestConcurrentserial(t *testing.T) {
	ctx := context.Background()
	rand.Seed(time.Now().UnixNano())
	dbG, dbR := createSampleDB()
	defer removeSampleDB(dbG, dbR)
	serialsG, serialsR := createSampleSerials(dbG, dbR)
	defer removeSampleSerials(serialsG, serialsR)

	var concurrent = 3
	var wg sync.WaitGroup
	wg.Add(concurrent)
	createserial := func() {
		db, _ := NewSampleGlobalDB(ctx)
		defer db.Close()
		serials := db.Serials()
		time.Sleep(time.Duration(rand.Intn(2)) * time.Second)

		for i := 0; i < 3; i++ {
			serial := serials.SampleSerial()

			err := db.Transaction(ctx, func(ctx context.Context) error {
				num, err := serial.NumberRX()
				if err != nil {
					t.Errorf("err should be nil, got %v", err)
				}
				fmt.Printf("num:%v\n", num)
				return serial.NumberWX()
			})
			if err != nil {
				t.Errorf("err should be nil, got %v", err)
			}
			// serial update need to be low frequency
			time.Sleep(time.Duration(rand.Intn(2)) * time.Second)
		}
		wg.Done()
	}

	//create go routing to do counting
	for i := 0; i < concurrent; i++ {
		go createserial()
	}
	wg.Wait()

	err := dbG.Transaction(ctx, func(ctx context.Context) error {
		serial := serialsG.SampleSerial()
		num, err := serial.NumberRX()
		if err != nil {
			t.Errorf("err should be nil, got %v", err)
		}
		if num != 10 {
			t.Errorf("serial = %d; want 10", num)
		}
		return serial.NumberWX()
	})
	if err != nil {
		t.Errorf("err should be nil, got %v", err)
		return
	}
}
