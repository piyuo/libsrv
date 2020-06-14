package data

// Object is base data interface
type Object interface {
	//ID is object unique id
	//
	//	d := &Sample{}
	//	id := d.ID()
	//
	ID() string

	// SetID set object unique id
	//
	//	d := &Sample{}
	//	id := d.SetID("uniqueID")
	//
	SetID(id string)

	// ModelName return object name in data store
	//
	//	type Sample struct {
	//	StoredObject
	//	}
	//	func (do *Sample) ModelName() string {
	//		return "Sample"
	//	}
	//
	ModelName() string
}

// StoredObject mean object will be saved to data store
//
//	type Sample struct {
//	StoredObject
//	}
//	func (do *Sample) ModelName() string {
//		return "Sample"
//	}
//
type StoredObject struct {
	id string
}

//ID is object unique id
//
//	d := &Sample{}
//	id := d.ID()
//
func (o *StoredObject) ID() string {
	return o.id
}

// SetID set object unique id
//
//	d := &Sample{}
//	id := d.SetID("uniqueID")
//
func (o *StoredObject) SetID(newID string) {
	o.id = newID
}

// ModelName return object name in data store
//
//	type Sample struct {
//	StoredObject
//	}
//	func (do *Sample) ModelName() string {
//		return "Sample"
//	}
//
func (o *StoredObject) ModelName() string {
	panic("not implement ModelName()")
}
