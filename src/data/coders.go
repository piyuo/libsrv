package data

// Coders is collection of code
//
type Coders struct {
	Connection Connection

	//TableName is code table name
	//
	TableName string
}

// Coder return code from database, set numshards 100 times of concurrent usage. for example if you think concurrent use is 10/seconds then set numshards to 1000 to avoid too much retention error
//
//	coders := db.Coders()
//	productCoder,err = coders.Coder("product-code",100)
//
func (c *Coders) Coder(name string, numshards int) Coder {
	return &CoderFirestore{
		MetaFirestore: MetaFirestore{
			conn:      c.Connection.(*ConnectionFirestore),
			tableName: c.TableName,
			id:        name,
			numShards: numshards,
		},
	}
}
