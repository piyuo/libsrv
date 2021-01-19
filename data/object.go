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

	// GetCreateTime return object create time
	//
	//	t := d.GetCreateTime()
	//
	GetCreateTime() time.Time

	// SetCreateTime set object create time, create time will not change if it's not empty
	//
	//	d.SetCreateTime(time.Now().UTC())
	//
	SetCreateTime(t time.Time)

	// GetUpdateTime return object last update time
	//
	//	t := d.GetUpdateTime()
	//
	GetUpdateTime() time.Time

	// SetUpdateTime set object latest update time
	//
	//	d.SetUpdateTime(time.Now().UTC())
	//
	SetUpdateTime(t time.Time)

	// GetAccountID return owner's account id
	//
	//	accountID := d.GetAccountID()
	//
	GetAccountID() string

	// SetAccountID set owner's account id
	//
	//	d.SetAccountID(userID)
	//
	SetAccountID(accountID string)

	// GetAccountID return owner's user id
	//
	//	userID := d.GetUserID()
	//
	GetUserID() string

	// SetUserID set owner's user id
	//
	//	d.SetUserID(userID)
	//
	SetUserID(userID string)
}

// BaseObject represent object stored in document database
//
type BaseObject struct {
	Object `firestore:"-"`

	// ID is object unique identifier used for other object to reference
	//
	ID string `firestore:"-"`

	// reference used by connection implementation
	//
	Ref interface{} `firestore:"-"`

	// CreateTime is object create time
	// We keep our own create time, cause database provide create time like "snapshot.CreateTime" may not use in query
	//
	CreateTime time.Time

	// UpdateTime is object last update time
	// We keep our own create time, cause database provide update time like "snapshot.UpdateTime" may not use in query
	//
	UpdateTime time.Time

	// UserID is owner's user id
	//
	UserID string

	// AccountID is owner's account id
	//
	AccountID string
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
//	t := d.GetCreateTime()
//
func (c *BaseObject) GetCreateTime() time.Time {
	return c.CreateTime
}

// SetCreateTime set object create time
//
//	d.SetCreateTime(time.Now().UTC())
//
func (c *BaseObject) SetCreateTime(t time.Time) {
	if c.CreateTime.IsZero() {
		c.CreateTime = t
	}
}

// GetUpdateTime return object last update time
//
//	t := d.GetUpdateTime()
//
func (c *BaseObject) GetUpdateTime() time.Time {
	return c.UpdateTime
}

// SetUpdateTime set object latest update time
//
//	d.SetUpdateTime(time.Now().UTC())
//
func (c *BaseObject) SetUpdateTime(t time.Time) {
	c.UpdateTime = t
}

// GetUserID return owner's user id
//
//	userID := d.GetUserID()
//
func (c *BaseObject) GetUserID() string {
	return c.UserID
}

// SetUserID set owner's user id
//
//	d.SetUserID(userID)
//
func (c *BaseObject) SetUserID(userID string) {
	c.UserID = userID
}

// GetAccountID return owner's account id
//
//	accountID := d.GetAccountID()
//
func (c *BaseObject) GetAccountID() string {
	return c.AccountID
}

// SetAccountID set owner's account id
//
//	d.SetAccountID(accountID)
//
func (c *BaseObject) SetAccountID(accountID string) {
	c.AccountID = accountID
}
