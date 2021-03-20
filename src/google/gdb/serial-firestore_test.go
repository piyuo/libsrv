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

func TestSerial(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	client := sampleClient()
	name := "test-serial-" + identifier.RandomString(6)
	serial := client.Serial(name)
	defer serial.Clear(ctx)

	var firstSerial int64
	err := client.Transaction(ctx, func(ctx context.Context, tx db.Transaction) error {
		num, err := serial.NumberRX(ctx, tx)
		assert.Nil(err)
		firstSerial = num
		return serial.NumberWX(ctx, tx)
	})
	assert.Nil(err)
	assert.True(firstSerial >= 0)

	var failSerial int64
	err = client.Transaction(ctx, func(ctx context.Context, tx db.Transaction) error {
		num, err := serial.NumberRX(ctx, tx)
		assert.Nil(err)
		failSerial = num
		err = serial.NumberWX(ctx, tx)
		assert.Nil(err)
		return errors.New("fail")
	})
	assert.NotNil(err)

	var secondSerial int64
	err = client.Transaction(ctx, func(ctx context.Context, tx db.Transaction) error {
		num, err := serial.NumberRX(ctx, tx)
		assert.Nil(err)
		secondSerial = num
		return serial.NumberWX(ctx, tx)
	})
	assert.Nil(err)
	assert.True(secondSerial > firstSerial)
	assert.True(secondSerial >= failSerial)

	// reset serial
	cleared, err := serial.Clear(ctx)
	assert.Nil(err)
	assert.True(cleared)
}

func TestSerialInCanceledCtx(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	client := sampleClient()
	serial := client.Serial("serialInCancelCtx")

	ctxCanceled := util.CanceledCtx()
	cleared, err := serial.Clear(ctxCanceled)
	assert.NotNil(err)
	assert.False(cleared)
}

func TestSerialConcurrent(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	rand.Seed(time.Now().UnixNano())
	client := sampleClient()
	name := "test-serial-concurrent-" + identifier.RandomString(6)
	serial := client.Serial(name)
	defer serial.Clear(ctx)

	var concurrent = 3
	var wg sync.WaitGroup
	wg.Add(concurrent)
	createserial := func() {
		//		time.Sleep(time.Duration(rand.Intn(2)) * time.Second)
		for i := 0; i < 3; i++ {
			serial := client.Serial(name)

			err := client.Transaction(ctx, func(ctx context.Context, tx db.Transaction) error {
				_, err := serial.NumberRX(ctx, tx)
				if err != nil {
					t.Errorf("rx err should be nil, got %v", err)
				}
				//fmt.Printf("num:%v\n", num)
				return serial.NumberWX(ctx, tx)
			})
			if err != nil {
				t.Errorf("wx err should be nil, got %v", err)
			}
			// serial update need to be low frequency
			time.Sleep(time.Duration(rand.Intn(1)) * time.Second)
		}
		wg.Done()
	}

	//create go routing to do counting
	for i := 0; i < concurrent; i++ {
		go createserial()
	}
	wg.Wait()

	err := client.Transaction(ctx, func(ctx context.Context, tx db.Transaction) error {
		num, err := serial.NumberRX(ctx, tx)
		if err != nil {
			t.Errorf("rx err should be nil, got %v", err)
		}
		if num != 10 {
			t.Errorf("serial = %d; want 10", num)
		}
		return serial.NumberWX(ctx, tx)
	})
	if err != nil {
		t.Errorf("tx err should be nil, got %v", err)
		return
	}
}
