package data

// Object is base data interface
type Object interface {
	Class() string
	ID() string
	SetID(id string)
	New() interface{}
}

// StoredObject is basic data type
type StoredObject struct {
	id string
}

//ID can get object unique id
func (o *StoredObject) ID() string {
	return o.id
}

//SetID can set object unique id
func (o *StoredObject) SetID(newID string) {
	o.id = newID
}

//Class get object db represent name
func (o *StoredObject) Class() string {
	panic("not implement New()")
}

//New create new object instance
func (o *StoredObject) New() interface{} {
	panic("not implement New()")
}
