package gdb

import (
	"context"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/piyuo/libsrv/src/db"
	"github.com/piyuo/libsrv/src/identifier"
	"github.com/piyuo/libsrv/src/util"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestCoderInCanceledCtx(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	client := sampleClient()
	coder := client.Coder("cancelCtx", 3)
	assert.NotNil(coder)
	ctxCanceled := util.CanceledCtx()
	err := coder.Delete(ctxCanceled)
	assert.NotNil(err)
}

func TestCoderMustReadBeforeWrite(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	ctx := context.Background()
	client := sampleClient()
	err := client.Transaction(ctx, func(ctx context.Context, tx db.Transaction) error {
		coder := client.Coder("sample", 3)
		err := coder.NumberWX(ctx, tx)
		assert.NotNil(err)
		err = coder.CodeWX(ctx, tx)
		assert.NotNil(err)
		err = coder.Code16WX(ctx, tx)
		assert.NotNil(err)
		err = coder.Code64WX(ctx, tx)
		assert.NotNil(err)
		return nil
	})
	assert.Nil(err)
}

func TestCoderNum(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	ctx := context.Background()
	client := sampleClient()

	name := "test-coder-num" + identifier.RandomString(8)
	coder := client.Coder(name, 1)
	// success
	var firstNum int64
	err := client.Transaction(ctx, func(ctx context.Context, tx db.Transaction) error {
		num, err := coder.NumberRX(ctx, tx)
		assert.Nil(err)
		assert.True(num > 0)
		firstNum = num

		err = coder.NumberWX(ctx, tx)
		assert.Nil(err)
		return nil
	})
	assert.Nil(err)

	shardsCount, err := coder.ShardsCount(ctx)
	assert.Nil(err)
	assert.Equal(1, shardsCount)

	var failNum int64
	// fail should not change number
	err = client.Transaction(ctx, func(ctx context.Context, tx db.Transaction) error {
		num, err := coder.NumberRX(ctx, tx)
		assert.Nil(err)
		failNum = num

		err = coder.NumberWX(ctx, tx)
		assert.Nil(err)
		return errors.New("fail")
	})
	assert.NotNil(err)

	var currentNum int64
	err = client.Transaction(ctx, func(ctx context.Context, tx db.Transaction) error {
		num, err := coder.NumberRX(ctx, tx)
		assert.Nil(err)
		currentNum = num

		err = coder.NumberWX(ctx, tx)
		assert.Nil(err)
		return nil
	})
	assert.Nil(err)
	assert.Equal(failNum, currentNum)
	assert.NotEqual(firstNum, currentNum)

	err = coder.Delete(ctx)
	assert.Nil(err)

	shardsCount, err = coder.ShardsCount(ctx)
	assert.Nil(err)
	assert.Equal(0, shardsCount)
}

func TestCoderCode(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	ctx := context.Background()
	client := sampleClient()
	name := "test-coder-code" + identifier.RandomString(8)
	coder := client.Coder(name, 1)
	defer coder.Delete(ctx)
	var firstCode string
	err := client.Transaction(ctx, func(ctx context.Context, tx db.Transaction) error {
		code, err := coder.CodeRX(ctx, tx)
		assert.Nil(err)
		assert.NotEmpty(code)
		firstCode = code

		err = coder.CodeWX(ctx, tx)
		assert.Nil(err)
		return nil
	})
	assert.Nil(err)

	var currentCode string
	err = client.Transaction(ctx, func(ctx context.Context, tx db.Transaction) error {
		code, err := coder.CodeRX(ctx, tx)
		assert.Nil(err)
		currentCode = code

		err = coder.CodeWX(ctx, tx)
		assert.Nil(err)
		return nil
	})
	assert.Nil(err)
	assert.NotEqual(firstCode, currentCode)
}

func TestCoderCode16(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	ctx := context.Background()
	client := sampleClient()

	name := "test-coder-code16" + identifier.RandomString(6)
	coder := client.Coder(name, 1)
	defer coder.Delete(ctx)
	var firstCode string
	err := client.Transaction(ctx, func(ctx context.Context, tx db.Transaction) error {
		code, err := coder.Code16RX(ctx, tx)
		assert.Nil(err)
		assert.NotEmpty(code)
		firstCode = code

		err = coder.Code16WX(ctx, tx)
		assert.Nil(err)
		return nil
	})
	assert.Nil(err)

	var currentCode string
	err = client.Transaction(ctx, func(ctx context.Context, tx db.Transaction) error {
		code, err := coder.Code16RX(ctx, tx)
		assert.Nil(err)
		currentCode = code

		err = coder.Code16WX(ctx, tx)
		assert.Nil(err)
		return nil
	})
	assert.Nil(err)
	assert.NotEqual(firstCode, currentCode)
}

func TestCoderCode64(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	ctx := context.Background()
	client := sampleClient()

	name := "test-coder-code64" + identifier.RandomString(6)
	coder := client.Coder(name, 1)
	defer coder.Delete(ctx)
	var firstCode string
	err := client.Transaction(ctx, func(ctx context.Context, tx db.Transaction) error {
		code, err := coder.Code64RX(ctx, tx)
		assert.Nil(err)
		assert.NotEmpty(code)
		firstCode = code

		err = coder.Code64WX(ctx, tx)
		assert.Nil(err)
		return nil
	})
	assert.Nil(err)

	var currentCode string
	err = client.Transaction(ctx, func(ctx context.Context, tx db.Transaction) error {
		code, err := coder.Code64RX(ctx, tx)
		assert.Nil(err)
		currentCode = code

		err = coder.Code64WX(ctx, tx)
		assert.Nil(err)
		return nil
	})
	assert.Nil(err)
	assert.NotEqual(firstCode, currentCode)
}

func TestConcurrentCoder(t *testing.T) {
	rand.Seed(time.Now().UTC().UnixNano())
	ctx := context.Background()
	client := sampleClient()
	name := "test-coder-concurrent" + identifier.RandomString(6)
	result := make(map[int64]int64)
	resultMutex := sync.RWMutex{}

	coder := client.Coder(name, 30)
	defer coder.Delete(ctx)

	var concurrent = 3
	var wg sync.WaitGroup
	wg.Add(concurrent)
	createCode := func() {
		for i := 0; i < 3; i++ {
			err := client.Transaction(ctx, func(ctx context.Context, tx db.Transaction) error {
				coder := client.Coder(name, 30)
				num, err := coder.NumberRX(ctx, tx)
				if err != nil {
					t.Errorf("rx err should be nil, got %v", err)
				}
				err = coder.NumberWX(ctx, tx)
				if err != nil {
					t.Errorf("wx err should be nil, got %v", err)
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
