package gstore

import (
	"context"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/piyuo/libsrv/src/util"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestCoderInCanceledCtx(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	g, err := NewSampleGlobalDB(ctx)
	assert.Nil(err)
	defer g.Close()
	coders := g.Coders()

	coder := coders.SampleCoder()
	assert.NotNil(coder)

	ctxCanceled := util.CanceledCtx()
	err = coder.Clear(ctxCanceled)
	assert.NotNil(err)
}

func TestMustUseWithInTransacton(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	g, err := NewSampleGlobalDB(ctx)
	assert.Nil(err)
	defer g.Close()

	coders := g.Coders()
	coder := coders.SampleCoder()
	num, err := coder.NumberRX(ctx)
	assert.NotNil(err)
	assert.Equal(int64(0), num)
	err = coder.NumberWX(ctx)
	assert.NotNil(err)

	code, err := coder.CodeRX(ctx)
	assert.NotNil(err)
	assert.Empty(code)
	err = coder.CodeWX(ctx)
	assert.NotNil(err)

	code, err = coder.Code16RX(ctx)
	assert.NotNil(err)
	assert.Empty(code)
	err = coder.Code16WX(ctx)
	assert.NotNil(err)

	code, err = coder.Code64RX(ctx)
	assert.NotNil(err)
	assert.Empty(code)
	err = coder.Code64WX(ctx)
	assert.NotNil(err)
}

func TestMustReadBeforeWrite(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	g, err := NewSampleGlobalDB(ctx)
	assert.Nil(err)
	defer g.Close()
	coders := g.Coders()

	g.Transaction(ctx, func(ctx context.Context) error {
		coder := coders.SampleCoder()
		err := coder.NumberWX(ctx)
		assert.NotNil(err)
		err = coder.CodeWX(ctx)
		assert.NotNil(err)
		err = coder.Code16WX(ctx)
		assert.NotNil(err)
		err = coder.Code64WX(ctx)
		assert.NotNil(err)
		return nil
	})
}

func TestInFailTransaction(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	g, err := NewSampleGlobalDB(ctx)
	assert.Nil(err)
	defer g.Close()
	coders := g.Coders()

	coder := coders.SampleCoder()
	err = coder.Clear(ctx)
	assert.Nil(err)
	defer coder.Clear(ctx)

	shardsCount, err := coder.ShardsCount(ctx)
	assert.Nil(err)
	assert.Equal(0, shardsCount)

	// success
	err = g.Transaction(ctx, func(ctx context.Context) error {
		coder := coders.SampleCoder()
		num, err := coder.NumberRX(ctx)
		assert.Nil(err)
		assert.GreaterOrEqual(num, int64(10))
		return coder.NumberWX(ctx)
	})
	assert.Nil(err)

	shardsCount, err = coder.ShardsCount(ctx)
	assert.Nil(err)
	assert.Equal(1, shardsCount)

	// fail
	err = g.Transaction(ctx, func(ctx context.Context) error {
		coder := coders.SampleCoder()
		num, err := coder.NumberRX(ctx)
		assert.Nil(err)
		assert.GreaterOrEqual(num, int64(10))
		err = coder.NumberWX(ctx)
		assert.Nil(err)
		return errors.New("make transation fail")
	})
	assert.NotNil(err)

	shardsCount, err = coder.ShardsCount(ctx)
	assert.Nil(err)
	assert.Equal(1, shardsCount)

}

func TestInTransaction(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	g, err := NewSampleGlobalDB(ctx)
	assert.Nil(err)
	defer g.Close()
	coders := g.Coders()

	coder := coders.SampleCoder()
	err = coder.Clear(ctx)
	assert.Nil(err)
	defer coder.Clear(ctx)

	var num1 int64
	var num2 int64
	err = g.Transaction(ctx, func(ctx context.Context) error {
		coder := coders.SampleCoder()
		num1, err = coder.NumberRX(ctx)
		assert.Nil(err)
		assert.GreaterOrEqual(num1, int64(10))
		return coder.NumberWX(ctx)
	})
	assert.Nil(err)

	// call second time
	err = g.Transaction(ctx, func(ctx context.Context) error {
		coder := coders.SampleCoder()
		num2, err = coder.NumberRX(ctx)
		assert.Nil(err)
		assert.GreaterOrEqual(num2, int64(10))
		return coder.NumberWX(ctx)
	})
	assert.Nil(err)
	assert.NotEqual(num1, num2)

	var code1 string
	var code2 string
	err = g.Transaction(ctx, func(ctx context.Context) error {
		coder := coders.SampleCoder()
		code1, err = coder.CodeRX(ctx)
		assert.Nil(err)
		assert.NotEmpty(code1)
		return coder.CodeWX(ctx)
	})
	assert.Nil(err)
	err = g.Transaction(ctx, func(ctx context.Context) error {
		coder := coders.SampleCoder()
		code2, err = coder.CodeRX(ctx)
		assert.Nil(err)
		assert.NotEmpty(code2)
		return coder.CodeWX(ctx)
	})
	assert.Nil(err)
	assert.NotEqual(code1, code2)

	var code161 string
	var code162 string
	err = g.Transaction(ctx, func(ctx context.Context) error {
		coder := coders.SampleCoder()
		code161, err = coder.Code16RX(ctx)
		assert.Nil(err)
		assert.NotEmpty(code161)
		return coder.Code16WX(ctx)
	})
	assert.Nil(err)
	err = g.Transaction(ctx, func(ctx context.Context) error {
		coder := coders.SampleCoder()
		code162, err = coder.Code16RX(ctx)
		assert.Nil(err)
		assert.NotEmpty(code162)
		return coder.Code16WX(ctx)
	})
	assert.Nil(err)
	assert.NotEqual(code161, code162)

	var code641 string
	var code642 string
	err = g.Transaction(ctx, func(ctx context.Context) error {
		coder := coders.SampleCoder()
		code641, err = coder.Code64RX(ctx)
		assert.Nil(err)
		assert.NotEmpty(code641)
		return coder.Code16WX(ctx)
	})
	assert.Nil(err)
	err = g.Transaction(ctx, func(ctx context.Context) error {
		coder := coders.SampleCoder()
		code642, err = coder.Code64RX(ctx)
		assert.Nil(err)
		assert.NotEmpty(code642)
		return coder.Code64WX(ctx)
	})
	assert.Nil(err)
	assert.NotEqual(code641, code642)
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

func TestCoderReset(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	g, err := NewSampleGlobalDB(ctx)
	assert.Nil(err)
	defer g.Close()
	coders := g.Coders()

	coder := coders.SampleCoder()
	err = coder.Clear(ctx)
	assert.Nil(err)
	defer coder.Clear(ctx)

	var num1 int64
	err = g.Transaction(ctx, func(ctx context.Context) error {
		coder := coders.SampleCoder()
		num1, err = coder.NumberRX(ctx)
		assert.Nil(err)
		assert.GreaterOrEqual(num1, int64(10))
		return coder.NumberWX(ctx)
	})
	assert.Nil(err)

	shardsCount, err := coder.ShardsCount(ctx)
	assert.Nil(err)
	assert.Equal(1, shardsCount)

	// reset
	coder = coders.SampleCoder()
	coder.Clear(ctx)

	shardsCount, err = coder.ShardsCount(ctx)
	assert.Nil(err)
	assert.Equal(0, shardsCount)

	err = g.Transaction(ctx, func(ctx context.Context) error {
		coder := coders.SampleCoder()
		num1, err = coder.NumberRX(ctx)
		assert.Nil(err)
		assert.GreaterOrEqual(num1, int64(10))
		return coder.NumberWX(ctx)
	})
	assert.Nil(err)

	shardsCount, err = coder.ShardsCount(ctx)
	assert.Nil(err)
	assert.Equal(1, shardsCount)

	// reset in transaction
	coder = coders.SampleCoder()
	err = g.Transaction(ctx, func(ctx context.Context) error {
		return coder.Clear(ctx)
	})

	shardsCount, err = coder.ShardsCount(ctx)
	assert.Nil(err)
	assert.Equal(0, shardsCount)

}
