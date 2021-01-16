package identifier

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUUID(t *testing.T) {
	assert := assert.New(t)
	id := UUID()
	assert.NotEmpty(id)
	//fmt.Printf("%v, %v\n", id, len(id))
}

func TestSerialID16(t *testing.T) {
	assert := assert.New(t)
	id := SerialID16(uint16(42))
	assert.NotEmpty(id)
	//fmt.Printf("%v, %v\n", id, len(id))
}

func TestSerialID32(t *testing.T) {
	assert := assert.New(t)
	//		for i := 0; i < 10000; i++ {
	id := SerialID32(uint32(42))
	assert.NotEmpty(id)
	//		fmt.Printf("%v, %v\n", id, len(id))
	//		}
}

func TestSerialID64(t *testing.T) {
	assert := assert.New(t)
	//for i := 0; i < 10000; i++ {
	id := SerialID64(uint64(42))
	assert.NotEmpty(id)
	//fmt.Printf("%v, %v\n", id, len(id))
	//}
}

func TestRandomNumber(t *testing.T) {
	assert := assert.New(t)
	id := RandomNumber(6)
	assert.NotEmpty(id)
	assert.Equal(6, len(id))
	fmt.Printf("%v, %v\n", id, len(id))
}

func TestIdenticalNumberString(t *testing.T) {
	assert := assert.New(t)
	//should be Identical
	assert.True(IsNumberStringIdentical("111111"))
	assert.True(IsNumberStringIdentical("111122"))
	//should Not be Identical
	assert.False(IsNumberStringIdentical("111124"))
	assert.False(IsNumberStringIdentical("123456"))
	assert.False(IsNumberStringIdentical("177756"))
	assert.False(IsNumberStringIdentical("211311"))
	assert.False(IsNumberStringIdentical("111311"))
	//should Not Identical random
	assert.NotEmpty(NotIdenticalRandomNumber(6))
}

func BenchmarkRandomNumber(b *testing.B) {
	for i := 0; i < 10000; i++ {
		RandomNumber(6)
	}
}

func BenchmarkNotIdenticalRandomNumber(b *testing.B) {
	for i := 0; i < 10000; i++ {
		NotIdenticalRandomNumber(6)
	}
}

func TestMapID(t *testing.T) {
	assert := assert.New(t)
	m := map[string]string{}
	id, err := MapID(m)
	assert.Nil(err)
	assert.Equal("1", id)
	m[id] = "a"

	id, err = MapID(m)
	assert.Nil(err)
	assert.Equal("2", id)
	m[id] = "b"

	id, err = MapID(m)
	assert.Nil(err)
	assert.Equal("3", id)
	m[id] = "c"

}

func BenchmarkMapID(b *testing.B) {
	m := map[string]string{}
	for i := 0; i < 100; i++ {
		id, _ := MapID(m)
		m[id] = id
	}
}
