package data

import (
	"context"

	"github.com/pkg/errors"
)

//DB interface
type DB interface {
	Close()
	Put(ctx context.Context, obj Object) error
	Update(ctx context.Context, objClass string, objID string, fields map[string]interface{}) error
	Get(ctx context.Context, obj Object) error
	GetByModelName(ctx context.Context, modelName string, obj Object) error
	GetAll(ctx context.Context, factory func() Object, callback func(o Object), limit int) error
	ListAll(ctx context.Context, factory func() Object, limit int) ([]Object, error)
	Delete(ctx context.Context, obj Object) error
	DeleteAll(ctx context.Context, class string, timeout int) (int, error)
	Select(ctx context.Context, factory func() Object) Query
	RunTransaction(ctx context.Context, callback func(ctx context.Context, tx Transaction) error) error
	Exist(ctx context.Context, path, field, op string, value interface{}) (bool, error)
	Count10(ctx context.Context, path, field, op string, value interface{}) (int, error)
	Increment(ctx context.Context, modelName, modelField, objectID string, value int) error
	//AddCount(className string) (int, int, error)
	//Counter(className string) (int, error)
}

// AbstractDB is parent class for all DB child
type AbstractDB struct {
	DB
}

// Shard use by Counter()
type Shard struct {
	Object
	c int // counter
}

// ErrOperationTimeout is returned by DeleteAll method when the method is run too long
var ErrOperationTimeout = errors.New("db operation timeout")

// ErrObjectNotFound is returned by Get method object not exist
var ErrObjectNotFound = errors.New("object not found")

//AddCount implement sharding counter, shards limit usually 20
//return shard number, shard count, error
/*
func (d *db) AddCount(className string, shards int) (int, int, error) {
	rand.Seed(time.Now().UTC().UnixNano())
	shardNumber := rand.Intn(shards)
	shard := Shard{}
	shard.SetID(strconv.Itoa(shardNumber))

	err := d.Get(&shard)
	if err != nil {
		if err == ErrObjectNotFound {
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
*/
