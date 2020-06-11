package data

// Object is base data interface
type Object interface {
	Class() string
	ID() string
	SetID(id string)
	New() interface{}
}

// DBObject is basic data type
type DBObject struct {
	id string
}

//ID can get object unique id
func (o *DBObject) ID() string {
	return o.id
}

//SetID can set object unique id
func (o *DBObject) SetID(newID string) {
	o.id = newID
}

//Class get object db represent name
func (o *DBObject) Class() string {
	panic("not implement New()")
}

//New create new object instance
func (o *DBObject) New() interface{} {
	panic("not implement New()")
}
