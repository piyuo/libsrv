package data

import (
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
	return &CounterFirestore{
		MetaFirestore: MetaFirestore{
			conn:      c.Connection.(*ConnectionFirestore),
			tableName: c.TableName,
			id:        name,
			numShards: numshards,
		},
		loc:    loc,
		native: time.Now().In(loc),
	}
}
