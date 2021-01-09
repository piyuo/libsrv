package data

// Counters is collection of counter
//
type Counters struct {
	Connection Connection

	//TableName is counter table name
	//
	TableName string
}

// DateHierarchy used in create counter
//
type DateHierarchy int8

const (
	// DateHierarchyNone create counter without date hierarchy, only total count
	//
	DateHierarchyNone DateHierarchy = 1

	// DateHierarchyFull create counter with year/month/day/hour hierarchy and total count
	//
	DateHierarchyFull = 2
)

// Counter return counter from database, create one if not exist, set numshards 100 times of concurrent usage. for example if you think concurrent use is 10/seconds then set numshards to 1000 to avoid too much retention error
// if keepDateHierarchy is true, counter will automatically generate year/month/day/hour hierarchy in utc timezone
//
//	counters := db.Counters()
//	orderCountCounter,err = counters.Counter("order-count",100,true) // utc timezone
//
func (c *Counters) Counter(name string, numshards int, hierarchy DateHierarchy) Counter {
	if numshards <= 0 {
		numshards = 10
	}

	return &CounterFirestore{
		MetaFirestore: MetaFirestore{
			conn:      c.Connection.(*ConnectionFirestore),
			tableName: c.TableName,
			id:        name,
			numShards: numshards,
		},
		keepDateHierarchy: hierarchy == DateHierarchyFull,
	}
}
