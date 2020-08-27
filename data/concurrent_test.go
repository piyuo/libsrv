package data

import (
	"context"
	"math/rand"
	"sync"
	"testing"
	"time"

	identifier "github.com/piyuo/libsrv/identifier"

	"github.com/pkg/errors"
)

//TestConcurrentDB will turn root sample into 15 different sample and make sure every sample has different value
//
func TestConcurrentDB(t *testing.T) {
	rand.Seed(time.Now().UTC().UnixNano())
	ctx := context.Background()

	dbG, dbR := createSampleDB()
	defer removeSampleDB(dbG, dbR)
	tableG, tableR := createSampleTable(dbG, dbR)

	removeSampleTable(tableG, tableR)
	defer removeSampleTable(tableG, tableR)
	countersG, countersR := createSampleCounters(dbG, dbR)
	codersG, codersR := createSampleCoders(dbG, dbR)

	counterG := countersG.SampleCounter()
	counterG.Clear(ctx)
	defer counterG.Clear(ctx)

	counterR := countersR.SampleCounter()
	counterR.Clear(ctx)
	defer counterR.Clear(ctx)

	coderG := codersG.SampleCoder()
	coderG.Clear(ctx)
	defer coderG.Clear(ctx)

	coderR := codersR.SampleCoder()
	coderR.Clear(ctx)
	defer coderR.Clear(ctx)

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
				stable := sdb.SampleTable()

				// read count first to avoid read after write error
				counter := sdb.Counters().SampleCounter1000()
				coder := sdb.Coders().SampleCoder1000()

				num, err2 := coder.NumberRX(ctx)
				if err2 != nil {
					t.Errorf("err should be nil, got %v", err2)
					return errors.Wrap(err, "failed to get code")
				}

				if err := counter.IncrementRX(ctx, 1); err != nil {
					t.Errorf("err should be nil, got %v", err)
					return errors.Wrap(err, "failed to IncrementRX")
				}

				code := identifier.SerialID32(uint32(num))
				sSample := &Sample{
					Name:  code,
					Value: int(num),
				}
				sSample.SetID(code)

				if err := stable.Set(ctx, sSample); err != nil {
					t.Errorf("err should be nil, got %v", err)
					return errors.Wrap(err, "failed to create sample")
				}

				if err := counter.IncrementWX(ctx); err != nil {
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

	count, err := counterG.CountAll(ctx)
	if err != nil {
		t.Errorf("err should be nil, got %v", err)
	}
	if count != 15 {
		t.Errorf("count = %v; want 15", count)
	}
}
