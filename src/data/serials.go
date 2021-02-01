package data

// Serials is collection of serial
//
type Serials struct {
	Connection Connection

	//TableName is serial table name
	//
	TableName string
}

// Serial return serial from database, create one if not exist, please be aware Serial can only generate 1 number per second, use serial with high frequency will cause too much retention error
//
//	serials := db.Serials()
//	productNo,err = serials.Serial("product-no")
//
func (c *Serials) Serial(name string) Serial {
	return &SerialFirestore{
		MetaFirestore: MetaFirestore{
			conn:      c.Connection.(*ConnectionFirestore),
			tableName: c.TableName,
			id:        name,
			numShards: 0,
		},
	}
}
