package data

import "time"

// ObjectRef is object reference for public method
//
type ObjectRef interface {

	// id is object unique identifier used for other object to reference
	//
	//	d := &Sample{}
	//	err = db.Get(ctx, d)
	//	id := d.ID()
	//
	GetID() string

	// SetID is object unique identifier used for other object to reference
	//
	//	d := &Sample{}
	//	err = db.Get(ctx, d)
	//	id := d.ID()
	//
	SetID(id string)

	// Ref return reference which used by db implementation
	//
	//	ref := d.Ref()
	//
	GetRef() interface{}

	// SetRef set reference which used by db implementation
	//
	//	d.setRef(ref)
	//
	SetRef(ref interface{})

	// CreateTime return object create time
	//
	GetCreateTime() time.Time

	// SetCreateTime set object create time
	//
	SetCreateTime(time.Time)

	// ReadTime return object read time
	//
	GetReadTime() time.Time

	// SetReadTime set object read time
	//
	SetReadTime(time.Time)

	// UpdateTime return object update time
	//
	GetUpdateTime() time.Time

	// SetUpdateTime set object update time
	//
	SetUpdateTime(time.Time)
}

// Object represent object stored in document database
//
type Object struct {
	ObjectRef

	// ID is object unique identifier used for other object to reference
	//
	ID string `firestore:"-"`

	// reference used by connection implementation
	//
	Ref interface{} `firestore:"-"`

	// CreateTime is object create time, this is readonly field
	//
	CreateTime time.Time `firestore:"-"`

	// ReadTime is object read time, this is readonly field
	//
	ReadTime time.Time `firestore:"-"`

	// UpdateTime is object update time, this is readonly field
	//
	UpdateTime time.Time `firestore:"-"`
}

// GetID return object unique identifier
//
//	d := &Sample{}
//	id := d.ID()
//
func (o *Object) GetID() string {
	return o.ID
}

// SetID set object unique identifier
//
//	d := &Sample{}
//	id := d.setID("uniqueID")
//
func (o *Object) SetID(id string) {
	o.ID = id
}

// GetRef return reference which used by db implementation
//
//	ref := d.Ref()
//
func (o *Object) GetRef() interface{} {
	return o.Ref
}

// SetRef set reference which used by db implementation
//
//	d.setRef(ref)
//
func (o *Object) SetRef(ref interface{}) {
	o.Ref = ref
}

// GetCreateTime return object create time
//
//	id := d.CreateTime()
//
func (o *Object) GetCreateTime() time.Time {
	return o.CreateTime
}

// SetCreateTime set object create time
//
//	id := d.SetCreateTime(time.Now())
//
func (o *Object) SetCreateTime(t time.Time) {
	o.CreateTime = t
}

// GetReadTime return object create time
//
//	id := d.ReadTime()
//
func (o *Object) GetReadTime() time.Time {
	return o.ReadTime
}

// SetReadTime set object read time
//
//	id := d.SetReadTime(time.Now())
//
func (o *Object) SetReadTime(t time.Time) {
	o.ReadTime = t
}

// GetUpdateTime return object update time
//
//	id := d.UpdateTime()
//
func (o *Object) GetUpdateTime() time.Time {
	return o.UpdateTime
}

// SetUpdateTime set object update time
//
//	id := d.SetUpdateTime(time.Now())
//
func (o *Object) SetUpdateTime(t time.Time) {
	o.UpdateTime = t
}
