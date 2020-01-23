package command

//TestMap mapping ID and object
type TestMap struct {
}

// IDToObject return object from id
func (r *TestMap) NewObjectByID(id uint16) interface{} {
	switch id {
	case 1:
		return new(TestAction)
	case 3:
		return new(TestActionNotRespond)
	}
	return nil
}
