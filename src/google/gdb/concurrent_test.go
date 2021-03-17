package gdb

import (
	"context"
	"math/rand"
	"sync"
	"testing"
	"time"

	identifier "github.com/piyuo/libsrv/src/identifier"

	"github.com/pkg/errors"
)

//TestConcurrentDB will turn root sample into 15 different sample and make sure every sample has different value
//
func TestConcurrentDB(t *testing.T) {
	rand.Seed(time.Now().UTC().UnixNano())
	ctx := context.Background()
	g, err := NewSampleGlobalDB(ctx)
	defer g.Close()

	counters := g.Counters()
	coders := g.Coders()

	counter := counters.SampleCounter()
	defer counter.Clear(ctx)

	coder := coders.SampleCoder()
	defer coder.Clear(ctx)

	//	begin concurrent run
	var concurrent = 3
	var wg sync.WaitGroup
	wg.Add(concurrent)
	sampling := func() {
		sdb, err := NewSampleGlobalDB(ctx)
		if err != nil {
			t.Errorf("err should be nil, got %v", err)
		}
		defer sdb.Close()

		for i := 0; i < 5; i++ {
			errTx := sdb.Transaction(ctx, func(ctx context.Context) error {

				// read count first to avoid read after write error
				counter := sdb.Counters().SampleCounter1000()
				coder := sdb.Coders().SampleCoder1000()

				num, err2 := coder.NumberRX(ctx)
				if err2 != nil {
					t.Errorf("err should be nil, got %v", err2)
					return errors.Wrap(err, "failed to get code")
				}

				if err := counter.IncrementRX(ctx); err != nil {
					t.Errorf("err should be nil, got %v", err)
					return errors.Wrap(err, "failed to IncrementRX")
				}

				code := identifier.SerialID32(uint32(num))
				sSample := &Sample{
					Name:  code,
					Value: int(num),
				}
				sSample.SetID(code)

				if err := sdb.SampleTable().Set(ctx, sSample); err != nil {
					t.Errorf("err should be nil, got %v", err)
					return errors.Wrap(err, "failed to create sample")
				}

				if err := counter.IncrementWX(ctx, 1); err != nil {
					t.Errorf("err should be nil, got %v", err)
					return errors.Wrap(err, "failed to IncrementWX")
				}

				return coder.NumberWX(ctx)
			})
			if errTx != nil {
				t.Errorf("failed to commit transaction %v", errTx)
			}
		}
		wg.Done()
	}

	//create go routing to do counting
	for i := 0; i < concurrent; i++ {
		go sampling()
	}
	wg.Wait()
	//finish concurrent run

	count, err := counter.CountAll(ctx)
	if err != nil {
		t.Errorf("err should be nil, got %v", err)
	}
	if count != 15 {
		t.Errorf("count = %v; want 15", count)
	}

	g.SampleTable().Clear(ctx)
}
