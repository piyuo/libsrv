package data

import (
	"context"
	"time"
)

// Counters is collection of counter
//
type Counters struct {
	Connection Connection

	//TableName is counter table name
	//
	TableName string
}

// Counter return counter from database, create one if not exist, set numshards 100 times of concurrent usage. for example if you think concurrent use is 10/seconds then set numshards to 1000 to avoid too much retention error
//
//	timezone is user timezone, cause counter will automatically generate year/month/day/hour count
//
//	counters := db.Counters()
//	orderCountCounter,err = counters.Counter("order-count",100,"UTC",0) // utc timezone
//
func (c *Counters) Counter(name string, numshards int, timezoneName string, timezoneOffset int) Counter {

	if numshards <= 0 {
		numshards = 10
	}

	loc := time.FixedZone(timezoneName, timezoneOffset)
	/*
		t := time.Now().In(loc)
		fmt.Printf(t.Format("2006-01-02 15:04:05") + "\n")
		fmt.Printf(strconv.Itoa(int(t.Month())) + "\n")
		name, offset := t.Zone()
		fmt.Printf(name + "," + strconv.Itoa(offset) + "\n")

		d := time.Date(t.Year(), time.Month(1), 01, 01, 0, 0, 0, loc)
		fmt.Printf(d.Format("2006-01-02 15:04:05") + "\n")
		name, offset = d.Zone()
		fmt.Printf(name + "," + strconv.Itoa(offset) + "\n")
	*/
	return &CounterFirestore{
		ShardsFirestore: ShardsFirestore{
			conn:      c.Connection.(*ConnectionFirestore),
			tableName: c.TableName,
			id:        name,
			numShards: numshards,
		},
		loc:    loc,
		native: time.Now().In(loc),
	}
}

// Delete counter from database
//
//	counters := db.Counters()
//	err = counters.Delete(ctx, "myCounter")
//
func (c *Counters) Delete(ctx context.Context, name string) error {
	shards := ShardsFirestore{
		conn:      c.Connection.(*ConnectionFirestore),
		tableName: c.TableName,
		id:        name,
		numShards: 0,
	}
	if err := shards.assert(ctx); err != nil {
		return err
	}

	return shards.deleteShards(ctx)
}
