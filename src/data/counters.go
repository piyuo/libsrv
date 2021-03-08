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

