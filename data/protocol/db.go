package protocol

import (
	"context"
	"math/rand"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

//DB interface
type DB interface {
	Close()
	Put(ctx context.Context, obj Object) error
	Update(ctx context.Context, objClass string, objID string, fields map[string]interface{}) error
	Get(ctx context.Context, obj Object) error
	GetByClass(ctx context.Context, class string, obj Object) error
	GetAll(ctx context.Context, factory func() Object, callback func(o Object), limit int) error
	ListAll(ctx context.Context, factory func() Object, limit int) ([]Object, error)
	Delete(ctx context.Context, obj Object) error
	DeleteAll(ctx context.Context, class string, timeout int) (int, error)
	Select(ctx context.Context, factory func() Object) Query
	RunTransaction(ctx context.Context, callback func(tx Transaction) error) error

	//AddCount(className string) (int, int, error)
	//Counter(className string) (int, error)
}

// DB simplify datastore create
type db struct {
}

//Put object into data store
func (d *db) Put(obj Object) error {
	panic("need implement Put")
}

//Get object from data store, return ErrNotFound if object not exist
func (d *db) Get(obj Object) error {
	panic("need implement Gut")
}

// Shard use by Counter()
type Shard struct {
	Object
	c int // counter
}

// ErrTimeout is returned by DeleteAll method when the method is run too long
var ErrTimeout = errors.New("db operation timeout")

// ErrNotFound is returned by Get method object not exist
var ErrNotFound = errors.New("object not found")

//AddCount implement sharding counter, shards limit usually 20
//return shard number, shard count, error
func (d *db) AddCount(className string, shards int) (int, int, error) {
	rand.Seed(time.Now().UTC().UnixNano())
	shardNumber := rand.Intn(shards)
	shard := Shard{}
	shard.SetID(strconv.Itoa(shardNumber))

	err := d.Get(&shard)
	if err != nil {
		if err == ErrNotFound {
		} else {
			return shardNumber, -1, err
		}
	}

	shard.c++

	if err := d.Put(&shard); err != nil {
		return shardNumber, -1, err
	}
	return shardNumber, shard.c, nil
}
