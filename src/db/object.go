package db

import "time"

// Object is any defined object in a database that is used to store or reference data
//
type Object interface {

	// Factory create a empty object, return object must be nil safe, no nil in any field
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

	// SetCreateTime return object create time
	//
	//	d.SetCreateTime(time.Now().UTC())
	//
	SetCreateTime(value time.Time)

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

// Entity a class that has id
//
type Entity struct {
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
	Createtime time.Time `firestore:"CreateTime"`
}

// ID return object unique identifier
//
//	d := &Sample{}
//	id := d.ID()
//
func (c *Entity) ID() string {
	return c.id
}

// SetID set object unique identifier
//
//	d := &Sample{}
//	id := d.setID("uniqueID")
//
func (c *Entity) SetID(id string) {
	c.id = id
}

// Ref return reference which used by db implementation
//
//	ref := d.Ref()
//
func (c *Entity) Ref() interface{} {
	return c.ref
}

// SetRef set reference which used by db implementation
//
//	d.SetRef(ref)
//
func (c *Entity) SetRef(ref interface{}) {
	c.ref = ref
}

// CreateTime return object create time
//
//	t := d.CreateTime()
//
func (c *Entity) CreateTime() time.Time {
	return c.Createtime
}

// SetCreateTime return object create time
//
//	d.SetCreateTime(time.Now().UTC())
//
func (c *Entity) SetCreateTime(value time.Time) {
	if c.Createtime.IsZero() {
		c.Createtime = value
	}
}

// UserID return owner's user id
//
func (c *Entity) UserID() string {
	return ""
}

// SetUserID set owner's user id
//
func (c *Entity) SetUserID(userID string) {
}

// AccountID return owner's account id
//
func (c *Entity) AccountID() string {
	return ""
}

// SetAccountID set owner's account id
//
//	d.SetAccountID(accountID)
//
func (c *Entity) SetAccountID(accountID string) {
}

// UpdateTime return object last update time
//
//	t := d.UpdateTime()
//
func (c *Entity) UpdateTime() time.Time {
	return time.Time{}
}

// SetUpdateTime set object latest update time
//
//	d.SetUpdateTime(time.Now().UTC())
//
func (c *Entity) SetUpdateTime(t time.Time) {
}

// Model is class with AccountID and UserID
//
type Model struct {
	Entity

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
func (c *Model) UserID() string {
	return c.UserId
}

// SetUserID set owner's user id
//
//	d.SetUserID(userID)
//
func (c *Model) SetUserID(userID string) {
	c.UserId = userID
}

// AccountID return owner's account id
//
//	accountID := d.AccountID()
//
func (c *Model) AccountID() string {
	return c.AccountId
}

// SetAccountID set owner's account id
//
//	d.SetAccountID(accountID)
//
func (c *Model) SetAccountID(accountID string) {
	c.AccountId = accountID
}

// UpdateTime return object last update time
//
//	t := d.UpdateTime()
//
func (c *Model) UpdateTime() time.Time {
	return c.Updatetime
}

// SetUpdateTime set object latest update time
//
//	d.SetUpdateTime(time.Now().UTC())
//
func (c *Model) SetUpdateTime(t time.Time) {
	c.Updatetime = t
}
