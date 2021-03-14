package data

// LimitQueryDefault limit query return item
//
const LimitQueryDefault = 10

// LimitTransactionClear limit transaction clear count
//
const LimitTransactionClear = 50

// LimitClear limit clear count
//
const LimitClear = 500

// DB represent DB public method
//
type DB interface {

	// Close connection
	//
	//	c.Close()
	//
	Close()

	// Transaction create transaction
	//
	Transaction() Transaction

	// Batch create batch
	//
	Batch() Batch
	// GetConnection return current connection
	//
	//	conn := c.GetConnection()
	//
	GetConnection() Connection

	// Usage return usage object
	//
	//	usage := c.Usage()
	//
	//	Usage() Usage
}

// BaseDB represent document database
//
type BaseDB struct {
	DB

	// Conn is current database connection
	//
	Connection Connection
}

// GetConnection return current connection
//
//	conn := c.GetConnection()
//
func (c *BaseDB) GetConnection() Connection {
	return c.Connection
}

// Close connection
//
//	c.Close()
//
func (c *BaseDB) Close() {
	if c.Connection != nil {
		c.Connection.Close()
		c.Connection = nil
	}
}

// Transaction create transaction
//
func (c *BaseDB) Transaction() Transaction {
	return c.Connection.CreateTransaction()
}

// Batch create batch
//
func (c *BaseDB) Batch() Batch {
	return c.Connection.CreateBatch()
}
