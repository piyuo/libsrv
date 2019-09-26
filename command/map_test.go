package command

//TestMap mapping ID and object
type TestMap struct {
}

// IDToObject return object from id
func (r *TestMap) IDToObject(id uint16) interface{} {
	switch id {
	case 1:
		return new(TestAction)
	case 2:
		return new(TestResponse)
	case 3:
		return new(TestActionNotRespond)
	}
	return nil
}

// IDCount give total id count
func (r *TestMap) IDCount() uint16 {
	return 3
}
