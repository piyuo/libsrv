package gdb

import (
	"context"
	"math/rand"
	"strconv"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/piyuo/libsrv/src/db"
	"github.com/piyuo/libsrv/src/log"
	"github.com/piyuo/libsrv/src/util"
	"github.com/pkg/errors"
	"google.golang.org/api/iterator"
)

// CounterFirestore implement Counter
//
type CounterFirestore struct {
	db.Counter `firestore:"-"`

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

// shardAllRef return picked all period ref
//
func (c *CounterFirestore) shardAllRef() *firestore.DocumentRef {
	return c.client.getDocRef(c.collection, c.id+string(db.HierarchyTotal)+"_"+c.pickedShard)
}

// IncrementRX increments a randomly picked shard. must used it in transaction with IncrementWX()
//
//	err := Transaction(ctx, func(ctx context.Context,tx db.Transaction) error {
//		err = counter.IncrementRX(ctx,transaction)
//		err = counter.IncrementWX(ctx,transaction,1)
//	})
//
func (c *CounterFirestore) IncrementRX(ctx context.Context, transaction db.Transaction) error {
	tx := transaction.(*TransactionFirestore)
	c.callRX = true
	c.pickedShard = strconv.Itoa(rand.Intn(c.numShards)) //random pick a shard
	log.Debug(ctx, "counter pick %v from %v shards", c.pickedShard, c.numShards)

	var err error
	if c.keepDateHierarchy {
		utcNow := time.Now().UTC()
		year := strconv.Itoa(utcNow.Year())
		month := strconv.Itoa(int(utcNow.Month()))
		day := strconv.Itoa(int(utcNow.Day()))
		hour := strconv.Itoa(int(utcNow.Hour()))
		yearRef := c.client.getDocRef(c.collection, c.id+year+"_"+c.pickedShard)
		monthRef := c.client.getDocRef(c.collection, c.id+year+"-"+month+"_"+c.pickedShard)
		dayRef := c.client.getDocRef(c.collection, c.id+year+"-"+month+"-"+day+"_"+c.pickedShard)
		hourRef := c.client.getDocRef(c.collection, c.id+year+"-"+month+"-"+day+"-"+hour+"_"+c.pickedShard)

		c.shardHourExist, err = tx.isShardExists(ctx, hourRef)
		if err != nil {
			return errors.Wrap(err, "hour")
		}
		if c.shardHourExist {
			c.shardDayExist = true
			c.shardMonthExist = true
			c.shardYearExist = true
			c.shardAllExist = true
			return nil
		}

		c.shardDayExist, err = tx.isShardExists(ctx, dayRef)
		if err != nil {
			return errors.Wrap(err, "day")
		}
		if c.shardDayExist {
			c.shardMonthExist = true
			c.shardYearExist = true
			c.shardAllExist = true
			return nil
		}

		c.shardMonthExist, err = tx.isShardExists(ctx, monthRef)
		if err != nil {
			return errors.Wrap(err, "month")
		}
		if c.shardMonthExist {
			c.shardYearExist = true
			c.shardAllExist = true
			return nil
		}

		c.shardYearExist, err = tx.isShardExists(ctx, yearRef)
		if err != nil {
			return errors.Wrap(err, "year")
		}
		if c.shardYearExist {
			c.shardAllExist = true
			return nil
		}
	}

	c.shardAllExist, err = tx.isShardExists(ctx, c.shardAllRef())
	if err != nil {
		return errors.Wrap(err, "all")
	}
	return nil
}

// IncrementWX commit IncrementRX()
//
//	err := Transaction(ctx, func(ctx context.Context,tx db.Transaction) error {
//		err = counter.IncrementRX(ctx,transaction)
//		err = counter.IncrementWX(ctx,transaction,1)
//	})
//
func (c *CounterFirestore) IncrementWX(ctx context.Context, transaction db.Transaction, value interface{}) error {
	tx := transaction.(*TransactionFirestore)
	if c.callRX == false {
		return errors.New("must call RX first")
	}

	utcNow := time.Now().UTC()
	shard := map[string]interface{}{
		db.MetaID:      c.id,
		db.MetaValue:   value,
		db.CounterTime: utcNow,
	}
	if c.keepDateHierarchy {
		year := strconv.Itoa(utcNow.Year())
		month := strconv.Itoa(int(utcNow.Month()))
		day := strconv.Itoa(int(utcNow.Day()))
		hour := strconv.Itoa(int(utcNow.Hour()))
		yearRef := c.client.getDocRef(c.collection, c.id+year+"_"+c.pickedShard)
		monthRef := c.client.getDocRef(c.collection, c.id+year+"-"+month+"_"+c.pickedShard)
		dayRef := c.client.getDocRef(c.collection, c.id+year+"-"+month+"-"+day+"_"+c.pickedShard)
		hourRef := c.client.getDocRef(c.collection, c.id+year+"-"+month+"-"+day+"-"+hour+"_"+c.pickedShard)

		if c.shardHourExist {
			if err := tx.incrementShard(hourRef, value); err != nil {
				return errors.Wrap(err, "inc hour")
			}
		} else {
			shard[db.CounterDateLevel] = db.HierarchyHour
			if err := tx.createShard(hourRef, shard); err != nil {
				return errors.Wrap(err, "create hour")
			}
		}

		if c.shardDayExist {
			if err := tx.incrementShard(dayRef, value); err != nil {
				return errors.Wrap(err, "inc day")
			}
		} else {
			shard[db.CounterDateLevel] = db.HierarchyDay
			if err := tx.createShard(dayRef, shard); err != nil {
				return errors.Wrap(err, "create day")
			}
		}

		if c.shardMonthExist {
			if err := tx.incrementShard(monthRef, value); err != nil {
				return errors.Wrap(err, "inc month")
			}
		} else {
			shard[db.CounterDateLevel] = db.HierarchyMonth
			if err := tx.createShard(monthRef, shard); err != nil {
				return errors.Wrap(err, "create month")
			}
		}

		if c.shardYearExist {
			if err := tx.incrementShard(yearRef, value); err != nil {
				return errors.Wrap(err, "inc year")
			}
		} else {
			shard[db.CounterDateLevel] = db.HierarchyYear
			if err := tx.createShard(yearRef, shard); err != nil {
				return errors.Wrap(err, "create year")
			}
		}
	}

	if c.shardAllExist {
		if err := tx.incrementShard(c.shardAllRef(), value); err != nil {
			return errors.Wrap(err, "inc all")
		}
	} else {
		shard[db.CounterDateLevel] = db.HierarchyTotal
		if err := tx.createShard(c.shardAllRef(), shard); err != nil {
			return errors.Wrap(err, "create all")
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

// CountAll return a total count across all period. this function not support transation cause it easily cause "Too much contention on these documents"
//
//	count, err = counter.CountAll(ctx)
//
func (c *CounterFirestore) CountAll(ctx context.Context) (float64, error) {
	tableRef := c.client.getCollectionRef(c.collection)
	shards := tableRef.Where(db.MetaID, "==", c.id).Where(db.CounterDateLevel, "==", db.HierarchyTotal).Documents(ctx)
	defer shards.Stop()
	return c.countValue(shards)
}

// CountPeriod return count between from and to. this function not support transation cause it easily cause "Too much contention on these documents"
//
//	from := time.Date(now.Year()-1, 01, 01, 0, 0, 0, 0, time.UTC)
//	to := time.Date(now.Year()+1, 01, 01, 0, 0, 0, 0, time.UTC)
//	count, err := counter.CountPeriod(ctx, HierarchyYear, from, to)
//
func (c *CounterFirestore) CountPeriod(ctx context.Context, hierarchy db.Hierarchy, from, to time.Time) (float64, error) {
	tableRef := c.client.getCollectionRef(c.collection)
	shards := tableRef.Where(db.MetaID, "==", c.id).Where(db.CounterDateLevel, "==", string(hierarchy)).Where(db.CounterTime, ">=", from).Where(db.CounterTime, "<=", to).Documents(ctx)
	defer shards.Stop()
	return c.countValue(shards)
}

// DetailPeriod return detail between from and to. this function not support transation cause it easily cause "Too much contention on these documents"
//
//	dict, err = counter.DetailPeriod(ctx)
//
func (c *CounterFirestore) DetailPeriod(ctx context.Context, hierarchy db.Hierarchy, from, to time.Time) (map[time.Time]float64, error) {
	result := map[time.Time]float64{}

	tableRef := c.client.getCollectionRef(c.collection)
	shards := tableRef.Where(db.MetaID, "==", c.id).Where(db.CounterDateLevel, "==", string(hierarchy)).Where(db.CounterTime, ">=", from).Where(db.CounterTime, "<=", to).Documents(ctx)
	defer shards.Stop()
	for {
		snotshot, err := shards.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, errors.Wrapf(err, "iter next %v-%v", c.collection, c.id)
		}

		obj := snotshot.Data()
		iValue := obj[db.MetaValue]
		value, err := util.ToFloat64(iValue)
		if err != nil {
			return nil, errors.Wrapf(err, "invalid dataType %T want float64 %v-%v", iValue, c.collection, c.id)
		}
		iDate := obj[db.CounterTime]
		date := iDate.(time.Time)

		if val, ok := result[date]; ok {
			result[date] = value + val
		} else {
			result[date] = value
		}
	}
	return result, nil
}

// Delete delete counter
//
//	err = Delete(ctx)
//
func (c *CounterFirestore) Delete(ctx context.Context) error {
	return c.deleteShards(ctx)
}

// ShardsCount returns shards count
//
//	count, err = ShardsCount(ctx)
//
func (c *CounterFirestore) ShardsCount(ctx context.Context) (int, error) {
	return c.shardsCount(ctx)
}
