package data

// Object is base data interface
type Object interface {
	ModelName() string
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

// ModelName get object db represent name
//
//	greet.ModelName() // return "greet"
//
func (o *StoredObject) ModelName() string {
	panic("not implement ModelName()")
}

//New create new object instance
func (o *StoredObject) New() interface{} {
	panic("not implement New()")
}
