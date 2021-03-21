package gdb

import (
	"context"

	"github.com/piyuo/libsrv/src/db"
	"github.com/piyuo/libsrv/src/google/gaccount"
)

type PlainObject struct {
	ID   string
	Name string
}

// SampleNoFactory with no factory and collection
//
type SampleNoFactory struct {
	db.DomainObject
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

// SampleClear for test clear
//
type SampleClear struct {
	db.DomainObject
	Name string `firestore:"Name,omitempty"`
}

func (c *SampleClear) Factory() db.Object {
	return &SampleClear{}
}

func (c *SampleClear) Collection() string {
	return "SampleClear"
}

// Sample for test
//
type Sample struct {
	db.DomainObject
	Name    string            `firestore:"Name,omitempty"`
	Tag     string            `firestore:"Tag,omitempty"`
	Value   int               `firestore:"Value,omitempty"`
	Map     map[string]string `firestore:"Map,omitempty"`
	Array   []string          `firestore:"Array,omitempty"`
	Numbers []int             `firestore:"Numbers,omitempty"`
	PObj    *PlainObject      `firestore:"PObj,omitempty"`
}

func (c *Sample) Factory() db.Object {
	return &Sample{}
}

func (c *Sample) Collection() string {
	return "Sample"
}

var sampleClientInstance *ClientFirestore

// sample client create db client use for test, it will keep client instance to resuse, recreate new instance if client is close
//
func sampleClient() db.Client {
	if sampleClientInstance != nil && sampleClientInstance.firestoreClient != nil {
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
	return client.Counter("SampleCount", 3, db.DateHierarchyNone)
}

func sampleCounter1000(client db.Client) db.Counter {
	return client.Counter("SampleCount", 1000, db.DateHierarchyNone)
}
