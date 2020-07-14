package data

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/piyuo/libsrv/util"

	"github.com/pkg/errors"
)

//TestConcurrentDB will turn root sample into 15 different sample and make sure every sample has different value
//
func TestConcurrentDB(t *testing.T) {
	ctx := context.Background()

	dbG, dbR := createSampleDB()
	defer removeSampleDB(dbG, dbR)
	tableG, tableR := createSampleTable(dbG, dbR)
	defer removeSampleTable(tableG, tableR)
	countersG, countersR := createSampleCounters(dbG, dbR)
	defer removeSampleCounters(countersG, countersR)
	codesG, codesR := createSampleCodes(dbG, dbR)
	defer removeSampleCodes(codesG, codesR)

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

	counter := countersG.SampleCounter()
	if err := counter.CreateShards(ctx); err != nil {
		t.Errorf("err should be nil, got %v", err)
		return
	}

	//code := codesG.SampleCode()
	//if err := code.CreateShards(ctx); err != nil {
	//	t.Errorf("err should be nil, got %v", err)
	//	return
	//}

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
				sRootRef, err := stable.Find(ctx, "Name", "==", "root")
				if err != nil {
					t.Errorf("err should be nil, got %v", err)
					return errors.New("failed to find")
				}
				sRoot := sRootRef.(*Sample)

				// read count first to avoid read after write error
				counters := sdb.Counters()
				counter := counters.SampleCounter()
				codes := sdb.Codes()
				code := codes.SampleCode()

				num, err2 := code.Number(ctx)
				if err2 != nil {
					t.Errorf("err should be nil, got %v", err)
					return errors.New("failed to get code")
				}
				fmt.Printf("sampling:%v\n", num)

				if err := counter.Increment(ctx, 1); err != nil {
					t.Errorf("err should be nil, got %v", err)
					return errors.New("failed to increment")
				}

				sSample := &Sample{
					Name:  util.SerialID32(uint32(num)),
					Value: int(num),
				}
				//sSample.SetID(code)

				if err := stable.Set(ctx, sSample); err != nil {
					t.Errorf("err should be nil, got %v", err)
					return errors.New("failed to create sample")
				}

				sRoot.Value--
				if err := stable.Set(ctx, sRoot); err != nil {
					t.Errorf("err should be nil, got %v", err)
					return errors.New("failed to update root sample")
				}
				return nil
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
