package data

import (
	"context"
	"math/rand"
	"strconv"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/piyuo/libsrv/util"
	"github.com/pkg/errors"
	"google.golang.org/api/iterator"
)

// CounterFirestore implement Counter
//
type CounterFirestore struct {
	Counter `firestore:"-"`

	MetaFirestore `firestore:"-"`

	callRX bool

	// pickedShard is a shard random picked
	//
	pickedShard string

	// keepDateHierarchy set to true mean keep count in date keepDateHierarchy
	//
	keepDateHierarchy bool

	shardAllExist bool

	shardYearExist bool

	shardMonthExist bool

	shardDayExist bool

	shardHourExist bool
}

// isShardExist return true if shard already exist
//
func (c *CounterFirestore) isShardExist(ctx context.Context, ref *firestore.DocumentRef) (bool, error) {
	snapshot, err := c.conn.tx.Get(ref)
	if snapshot != nil && !snapshot.Exists() {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// shardAllRef return picked all period ref
//
func (c *CounterFirestore) shardAllRef() *firestore.DocumentRef {
	return c.conn.getDocRef(c.tableName, c.id+string(HierarchyTotal)+"_"+c.pickedShard)
}

// IncrementRX increments a randomly picked shard. must used it in transaction with IncrementWX()
//
//	err = counter.IncrementRX(ctx,1)
//
func (c *CounterFirestore) IncrementRX(ctx context.Context) error {
	if c.conn.tx == nil {
		return errors.New("IncrementRX() must run in transaction")
	}
	c.callRX = true
	if c.pickedShard == "" {
		c.pickedShard = strconv.Itoa(rand.Intn(c.numShards)) //random pick a shard
	}
	//fmt.Printf("counter pick:" + c.shardPick + "\n")

	var err error
	if c.keepDateHierarchy {
		utcNow := time.Now().UTC()
		year := strconv.Itoa(utcNow.Year())
		month := strconv.Itoa(int(utcNow.Month()))
		day := strconv.Itoa(int(utcNow.Day()))
		hour := strconv.Itoa(int(utcNow.Hour()))
		yearRef := c.conn.getDocRef(c.tableName, c.id+year+"_"+c.pickedShard)
		monthRef := c.conn.getDocRef(c.tableName, c.id+year+"-"+month+"_"+c.pickedShard)
		dayRef := c.conn.getDocRef(c.tableName, c.id+year+"-"+month+"-"+day+"_"+c.pickedShard)
		hourRef := c.conn.getDocRef(c.tableName, c.id+year+"-"+month+"-"+day+"-"+hour+"_"+c.pickedShard)

		c.shardHourExist, err = c.isShardExist(ctx, hourRef)
		if err != nil {
			return err
		}
		if c.shardHourExist {
			c.shardDayExist = true
			c.shardMonthExist = true
			c.shardYearExist = true
			c.shardAllExist = true
			return nil
		}

		c.shardDayExist, err = c.isShardExist(ctx, dayRef)
		if err != nil {
			return err
		}
		if c.shardDayExist {
			c.shardMonthExist = true
			c.shardYearExist = true
			c.shardAllExist = true
			return nil
		}

		c.shardMonthExist, err = c.isShardExist(ctx, monthRef)
		if err != nil {
			return err
		}
		if c.shardMonthExist {
			c.shardYearExist = true
			c.shardAllExist = true
			return nil
		}

		c.shardYearExist, err = c.isShardExist(ctx, yearRef)
		if err != nil {
			return err
		}
		if c.shardYearExist {
			c.shardAllExist = true
			return nil
		}
	}

	c.shardAllExist, err = c.isShardExist(ctx, c.shardAllRef())
	if err != nil {
		return err
	}
	return nil
}

// IncrementWX commit IncrementRX()
//
//	err = counter.IncrementWX(ctx)
//
func (c *CounterFirestore) IncrementWX(ctx context.Context, value interface{}) error {
	if c.conn.tx == nil {
		return errors.New("IncrementWX() must run in transaction")
	}
	if c.callRX == false {
		return errors.New("IncrementWX() need call IncrementRX() first")
	}

	utcNow := time.Now().UTC()
	shard := map[string]interface{}{
		MetaID:      c.id,
		MetaValue:   value,
		CounterTime: utcNow,
	}
	if c.keepDateHierarchy {
		year := strconv.Itoa(utcNow.Year())
		month := strconv.Itoa(int(utcNow.Month()))
		day := strconv.Itoa(int(utcNow.Day()))
		hour := strconv.Itoa(int(utcNow.Hour()))
		yearRef := c.conn.getDocRef(c.tableName, c.id+year+"_"+c.pickedShard)
		monthRef := c.conn.getDocRef(c.tableName, c.id+year+"-"+month+"_"+c.pickedShard)
		dayRef := c.conn.getDocRef(c.tableName, c.id+year+"-"+month+"-"+day+"_"+c.pickedShard)
		hourRef := c.conn.getDocRef(c.tableName, c.id+year+"-"+month+"-"+day+"-"+hour+"_"+c.pickedShard)

		if c.shardHourExist {
			if err := c.incrementShard(hourRef, value); err != nil {
				return errors.Wrap(err, "Failed to increment shard hour")
			}
		} else {
			shard[CounterDateLevel] = HierarchyHour
			if err := c.createShard(hourRef, shard); err != nil {
				return errors.Wrap(err, "Failed to create shard hour")
			}
		}

		if c.shardDayExist {
			if err := c.incrementShard(dayRef, value); err != nil {
				return errors.Wrap(err, "Failed to increment shard day")
			}
		} else {
			shard[CounterDateLevel] = HierarchyDay
			if err := c.createShard(dayRef, shard); err != nil {
				return errors.Wrap(err, "Failed to create shard day")
			}
		}

		if c.shardMonthExist {
			if err := c.incrementShard(monthRef, value); err != nil {
				return errors.Wrap(err, "Failed to increment shard month")
			}
		} else {
			shard[CounterDateLevel] = HierarchyMonth
			if err := c.createShard(monthRef, shard); err != nil {
				return errors.Wrap(err, "Failed to create shard month")
			}
		}

		if c.shardYearExist {
			if err := c.incrementShard(yearRef, value); err != nil {
				return errors.Wrap(err, "Failed to increment shard year")
			}
		} else {
			shard[CounterDateLevel] = HierarchyYear
			if err := c.createShard(yearRef, shard); err != nil {
				return errors.Wrap(err, "Failed to create shard year")
			}
		}
	}

	if c.shardAllExist {
		if err := c.incrementShard(c.shardAllRef(), value); err != nil {
			return errors.Wrap(err, "Failed to increment shard all")
		}
	} else {
		shard[CounterDateLevel] = HierarchyTotal
		if err := c.createShard(c.shardAllRef(), shard); err != nil {
			return errors.Wrap(err, "Failed to create shard all")
		}

	}
	c.callRX = false
	c.shardAllExist = false
	c.shardYearExist = false
	c.shardMonthExist = false
	c.shardDayExist = false
	c.shardHourExist = false
	return nil
}

// mock create mock data
//
func (c *CounterFirestore) mock(hierarchy Hierarchy, date time.Time, pick int, value interface{}) error {

	c.pickedShard = strconv.Itoa(pick)
	shard := map[string]interface{}{
		MetaID:           c.id,
		MetaValue:        value,
		CounterDateLevel: hierarchy,
		CounterTime:      date,
	}
	if err := c.createShard(c.shardAllRef(), shard); err != nil {
		return errors.Wrap(err, "Failed to create shard at mock")
	}
	return nil
}

// CountAll return a total count across all period. this function not support transation cause it easily cause "Too much contention on these documents"
//
//	count, err = counter.CountAll(ctx)
//
func (c *CounterFirestore) CountAll(ctx context.Context) (float64, error) {
	tableRef := c.conn.getCollectionRef(c.tableName)
	shards := tableRef.Where(MetaID, "==", c.id).Where(CounterDateLevel, "==", HierarchyTotal).Documents(ctx)
	return c.countValue(shards)
}

// CountPeriod return count between from and to. this function not support transation cause it easily cause "Too much contention on these documents"
//
//	from := time.Date(now.Year()-1, 01, 01, 0, 0, 0, 0, time.UTC)
//	to := time.Date(now.Year()+1, 01, 01, 0, 0, 0, 0, time.UTC)
//	count, err := counter.CountPeriod(ctx, HierarchyYear, from, to)
//
func (c *CounterFirestore) CountPeriod(ctx context.Context, hierarchy Hierarchy, from, to time.Time) (float64, error) {
	tableRef := c.conn.getCollectionRef(c.tableName)
	shards := tableRef.Where(MetaID, "==", c.id).Where(CounterDateLevel, "==", string(hierarchy)).Where(CounterTime, ">=", from).Where(CounterTime, "<=", to).Documents(ctx)
	return c.countValue(shards)
}

// DetailPeriod return detail between from and to. this function not support transation cause it easily cause "Too much contention on these documents"
//
//	dict, err = counter.DetailPeriod(ctx)
//
func (c *CounterFirestore) DetailPeriod(ctx context.Context, hierarchy Hierarchy, from, to time.Time) (map[time.Time]float64, error) {
	result := map[time.Time]float64{}

	tableRef := c.conn.getCollectionRef(c.tableName)
	shards := tableRef.Where(MetaID, "==", c.id).Where(CounterDateLevel, "==", string(hierarchy)).Where(CounterTime, ">=", from).Where(CounterTime, "<=", to).Documents(ctx)
	defer shards.Stop()
	for {
		snotshot, err := shards.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, errors.Wrap(err, "failed to iterator shards at detail period: "+c.errorID())
		}

		obj := snotshot.Data()
		iValue := obj[MetaValue]
		value, err := util.ToFloat64(iValue)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get value on shards, invalid dataType %T, want float64: "+c.errorID(), iValue)
		}
		iDate := obj[CounterTime]
		date := iDate.(time.Time)

		if val, ok := result[date]; ok {
			result[date] = value + val
		} else {
			result[date] = value
		}
	}
	return result, nil
}

// Clear all shards
//
//	err = c.Clear(ctx)
//
func (c *CounterFirestore) Clear(ctx context.Context) error {
	return c.clear(ctx)
}

// ShardsCount returns shards count
//
//	count, err = coder.ShardsCount()
//
func (c *CounterFirestore) ShardsCount(ctx context.Context) (int, error) {
	return c.shardsCount(ctx)
}
