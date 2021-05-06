package gdb

import (
	"context"

	"github.com/piyuo/libsrv/db"
	"github.com/piyuo/libsrv/google/gaccount"
)

type PlainObject struct {
	ID   string
	Name string
}

// SampleNoFactory with no factory and collection
//
type SampleNoFactory struct {
	db.Model
	Name    string
	Value   int
	Map     map[string]string
	Array   []string
	Numbers []int
	pObj    *PlainObject
}

func (c *SampleNoFactory) Factory() db.Object {
	return nil
}

func (c *SampleNoFactory) Collection() string {
	return ""
}

// SampleEmpty for test clear
//
type SampleEmpty struct {
	db.Model
	Name string `firestore:"Name,omitempty"`
}

func (c *SampleEmpty) Factory() db.Object {
	return &SampleEmpty{}
}

func (c *SampleEmpty) Collection() string {
	return "SampleEmpty"
}

// SampleDeleteAll for test clear
//
type SampleDeleteAll struct {
	db.Model
	Name string `firestore:"Name,omitempty"`
}

func (c *SampleDeleteAll) Factory() db.Object {
	return &SampleDeleteAll{}
}

func (c *SampleDeleteAll) Collection() string {
	return "SampleDeleteAll"
}

// Sample for test
//
type Sample struct {
	db.Model
	Name    string            `firestore:"Name,omitempty"`
	Tag     string            `firestore:"Tag,omitempty"`
	Value   int               `firestore:"Value,omitempty"`
	Map     map[string]string `firestore:"Map,omitempty"`
	Array   []string          `firestore:"Array,omitempty"`
	Numbers []int             `firestore:"Numbers,omitempty"`
	PObj    *PlainObject      `firestore:"PObj,omitempty"`
}

// Factory create a empty object, return object must be nil safe, no nil in any field
//
func (c *Sample) Factory() db.Object {
	return &Sample{
		Map:     map[string]string{},
		Array:   []string{},
		Numbers: []int{},
		PObj:    &PlainObject{},
	}
}

// Collection is name in the database
//
func (c *Sample) Collection() string {
	return "Sample"
}

var sampleClientInstance *ClientFirestore

// sample client create db client use for test, it will keep client instance to resuse, recreate new instance if client is close
//
func sampleClient() db.Client {
	if sampleClientInstance != nil && !sampleClientInstance.IsClose() {
		return sampleClientInstance
	}
	ctx := context.Background()
	cred, err := gaccount.GlobalCredential(ctx)
	if err != nil {
		return nil
	}
	client, err := NewClient(ctx, cred)
	if err != nil {
		return nil
	}
	sampleClientInstance = client.(*ClientFirestore)
	return sampleClientInstance
}

func sampleCoder(client db.Client) db.Coder {
	return client.Coder("SampleCode", 10)
}

func sampleCoder1000(client db.Client) db.Coder {
	return client.Coder("SampleCode1000", 1000)
}

func sampleSerial(client db.Client) db.Serial {
	return client.Serial("SampleSerial")
}

func sampleCounter(client db.Client) db.Counter {
	return client.Counter("SampleCount", 3)
}

func sampleCounter1000(client db.Client) db.Counter {
	return client.Counter("SampleCount", 1000)
}
