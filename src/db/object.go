package db

import "time"

// Object is any defined object in a database that is used to store or reference data
//
type Object interface {

	// Factory create a empty object
	//
	//	newSample = sample.Factory()
	//
	Factory() Object

	// Collection name
	//
	//	collection = sample.Collection() // "Sample"
	//
	Collection() string

	// ID is object unique identifier used for other object to reference
	//
	//	id := ID()
	//
	ID() string

	// SetID is object unique identifier used for other object to reference
	//
	//	id := sample.SetID("id")
	//
	SetID(id string)

	// Ref return reference which used by db implementation
	//
	//	ref := sample.Ref()
	//
	Ref() interface{}

	// SetRef set reference which used by db implementation
	//
	//	sample.SetRef(ref)
	//
	SetRef(ref interface{})

	// CreateTime return object create time
	//
	//	t := sample.CreateTime()
	//
	CreateTime() time.Time

	// UpdateTime return object last update time
	//
	//	t := sample.UpdateTime()
	//
	UpdateTime() time.Time

	// SetUpdateTime set object latest update time
	//
	//	sample.SetUpdateTime(time.Now().UTC())
	//
	SetUpdateTime(t time.Time)

	// AccountID return owner's account id
	//
	//	accountID := sample.AccountID()
	//
	AccountID() string

	// SetAccountID set owner's account id
	//
	//	sample.SetAccountID(userID)
	//
	SetAccountID(accountID string)

	// UserID return owner's user id
	//
	//	userID := sample.UserID()
	//
	UserID() string

	// SetUserID set owner's user id
	//
	//	sample.SetUserID(userID)
	//
	SetUserID(userID string)
}

// BaseObject represent object stored in document database
//
type BaseObject struct {
	Object `firestore:"-"`

	// id is object unique identifier
	//
	id string // lowercase private field will not save to database

	// ref use in connection implementation
	//
	ref interface{} // lowercase private field will not save to database

	// Createtime is object create time, you should always use CreateTime() SetCreateTime() to access this field
	// We keep our own create time, cause database provide create time like "snapshot.CreateTime" may not use in query
	//
	Createtime time.Time `firestore:"CreateTime,serverTimestamp"`
}

// ID return object unique identifier
//
//	d := &Sample{}
//	id := d.ID()
//
func (c *BaseObject) ID() string {
	return c.id
}

// SetID set object unique identifier
//
//	d := &Sample{}
//	id := d.setID("uniqueID")
//
func (c *BaseObject) SetID(id string) {
	c.id = id
}

// Ref return reference which used by db implementation
//
//	ref := d.Ref()
//
func (c *BaseObject) Ref() interface{} {
	return c.ref
}

// SetRef set reference which used by db implementation
//
//	d.SetRef(ref)
//
func (c *BaseObject) SetRef(ref interface{}) {
	c.ref = ref
}

// CreateTime return object create time
//
//	t := d.CreateTime()
//
func (c *BaseObject) CreateTime() time.Time {
	return c.Createtime
}

// UserID return owner's user id
//
func (c *BaseObject) UserID() string {
	return ""
}

// SetUserID set owner's user id
//
func (c *BaseObject) SetUserID(userID string) {
}

// AccountID return owner's account id
//
func (c *BaseObject) AccountID() string {
	return ""
}

// SetAccountID set owner's account id
//
//	d.SetAccountID(accountID)
//
func (c *BaseObject) SetAccountID(accountID string) {
}

// UpdateTime return object last update time
//
//	t := d.UpdateTime()
//
func (c *BaseObject) UpdateTime() time.Time {
	return time.Time{}
}

// SetUpdateTime set object latest update time
//
//	d.SetUpdateTime(time.Now().UTC())
//
func (c *BaseObject) SetUpdateTime(t time.Time) {
}

// DomainObject is object with AccountID and UserID
//
type DomainObject struct {
	BaseObject

	// Updatetime is object last update time, you should use UpdateTime() SetUpdateTime() to access this field
	// We keep our own update time, cause database provide update time like "snapshot.UpdateTime" may not use in query
	//
	Updatetime time.Time `firestore:"UpdateTime,omitempty"`

	// AccountId is owner's account id, you should use AccountID() SetAccountID() to access this field
	//
	AccountId string `firestore:"AccountID,omitempty"`

	// UserId is owner's user id, you should use UserID() SetUserID() to access this field
	//
	UserId string `firestore:"UserID,omitempty"`
}

// UserID return owner's user id
//
//	userID := d.UserID()
//
func (c *DomainObject) UserID() string {
	return c.UserId
}

// SetUserID set owner's user id
//
//	d.SetUserID(userID)
//
func (c *DomainObject) SetUserID(userID string) {
	c.UserId = userID
}

// AccountID return owner's account id
//
//	accountID := d.AccountID()
//
func (c *DomainObject) AccountID() string {
	return c.AccountId
}

// SetAccountID set owner's account id
//
//	d.SetAccountID(accountID)
//
func (c *DomainObject) SetAccountID(accountID string) {
	c.AccountId = accountID
}

// UpdateTime return object last update time
//
//	t := d.UpdateTime()
//
func (c *DomainObject) UpdateTime() time.Time {
	return c.Updatetime
}

// SetUpdateTime set object latest update time
//
//	d.SetUpdateTime(time.Now().UTC())
//
func (c *DomainObject) SetUpdateTime(t time.Time) {
	c.Updatetime = t
}
