package data

import "time"

// Object represent single database object
//
type Object interface {

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

	// UpdateTime return object update time
	//
	GetUpdateTime() time.Time

	// SetUpdateTime set object update time
	//
	SetUpdateTime(time.Time)
}

// BaseObject represent object stored in document database
//
type BaseObject struct {
	Object

	// ID is object unique identifier used for other object to reference
	//
	ID string `firestore:"-"`

	// reference used by connection implementation
	//
	Ref interface{} `firestore:"-"`

	// CreateTime is object create time, this is readonly field
	//
	CreateTime time.Time `firestore:"-"`

	// UpdateTime is object update time, this is readonly field
	//
	UpdateTime time.Time `firestore:"-"`
}

// GetID return object unique identifier
//
//	d := &Sample{}
//	id := d.ID()
//
func (c *BaseObject) GetID() string {
	return c.ID
}

// SetID set object unique identifier
//
//	d := &Sample{}
//	id := d.setID("uniqueID")
//
func (c *BaseObject) SetID(id string) {
	c.ID = id
}

// GetRef return reference which used by db implementation
//
//	ref := d.Ref()
//
func (c *BaseObject) GetRef() interface{} {
	return c.Ref
}

// SetRef set reference which used by db implementation
//
//	d.setRef(ref)
//
func (c *BaseObject) SetRef(ref interface{}) {
	c.Ref = ref
}

// GetCreateTime return object create time
//
//	id := d.CreateTime()
//
func (c *BaseObject) GetCreateTime() time.Time {
	return c.CreateTime
}

// SetCreateTime set object create time
//
//	id := d.SetCreateTime(time.Now())
//
func (c *BaseObject) SetCreateTime(t time.Time) {
	c.CreateTime = t
}

// GetUpdateTime return object update time
//
//	id := d.UpdateTime()
//
func (c *BaseObject) GetUpdateTime() time.Time {
	return c.UpdateTime
}

// SetUpdateTime set object update time
//
//	id := d.SetUpdateTime(time.Now())
//
func (c *BaseObject) SetUpdateTime(t time.Time) {
	c.UpdateTime = t
}
