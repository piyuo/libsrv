package data

import "time"

// Object represent single document in table
//
type Object interface {

	// id is object unique identifier used for other object to reference
	//
	//	d := &Sample{}
	//	err = db.Get(ctx, d)
	//	id := d.ID()
	//
	ID() string

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
	Ref() interface{}

	// SetRef set reference which used by db implementation
	//
	//	d.setRef(ref)
	//
	SetRef(ref interface{})

	// CreateTime return object create time
	//
	CreateTime() time.Time

	// SetCreateTime set object create time
	//
	SetCreateTime(time.Time)

	// ReadTime return object read time
	//
	ReadTime() time.Time

	// SetReadTime set object read time
	//
	SetReadTime(time.Time)

	// UpdateTime return object update time
	//
	UpdateTime() time.Time

	// SetUpdateTime set object update time
	//
	SetUpdateTime(time.Time)
}

// DocObject represent object stored in document database
//
//	type Sample struct {
//	AbstractObject
//	}
//	func (do *Sample) ModelName() string {
//		return "Sample"
//	}
//
type DocObject struct {
	Object

	// id is object unique identifier used for other object to reference
	//
	id string

	// reference which used by db implementation
	//
	ref interface{}

	// createTime is object create time
	//
	createTime time.Time

	// readTime is object read time
	//
	readTime time.Time

	// updateTime is object update time
	//
	updateTime time.Time `firestore:"-"`
}

// ID return object unique identifier
//
//	d := &Sample{}
//	id := d.ID()
//
func (do *DocObject) ID() string {
	return do.id
}

// SetID set object unique identifier
//
//	d := &Sample{}
//	id := d.setID("uniqueID")
//
func (do *DocObject) SetID(id string) {
	do.id = id
}

// Ref return reference which used by db implementation
//
//	ref := d.Ref()
//
func (do *DocObject) Ref() interface{} {
	return do.ref
}

// SetRef set reference which used by db implementation
//
//	d.setRef(ref)
//
func (do *DocObject) SetRef(ref interface{}) {
	do.ref = ref
}

// CreateTime return object create time
//
//	id := d.CreateTime()
//
func (do *DocObject) CreateTime() time.Time {
	return do.createTime
}

// SetCreateTime set object create time
//
//	id := d.SetCreateTime(time.Now())
//
func (do *DocObject) SetCreateTime(t time.Time) {
	do.createTime = t
}

// ReadTime return object create time
//
//	id := d.ReadTime()
//
func (do *DocObject) ReadTime() time.Time {
	return do.readTime
}

// SetReadTime set object read time
//
//	id := d.SetReadTime(time.Now())
//
func (do *DocObject) SetReadTime(t time.Time) {
	do.readTime = t
}

// UpdateTime return object update time
//
//	id := d.UpdateTime()
//
func (do *DocObject) UpdateTime() time.Time {
	return do.updateTime
}

// SetUpdateTime set object update time
//
//	id := d.SetUpdateTime(time.Now())
//
func (do *DocObject) SetUpdateTime(t time.Time) {
	do.updateTime = t
}
