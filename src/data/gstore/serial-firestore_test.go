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

func TestSerialInCanceledCtx(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	g, err := NewSampleGlobalDB(ctx)
	assert.Nil(err)
	defer g.Close()
	serial := g.Serials().SampleSerial()

	ctxCanceled := util.CanceledCtx()
	err = serial.Clear(ctxCanceled)
	assert.NotNil(err)
}

func TestSerialMustUseWithInTransacton(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	g, err := NewSampleGlobalDB(ctx)
	assert.Nil(err)
	defer g.Close()
	serial := g.Serials().SampleSerial()

	num, err := serial.NumberRX()
	assert.NotNil(err)
	assert.Equal(int64(0), num)
	err = serial.NumberWX()
	assert.NotNil(err)
}

func TestSerialInTransactionTest(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	g, err := NewSampleGlobalDB(ctx)
	assert.Nil(err)
	defer g.Close()
	serial := g.Serials().SampleSerial()
	defer serial.Clear(ctx)

	err = g.Transaction(ctx, func(ctx context.Context) error {
		num, err := serial.NumberRX()
		assert.Nil(err)
		assert.Equal(int64(1), num)
		return serial.NumberWX()
	})
	assert.Nil(err)

	err = g.Transaction(ctx, func(ctx context.Context) error {
		num, err := serial.NumberRX()
		assert.Nil(err)
		assert.Equal(int64(2), num)
		err = serial.NumberWX()
		assert.Nil(err)
		return errors.New("make fail transaction")
	})
	assert.NotNil(err)

	err = g.Transaction(ctx, func(ctx context.Context) error {
		num, err := serial.NumberRX()
		assert.Nil(err)
		assert.Equal(int64(2), num)
		return serial.NumberWX()
	})
	assert.Nil(err)

	err = g.Transaction(ctx, func(ctx context.Context) error {
		num, err := serial.NumberRX()
		assert.Nil(err)
		assert.Equal(int64(3), num)
		return serial.NumberWX()
	})

	// reset serial
	err = serial.Clear(ctx)
	assert.Nil(err)

	err = g.Transaction(ctx, func(ctx context.Context) error {
		num, err := serial.NumberRX()
		assert.Nil(err)
		assert.Equal(int64(1), num)
		return serial.NumberWX()
	})

	// reset in transaction
	err = g.Transaction(ctx, func(ctx context.Context) error {
		return serial.Clear(ctx)
	})

	err = g.Transaction(ctx, func(ctx context.Context) error {
		num, err := serial.NumberRX()
		assert.Nil(err)
		assert.Equal(int64(1), num)
		return serial.NumberWX()
	})
	assert.Nil(err)
}

func TestConcurrentSerial(t *testing.T) {
	ctx := context.Background()
	rand.Seed(time.Now().UnixNano())
	g, err := NewSampleGlobalDB(ctx)
	defer g.Close()
	serial := g.Serials().SampleSerial()
	defer serial.Clear(ctx)

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
				_, err := serial.NumberRX()
				if err != nil {
					t.Errorf("err should be nil, got %v", err)
				}
				//fmt.Printf("num:%v\n", num)
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

	err = g.Transaction(ctx, func(ctx context.Context) error {
		serial := g.Serials().SampleSerial()
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
