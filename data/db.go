package data

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

//"github.com/pkg/errors"
//"math/rand"
//"time"

// IDB is db interface
type IDB interface {
	Close()
	Put(obj IObject) error
	Update(objClass string, objID string, fields map[string]interface{}) error
	Get(obj IObject) error
	GetByClass(class string, obj IObject) error
	GetAll(factory func() IObject, callback func(o IObject), limit int) error
	ListAll(factory func() IObject, limit int) ([]IObject, error)
	Delete(obj IObject) error
	DeleteAll(class string, timeout int) (int, error)
	Select(factory func() IObject) IQuery
	RunTransaction(callback func(tx ITransaction) error) error

	//AddCount(className string) (int, int, error)
	//Counter(className string) (int, error)
}

// DB simplify datastore create
type DB struct {
}

//Put object into data store
func (db *DB) Put(obj IObject) error {
	panic("need implement Put")
}

//Get object from data store, return ErrNotFound if object not exist
func (db *DB) Get(obj IObject) error {
	panic("need implement Gut")
}

// Shard use by Counter()
type Shard struct {
	IObject
	c int // counter
}

// ErrTimeout is returned by DeleteAll method when the method is run too long
var ErrTimeout = errors.New("db operation timeout")

// ErrNotFound is returned by Get method object not exist
var ErrNotFound = errors.New("object not found")

//AddCount implement sharding counter, shards limit usually 20
//return shard number, shard count, error
func (db *DB) AddCount(className string, shards int) (int, int, error) {
	rand.Seed(time.Now().UTC().UnixNano())
	shardNumber := rand.Intn(shards)
	shard := Shard{}
	shard.SetID(strconv.Itoa(shardNumber))

	err := db.Get(&shard)
	if err != nil {
		if err == ErrNotFound {
		} else {
			return shardNumber, -1, err
		}
	}

	shard.c++

	if err := db.Put(&shard); err != nil {
		return shardNumber, -1, err
	}
	return shardNumber, shard.c, nil
}
