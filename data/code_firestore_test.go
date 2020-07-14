package data

import (
	"context"
	"fmt"
	"sync"
	"testing"

	util "github.com/piyuo/libsrv/util"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCode(t *testing.T) {
	Convey("check code function", t, func() {
		ctx := context.Background()
		dbG, dbR := createSampleDB()
		defer removeSampleDB(dbG, dbR)
		codesG, codesR := createSampleCodes(dbG, dbR)
		defer removeSampleCodes(codesG, codesR)

		codeTest(ctx, codesG)
		codeTest(ctx, codesR)
		codeInTransactionTest(ctx, dbG, codesG)
		codeInTransactionTest(ctx, dbR, codesR)
		codeContextCanceled(codesG)
		codeContextCanceled(codesR)
	})
}

func codeTest(ctx context.Context, codes *SampleCodes) {
	code := codes.SampleCode()
	err := codes.DeleteSampleCode(ctx)
	So(err, ShouldBeNil)
	//	code.CreateShards(ctx)

	num, err := code.Number(ctx)
	So(err, ShouldBeNil)
	So(num, ShouldBeGreaterThanOrEqualTo, 10)

	num, err = code.Number(ctx)
	So(err, ShouldBeNil)
	So(num, ShouldBeGreaterThanOrEqualTo, 10)

	num, err = code.Number(ctx)
	So(err, ShouldBeNil)
	So(num, ShouldBeGreaterThanOrEqualTo, 10)

	c, err := code.Code(ctx)
	So(err, ShouldBeNil)
	So(c, ShouldNotBeEmpty)

	c, err = code.Code16(ctx)
	So(err, ShouldBeNil)
	So(c, ShouldNotBeEmpty)

	c, err = code.Code64(ctx)
	So(err, ShouldBeNil)
	So(c, ShouldNotBeEmpty)

	err = codes.DeleteSampleCode(ctx)
	So(err, ShouldBeNil)
}

func codeInTransactionTest(ctx context.Context, db SampleDB, codes *SampleCodes) {
	err := codes.DeleteSampleCode(ctx)
	So(err, ShouldBeNil)

	err = db.Transaction(ctx, func(ctx context.Context) error {
		code := codes.SampleCode()
		//code.CreateShards(ctx)
		num, err := code.Number(ctx)
		So(err, ShouldBeNil)
		So(num, ShouldBeGreaterThanOrEqualTo, 10)
		return nil
	})
	So(err, ShouldBeNil)

	err = db.Transaction(ctx, func(ctx context.Context) error {
		code := codes.SampleCode()
		num, err := code.Number(ctx)
		So(err, ShouldBeNil)
		So(num, ShouldBeGreaterThanOrEqualTo, 10)
		// this line will cause read after write error, only one serial can generated in transaction
		num, err = code.Number(ctx)
		So(err, ShouldNotBeNil)
		So(num, ShouldEqual, 0)
		return nil
	})
	So(err, ShouldNotBeNil)

	err = db.Transaction(ctx, func(ctx context.Context) error {
		code := codes.SampleCode()
		c, err := code.Code16(ctx)
		So(err, ShouldBeNil)
		So(c, ShouldNotBeEmpty)
		return nil
	})
	So(err, ShouldBeNil)

	err = db.Transaction(ctx, func(ctx context.Context) error {
		code := codes.SampleCode()
		c, err := code.Code(ctx)
		So(err, ShouldBeNil)
		So(c, ShouldNotBeEmpty)
		return nil
	})
	So(err, ShouldBeNil)

	err = db.Transaction(ctx, func(ctx context.Context) error {
		code := codes.SampleCode()
		c, err := code.Code64(ctx)
		So(err, ShouldBeNil)
		So(c, ShouldNotBeEmpty)
		return nil
	})
	So(err, ShouldBeNil)

	err = db.Transaction(ctx, func(ctx context.Context) error {
		err := codes.DeleteSampleCode(ctx)
		So(err, ShouldBeNil)
		return nil
	})
	So(err, ShouldBeNil)
}

func codeContextCanceled(codes *SampleCodes) {
	ctx := util.CanceledCtx()

	code := codes.SampleCode()
	num, err := code.Number(ctx)
	So(err, ShouldNotBeNil)
	So(num, ShouldEqual, 0)

	err = codes.DeleteSampleCode(ctx)
	So(err, ShouldNotBeNil)

	c, err := code.Code16(ctx)
	So(err, ShouldNotBeNil)
	So(c, ShouldBeEmpty)

	c, err = code.Code(ctx)
	So(err, ShouldNotBeNil)
	So(c, ShouldBeEmpty)

	c, err = code.Code64(ctx)
	So(err, ShouldNotBeNil)
	So(c, ShouldBeEmpty)

}

func TestConcurrentCode(t *testing.T) {
	ctx := context.Background()
	gdb, _ := NewSampleGlobalDB(ctx)
	defer gdb.Close()
	codes := gdb.Codes()
	err := codes.DeleteSampleCode(ctx)
	So(err, ShouldBeNil)

	var concurrent = 3
	var wg sync.WaitGroup
	wg.Add(concurrent)
	createCode := func() {
		db, _ := NewSampleGlobalDB(ctx)
		defer db.Close()
		code := codes.SampleCode()

		for i := 0; i < 5; i++ {

			num, err := code.Number(ctx)
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
		go createCode()
	}
	wg.Wait()
	code := codes.SampleCode()
	num, err := code.Number(ctx)
	if err != nil {
		t.Errorf("err should be nil, got %v", err)
		return
	}
	if num == 0 {
		t.Errorf("serial = %d; need more than 0", num)
	}
}
