package data

import (
	"context"
	"math/rand"
	"sync"
	"testing"
	"time"

	util "github.com/piyuo/libsrv/util"
	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCoder(t *testing.T) {
	Convey("check coder function", t, func() {
		ctx := context.Background()
		dbG, dbR := createSampleDB()
		defer removeSampleDB(dbG, dbR)
		codesG, codesR := createSampleCoders(dbG, dbR)

		coderMustUseWithInTransacton(codesG)
		coderMustReadBeforeWrite(ctx, dbG, codesG)

		coderInFailTransaction(ctx, dbG, codesG)
		coderInFailTransaction(ctx, dbR, codesR)

		coderInTransaction(ctx, dbG, codesG)
		coderInTransaction(ctx, dbR, codesR)

		coderReset(ctx, dbG, codesG)
		coderReset(ctx, dbR, codesR)

		testCoderInCanceledCtx(ctx, dbR, codesG)
		testCoderInCanceledCtx(ctx, dbR, codesR)
	})
}

func testCoderInCanceledCtx(ctx context.Context, db SampleDB, coders *SampleCoders) {
	coder := coders.SampleCoder()
	So(coder, ShouldNotBeNil)

	ctxCanceled := util.CanceledCtx()
	err := coder.Clear(ctxCanceled)
	So(err, ShouldNotBeNil)
}

func coderMustUseWithInTransacton(codes *SampleCoders) {
	coder := codes.SampleCoder()
	ctx := context.Background()
	num, err := coder.NumberRX(ctx)
	So(err, ShouldNotBeNil)
	So(num, ShouldEqual, 0)
	err = coder.NumberWX(ctx)
	So(err, ShouldNotBeNil)

	code, err := coder.CodeRX(ctx)
	So(err, ShouldNotBeNil)
	So(code, ShouldBeEmpty)
	err = coder.CodeWX(ctx)
	So(err, ShouldNotBeNil)

	code, err = coder.Code16RX(ctx)
	So(err, ShouldNotBeNil)
	So(code, ShouldBeEmpty)
	err = coder.Code16WX(ctx)
	So(err, ShouldNotBeNil)

	code, err = coder.Code64RX(ctx)
	So(err, ShouldNotBeNil)
	So(code, ShouldBeEmpty)
	err = coder.Code64WX(ctx)
	So(err, ShouldNotBeNil)
}

func coderMustReadBeforeWrite(ctx context.Context, db *SampleGlobalDB, codes *SampleCoders) {
	db.Transaction(ctx, func(ctx context.Context) error {
		coder := codes.SampleCoder()
		err := coder.NumberWX(ctx)
		So(err, ShouldNotBeNil)
		err = coder.CodeWX(ctx)
		So(err, ShouldNotBeNil)
		err = coder.Code16WX(ctx)
		So(err, ShouldNotBeNil)
		err = coder.Code64WX(ctx)
		So(err, ShouldNotBeNil)
		return nil
	})
}

func coderInFailTransaction(ctx context.Context, db SampleDB, codes *SampleCoders) {

	coder := codes.SampleCoder()
	err := coder.Clear(ctx)
	So(err, ShouldBeNil)
	defer coder.Clear(ctx)

	shardsCount, err := coder.ShardsCount(ctx)
	So(err, ShouldBeNil)
	So(shardsCount, ShouldEqual, 0)

	// success
	err = db.Transaction(ctx, func(ctx context.Context) error {
		coder := codes.SampleCoder()
		num, err := coder.NumberRX(ctx)
		So(err, ShouldBeNil)
		So(num, ShouldBeGreaterThanOrEqualTo, 10)
		return coder.NumberWX(ctx)
	})
	So(err, ShouldBeNil)

	shardsCount, err = coder.ShardsCount(ctx)
	So(err, ShouldBeNil)
	So(shardsCount, ShouldEqual, 1)

	// fail
	err = db.Transaction(ctx, func(ctx context.Context) error {
		coder := codes.SampleCoder()
		num, err := coder.NumberRX(ctx)
		So(err, ShouldBeNil)
		So(num, ShouldBeGreaterThanOrEqualTo, 10)
		err = coder.NumberWX(ctx)
		So(err, ShouldBeNil)
		return errors.New("make transation fail")
	})
	So(err, ShouldNotBeNil)

	shardsCount, err = coder.ShardsCount(ctx)
	So(err, ShouldBeNil)
	So(shardsCount, ShouldEqual, 1)
}

func coderInTransaction(ctx context.Context, db SampleDB, codes *SampleCoders) {
	coder := codes.SampleCoder()
	err := coder.Clear(ctx)
	So(err, ShouldBeNil)
	defer coder.Clear(ctx)

	var num1 int64
	var num2 int64
	err = db.Transaction(ctx, func(ctx context.Context) error {
		coder := codes.SampleCoder()
		num1, err = coder.NumberRX(ctx)
		So(err, ShouldBeNil)
		So(num1, ShouldBeGreaterThanOrEqualTo, 10)
		return coder.NumberWX(ctx)
	})
	So(err, ShouldBeNil)

	// call second time
	err = db.Transaction(ctx, func(ctx context.Context) error {
		coder := codes.SampleCoder()
		num2, err = coder.NumberRX(ctx)
		So(err, ShouldBeNil)
		So(num2, ShouldBeGreaterThanOrEqualTo, 10)
		return coder.NumberWX(ctx)
	})
	So(err, ShouldBeNil)
	So(num1, ShouldNotEqual, num2)

	var code1 string
	var code2 string
	err = db.Transaction(ctx, func(ctx context.Context) error {
		coder := codes.SampleCoder()
		code1, err = coder.CodeRX(ctx)
		So(err, ShouldBeNil)
		So(code1, ShouldNotBeEmpty)
		return coder.CodeWX(ctx)
	})
	So(err, ShouldBeNil)
	err = db.Transaction(ctx, func(ctx context.Context) error {
		coder := codes.SampleCoder()
		code2, err = coder.CodeRX(ctx)
		So(err, ShouldBeNil)
		So(code2, ShouldNotBeEmpty)
		return coder.CodeWX(ctx)
	})
	So(err, ShouldBeNil)
	So(code1, ShouldNotEqual, code2)

	var code161 string
	var code162 string
	err = db.Transaction(ctx, func(ctx context.Context) error {
		coder := codes.SampleCoder()
		code161, err = coder.Code16RX(ctx)
		So(err, ShouldBeNil)
		So(code161, ShouldNotBeEmpty)
		return coder.Code16WX(ctx)
	})
	So(err, ShouldBeNil)
	err = db.Transaction(ctx, func(ctx context.Context) error {
		coder := codes.SampleCoder()
		code162, err = coder.Code16RX(ctx)
		So(err, ShouldBeNil)
		So(code162, ShouldNotBeEmpty)
		return coder.Code16WX(ctx)
	})
	So(err, ShouldBeNil)
	So(code161, ShouldNotEqual, code162)

	var code641 string
	var code642 string
	err = db.Transaction(ctx, func(ctx context.Context) error {
		coder := codes.SampleCoder()
		code641, err = coder.Code64RX(ctx)
		So(err, ShouldBeNil)
		So(code641, ShouldNotBeEmpty)
		return coder.Code16WX(ctx)
	})
	So(err, ShouldBeNil)
	err = db.Transaction(ctx, func(ctx context.Context) error {
		coder := codes.SampleCoder()
		code642, err = coder.Code64RX(ctx)
		So(err, ShouldBeNil)
		So(code642, ShouldNotBeEmpty)
		return coder.Code64WX(ctx)
	})
	So(err, ShouldBeNil)
	So(code641, ShouldNotEqual, code642)
}

func TestConcurrentCoder(t *testing.T) {
	rand.Seed(time.Now().UTC().UnixNano())
	ctx := context.Background()

	gdb, _ := NewSampleGlobalDB(ctx)
	defer gdb.Close()
	coders := gdb.Coders()
	coder := coders.SampleCoder()
	err := coder.Clear(ctx)
	defer coder.Clear(ctx)

	result := make(map[int64]int64)
	resultMutex := sync.RWMutex{}

	var concurrent = 3
	var wg sync.WaitGroup
	wg.Add(concurrent)
	createCode := func() {
		db, _ := NewSampleGlobalDB(ctx)
		defer db.Close()

		for i := 0; i < 3; i++ {
			err = db.Transaction(ctx, func(ctx context.Context) error {
				coders := db.Coders()
				coder := coders.SampleCoder()
				num, err := coder.NumberRX(ctx)
				defer coder.NumberWX(ctx)
				if err != nil {
					t.Errorf("err should be nil, got %v", err)
				}
				resultMutex.Lock()
				// this may print more than 9 time, cause transaction may rerun
				//fmt.Printf("num:%v\n", num)
				result[num] = num
				resultMutex.Unlock()
				return nil
			})
			if err != nil {
				t.Errorf("err should be nil, got %v", err)
			}
		}
		wg.Done()
	}
	//create go routing to do counting
	for i := 0; i < concurrent; i++ {
		go createCode()
	}
	wg.Wait()
	resultLen := len(result)
	if resultLen != 9 {
		t.Errorf("result = %d; need 9", resultLen)
	}
}

func coderReset(ctx context.Context, db SampleDB, codes *SampleCoders) {
	coder := codes.SampleCoder()
	err := coder.Clear(ctx)
	So(err, ShouldBeNil)
	defer coder.Clear(ctx)

	var num1 int64
	err = db.Transaction(ctx, func(ctx context.Context) error {
		coder := codes.SampleCoder()
		num1, err = coder.NumberRX(ctx)
		So(err, ShouldBeNil)
		So(num1, ShouldBeGreaterThanOrEqualTo, 10)
		return coder.NumberWX(ctx)
	})
	So(err, ShouldBeNil)

	shardsCount, err := coder.ShardsCount(ctx)
	So(err, ShouldBeNil)
	So(shardsCount, ShouldEqual, 1)

	// reset
	coder = codes.SampleCoder()
	coder.Clear(ctx)

	shardsCount, err = coder.ShardsCount(ctx)
	So(err, ShouldBeNil)
	So(shardsCount, ShouldEqual, 0)

	err = db.Transaction(ctx, func(ctx context.Context) error {
		coder := codes.SampleCoder()
		num1, err = coder.NumberRX(ctx)
		So(err, ShouldBeNil)
		So(num1, ShouldBeGreaterThanOrEqualTo, 10)
		return coder.NumberWX(ctx)
	})
	So(err, ShouldBeNil)

	shardsCount, err = coder.ShardsCount(ctx)
	So(err, ShouldBeNil)
	So(shardsCount, ShouldEqual, 1)

	// reset in transaction
	coder = codes.SampleCoder()
	err = db.Transaction(ctx, func(ctx context.Context) error {
		return coder.Clear(ctx)
	})

	shardsCount, err = coder.ShardsCount(ctx)
	So(err, ShouldBeNil)
	So(shardsCount, ShouldEqual, 0)
}
