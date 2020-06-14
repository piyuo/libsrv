package data

// ShortID create unique serial number, please be aware serial can only generate one number per second and use with transation to ensure unique
//
type ShortID struct {
	ID int64
}

// Number returns a serial number
//
func (c *ShortID) Number() int64 {

	return c.ID
}
