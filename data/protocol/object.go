package protocol

// Object is base data interface
type Object interface {
	Class() string
	ID() string
	SetID(id string)
	New() interface{}
}

// Object is basic data type
type object struct {
	id string
}

//ID can get object unique id
func (o *object) ID() string {
	return o.id
}

//SetID can set object unique id
func (o *object) SetID(newID string) {
	o.id = newID
}

//Class get object db represent name
func (o *object) Class() string {
	panic("not implement New()")
}

//New create new object instance
func (o *object) New() interface{} {
	panic("not implement New()")
}
