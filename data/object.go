package data

// IObject is base data interface
type IObject interface {
	Class() string
	ID() string
	SetID(id string)
	New() interface{}
}

// Object is basic data type
type Object struct {
	id string
}

//ID can get object unique id
func (o *Object) ID() string {
	return o.id
}

//SetID can set object unique id
func (o *Object) SetID(newID string) {
	o.id = newID
}

//Class get object db represent name
func (o *Object) Class() string {
	panic("not implement New()")
}

//New create new object instance
func (o *Object) New() interface{} {
	panic("not implement New()")
}
