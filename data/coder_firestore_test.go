package data

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCoder(t *testing.T) {
	Convey("check coder function", t, func() {
		ctx := context.Background()
		dbG, dbR := createSampleDB()
		defer removeSampleDB(dbG, dbR)
		codesG, codesR := createSampleCoders(dbG, dbR)
		defer removeSampleCoders(codesG, codesR)

		coderMustUseWithInTransacton(codesG)
		coderMustReadBeforeWrite(ctx, dbG, codesG)

		coderInFailTransaction(ctx, dbG, codesG)
		coderInFailTransaction(ctx, dbR, codesR)

		coderInTransaction(ctx, dbG, codesG)
		coderInTransaction(ctx, dbR, codesR)
	})
}

func coderMustUseWithInTransacton(codes *SampleCoders) {
	coder := codes.SampleCoder()

	num, err := coder.NumberRX()
	So(err, ShouldNotBeNil)
	So(num, ShouldEqual, 0)
	err = coder.NumberWX()
	So(err, ShouldNotBeNil)

	code, err := coder.CodeRX()
	So(err, ShouldNotBeNil)
	So(code, ShouldBeEmpty)
	err = coder.CodeWX()
	So(err, ShouldNotBeNil)

	code, err = coder.Code16RX()
	So(err, ShouldNotBeNil)
	So(code, ShouldBeEmpty)
	err = coder.Code16WX()
	So(err, ShouldNotBeNil)

	code, err = coder.Code64RX()
	So(err, ShouldNotBeNil)
	So(code, ShouldBeEmpty)
	err = coder.Code64WX()
	So(err, ShouldNotBeNil)
}

func coderMustReadBeforeWrite(ctx context.Context, db *SampleGlobalDB, codes *SampleCoders) {
	db.Transaction(ctx, func(ctx context.Context) error {
		coder := codes.SampleCoder()
		err := coder.NumberWX()
		So(err, ShouldNotBeNil)
		err = coder.CodeWX()
		So(err, ShouldNotBeNil)
		err = coder.Code16WX()
		So(err, ShouldNotBeNil)
		err = coder.Code64WX()
		So(err, ShouldNotBeNil)
		return nil
	})
}

func coderInFailTransaction(ctx context.Context, db SampleDB, codes *SampleCoders) {
	err := codes.DeleteSampleCode(ctx)
	So(err, ShouldBeNil)

	iCoder := codes.SampleCoder()
	coder := iCoder.(*CoderFirestore)
	docCount, shardsCount, err := coder.shardsInfo(ctx)
	So(err, ShouldBeNil)
	So(docCount, ShouldEqual, 0)
	So(shardsCount, ShouldEqual, 0)

	// success
	err = db.Transaction(ctx, func(ctx context.Context) error {
		coder := codes.SampleCoder()
		num, err := coder.NumberRX()
		So(err, ShouldBeNil)
		So(num, ShouldBeGreaterThanOrEqualTo, 10)
		return coder.NumberWX()
	})
	So(err, ShouldBeNil)

	docCount, shardsCount, err = coder.shardsInfo(ctx)
	So(err, ShouldBeNil)
	So(docCount, ShouldEqual, 1)
	So(shardsCount, ShouldEqual, 1)

	// fail
	err = db.Transaction(ctx, func(ctx context.Context) error {
		coder := codes.SampleCoder()
		num, err := coder.NumberRX()
		So(err, ShouldBeNil)
		So(num, ShouldBeGreaterThanOrEqualTo, 10)
		err = coder.NumberWX()
		So(err, ShouldBeNil)
		return errors.New("make transation fail")
	})
	So(err, ShouldNotBeNil)

	docCount, shardsCount, err = coder.shardsInfo(ctx)
	So(err, ShouldBeNil)
	So(docCount, ShouldEqual, 1)
	So(shardsCount, ShouldEqual, 1)
}

func coderInTransaction(ctx context.Context, db SampleDB, codes *SampleCoders) {
	err := codes.DeleteSampleCode(ctx)
	So(err, ShouldBeNil)

	var num1 int64
	var num2 int64
	err = db.Transaction(ctx, func(ctx context.Context) error {
		coder := codes.SampleCoder()
		num1, err = coder.NumberRX()
		So(err, ShouldBeNil)
		So(num1, ShouldBeGreaterThanOrEqualTo, 10)
		return coder.NumberWX()
	})
	So(err, ShouldBeNil)

	// call second time
	err = db.Transaction(ctx, func(ctx context.Context) error {
		coder := codes.SampleCoder()
		num2, err = coder.NumberRX()
		So(err, ShouldBeNil)
		So(num2, ShouldBeGreaterThanOrEqualTo, 10)
		return coder.NumberWX()
	})
	So(err, ShouldBeNil)
	So(num1, ShouldNotEqual, num2)

	var code1 string
	var code2 string
	err = db.Transaction(ctx, func(ctx context.Context) error {
		coder := codes.SampleCoder()
		code1, err = coder.CodeRX()
		So(err, ShouldBeNil)
		So(code1, ShouldNotBeEmpty)
		return coder.CodeWX()
	})
	So(err, ShouldBeNil)
	err = db.Transaction(ctx, func(ctx context.Context) error {
		coder := codes.SampleCoder()
		code2, err = coder.CodeRX()
		So(err, ShouldBeNil)
		So(code2, ShouldNotBeEmpty)
		return coder.CodeWX()
	})
	So(err, ShouldBeNil)
	So(code1, ShouldNotEqual, code2)

	var code161 string
	var code162 string
	err = db.Transaction(ctx, func(ctx context.Context) error {
		coder := codes.SampleCoder()
		code161, err = coder.Code16RX()
		So(err, ShouldBeNil)
		So(code161, ShouldNotBeEmpty)
		return coder.Code16WX()
	})
	So(err, ShouldBeNil)
	err = db.Transaction(ctx, func(ctx context.Context) error {
		coder := codes.SampleCoder()
		code162, err = coder.Code16RX()
		So(err, ShouldBeNil)
		So(code162, ShouldNotBeEmpty)
		return coder.Code16WX()
	})
	So(err, ShouldBeNil)
	So(code161, ShouldNotEqual, code162)

	var code641 string
	var code642 string
	err = db.Transaction(ctx, func(ctx context.Context) error {
		coder := codes.SampleCoder()
		code641, err = coder.Code64RX()
		So(err, ShouldBeNil)
		So(code641, ShouldNotBeEmpty)
		return coder.Code16WX()
	})
	So(err, ShouldBeNil)
	err = db.Transaction(ctx, func(ctx context.Context) error {
		coder := codes.SampleCoder()
		code642, err = coder.Code64RX()
		So(err, ShouldBeNil)
		So(code642, ShouldNotBeEmpty)
		return coder.Code64WX()
	})
	So(err, ShouldBeNil)
	So(code641, ShouldNotEqual, code642)

	err = db.Transaction(ctx, func(ctx context.Context) error {
		err := codes.DeleteSampleCode(ctx)
		So(err, ShouldBeNil)
		return nil
	})
	So(err, ShouldBeNil)
}

func TestConcurrentCoder(t *testing.T) {
	ctx := context.Background()
	gdb, _ := NewSampleGlobalDB(ctx)
	defer gdb.Close()
	coders := gdb.Coders()
	defer coders.DeleteSampleCode(ctx)
	err := coders.DeleteSampleCode(ctx)
	if err != nil {
		t.Errorf("err should be nil, got %v", err)
	}

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
				num, err := coder.NumberRX()
				defer coder.NumberWX()
				if err != nil {
					t.Errorf("err should be nil, got %v", err)
				}
				resultMutex.Lock()
				// this may print more than 9 time, cause transaction may rerun
				fmt.Printf("num:%v\n", num)
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
