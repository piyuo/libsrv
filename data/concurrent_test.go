package data

import (
	"context"
	"fmt"
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
	ctx := context.Background()

	dbG, dbR := createSampleDB()
	defer removeSampleDB(dbG, dbR)
	tableG, tableR := createSampleTable(dbG, dbR)
	removeSampleTable(tableG, tableR)
	defer removeSampleTable(tableG, tableR)
	countersG, countersR := createSampleCounters(dbG, dbR)
	removeSampleCounters(countersG, countersR)
	defer removeSampleCounters(countersG, countersR)
	codersG, codersR := createSampleCoders(dbG, dbR)
	removeSampleCoders(codersG, codersR)
	defer removeSampleCoders(codersG, codersR)

	//init test data
	table := tableG
	root := &Sample{
		Name:  "root",
		Value: 15,
	}
	root.SetID("root")

	if err := table.Set(ctx, root); err != nil {
		t.Errorf("err should be nil, got %v", err)
	}

	//	begin concurrent run
	var concurrent = 3
	var wg sync.WaitGroup
	wg.Add(concurrent)
	sampling := func() {
		rand.Seed(time.Now().UnixNano())
		sdb, err := NewSampleGlobalDB(ctx)
		if err != nil {
			t.Errorf("err should be nil, got %v", err)
		}
		defer sdb.Close()

		for i := 0; i < 5; i++ {
			errTx := sdb.Transaction(ctx, func(ctx context.Context) error {
				stable := sdb.SampleTable()
				sRootRef, err := stable.Find(ctx, "Name", "==", "root")
				if err != nil {
					t.Errorf("err should be nil, got %v", err)
					return errors.New("failed to find")
				}
				sRoot := sRootRef.(*Sample)

				// read count first to avoid read after write error
				counter := sdb.Counters().SampleCounter500()
				coder := sdb.Coders().SampleCoder500()

				num, err2 := coder.NumberRX()
				if err2 != nil {
					t.Errorf("err should be nil, got %v", err2)
					return errors.Wrap(err, "failed to get code")
				}
				fmt.Printf("sampling:%v\n", num)

				if err := counter.IncrementRX(1); err != nil {
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

				sRoot.Value--
				if err := stable.Set(ctx, sRoot); err != nil {
					t.Errorf("err should be nil, got %v", err)
					return errors.Wrap(err, "failed to update root sample")
				}

				if err := counter.IncrementWX(); err != nil {
					t.Errorf("err should be nil, got %v", err)
					return errors.Wrap(err, "failed to IncrementWX")
				}

				return coder.NumberWX()
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
	rootRef, err := table.Find(ctx, "Name", "==", "root")
	if err != nil {
		t.Errorf("err should be nil, got %v", err)
	}
	root = rootRef.(*Sample)
	if root.Value != 0 {
		t.Errorf("serial = %d; want 0", root.Value)
	}
}
