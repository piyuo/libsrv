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

	// TimeCreated return object create time
	//
	//	t := d.TimeCreated()
	//
	TimeCreated() time.Time

	// setCreated set object create time
	//
	//	d.setCreated(time.Now().UTC())
	//
	setCreated(t time.Time)

	// TimeUpdated return object latest update time
	//
	//	t := d.TimeUpdated()
	//
	TimeUpdated() time.Time

	// setUpdated set object latest update time
	//
	//	d.setUpdated(time.Now().UTC())
	//
	setUpdated(t time.Time)
}

// BaseObject represent object stored in document database
//
type BaseObject struct {
	Object `firestore:"-"`

	// Created is object create time
	//
	Created time.Time

	// Updated is object latest update time
	//
	Updated time.Time

	// ID is object unique identifier used for other object to reference
	//
	ID string `firestore:"-"`

	// reference used by connection implementation
	//
	Ref interface{} `firestore:"-"`
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

// TimeCreated return object create time
//
//	t := d.TimeCreated()
//
func (c *BaseObject) TimeCreated() time.Time {
	return c.Created
}

// setCreated set object create time
//
//	d.setCreated(time.Now().UTC())
//
func (c *BaseObject) setCreated(t time.Time) {
	if c.Created.IsZero() {
		c.Created = t
	}
}

// TimeUpdated return object latest update time
//
//	t := d.TimeUpdated()
//
func (c *BaseObject) TimeUpdated() time.Time {
	return c.Updated
}

// setUpdated set object latest update time
//
//	d.setCreated(time.Now().UTC())
//
func (c *BaseObject) setUpdated(t time.Time) {
	c.Updated = t
}
