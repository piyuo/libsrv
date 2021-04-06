package gdb

import (
	"context"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/piyuo/libsrv/src/db"
	identifier "github.com/piyuo/libsrv/src/identifier"
)

//TestConcurrent will turn root sample into 15 different sample and make sure every sample has different value
//
func TestConcurrent(t *testing.T) {
	t.Parallel()
	rand.Seed(time.Now().UTC().UnixNano())
	ctx := context.Background()
	client := sampleClient()
	rand := identifier.RandomString(8)
	counterName := "test-concurrent-counter-" + rand
	counter := client.Counter(counterName, 1000)
	defer counter.Delete(ctx)

	coderName := "test-concurrent-coder-" + rand
	coder := client.Coder(coderName, 1000)
	defer coder.Delete(ctx)

	//cleanup
	defer client.Query(&Sample{}).Where("Tag", "==", rand).Delete(ctx, 100)

	//	begin concurrent run
	var concurrent = 3
	var wg sync.WaitGroup
	wg.Add(concurrent)
	sampling := func() {
		for i := 0; i < 5; i++ {
			err := client.Transaction(ctx, func(ctx context.Context, tx db.Transaction) error {
				// read count first to avoid read after write error
				counter := client.Counter(counterName, 1000)
				coder := client.Coder(coderName, 1000)
				num, err := coder.NumberRX(ctx, tx)
				if err != nil {
					t.Errorf("err should be nil, got %v", err)
					return err
				}

				if err := counter.IncrementRX(ctx, tx); err != nil {
					t.Errorf("err should be nil, got %v", err)
					return err
				}

				code := identifier.SerialID32(uint32(num))
				sample := &Sample{
					Name:  code,
					Value: int(num),
					Tag:   rand,
				}
				sample.SetID(code)

				if err := client.Set(ctx, sample); err != nil {
					t.Errorf("err should be nil, got %v", err)
					return err
				}

				if err := counter.IncrementWX(ctx, tx, 1); err != nil {
					t.Errorf("err should be nil, got %v", err)
					return err
				}

				return coder.NumberWX(ctx, tx)
			})
			if err != nil {
				t.Errorf("commit transaction %v", err)
			}
		}
		wg.Done()
	}

	// create go routing to do counting
	for i := 0; i < concurrent; i++ {
		go sampling()
	}
	wg.Wait()

	// finish concurrent run
	count, err := counter.CountAll(ctx)
	if err != nil {
		t.Errorf("err should be nil, got %v", err)
	}
	if count != 15 {
		t.Errorf("count = %v; want 15", count)
	}
}
