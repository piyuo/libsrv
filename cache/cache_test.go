package cache

import (
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCacheLimit(t *testing.T) {
	assert := assert.New(t)
	//should only cache 1000 low frequency item
	Reset()
	for i := 0; i <= 1002; i++ {
		Set(LOW, "l"+strconv.Itoa(i), "1")
	}
	assert.Equal(1000, Count())
	Reset()
	//should only cache 2000 medium frequency item
	Reset()
	for i := 0; i <= 2002; i++ {
		Set(MEDIUM, "m"+strconv.Itoa(i), "1")
	}
	assert.Equal(2000, Count())
	Reset()
	//should only cache 2500 high frequency item
	Reset()
	for i := 0; i <= 2530; i++ {
		Set(HIGH, "m"+strconv.Itoa(i), "1")
	}
	assert.Equal(2500, Count())
	Reset()
}

func TestConcurrentCache(t *testing.T) {
	var concurrent = 5
	var wg sync.WaitGroup
	wg.Add(concurrent)
	runCache := func() {
		for i := 0; i < 20; i++ {
			Set(HIGH, "key"+strconv.Itoa(i), i)
			value, found := Get("key" + strconv.Itoa(i))
			if !found {
				t.Fatal("key" + strconv.Itoa(i) + " already set to cache, but not found in cache")
				return
			}
			if value != i {
				t.Fatal("key" + strconv.Itoa(i) + " get value is not equal to set value")
				return
			}
			//fmt.Print(strconv.Itoa(i) + "\n")

		}
		wg.Done()
	}

	//create go routing to do counting
	for i := 0; i < concurrent; i++ {
		go runCache()
	}
	wg.Wait()
}

func TestGetSetMethod(t *testing.T) {
	assert := assert.New(t)
	Reset()
	Set(HIGH, "key", 1)
	valueInt, found := GetInt("key")
	assert.True(found)
	assert.Equal(1, valueInt)

	Set(HIGH, "key32", uint32(1))
	valueUInt32, found := GetUInt32("key32")
	assert.True(found)
	assert.Equal(uint32(1), valueUInt32)

	Set(HIGH, "key64", uint64(1))
	valueUInt64, found := GetUInt64("key64")
	assert.True(found)
	assert.Equal(uint64(1), valueUInt64)

	Set(HIGH, "key", int64(2))
	valueInt64, found := GetInt64("key")
	assert.True(found)
	assert.Equal(int64(2), valueInt64)

	Set(HIGH, "key", true)
	valueBool, found := GetBool("key")
	assert.True(found)
	assert.True(valueBool)

	Set(HIGH, "key", "hi")
	valueString, found := GetString("key")
	assert.True(found)
	assert.Equal("hi", valueString)

	Set(HIGH, "key", []byte("hi"))
	valueBytes, found := GetBytes("key")
	assert.True(found)
	assert.NotNil(valueBytes)

	//test not exist
	valueInt, found = GetInt("not-exist")
	assert.Equal(0, valueInt)
	assert.False(found)

	valueUInt32, found = GetUInt32("not-exist")
	assert.Equal(uint32(0), valueUInt32)
	assert.False(found)

	valueUInt64, found = GetUInt64("not-exist")
	assert.Equal(uint64(0), valueUInt64)
	assert.False(found)

	valueInt64, found = GetInt64("not-exist")
	assert.Equal(int64(0), valueInt64)
	assert.False(found)

	valueBool, found = GetBool("not-exist")
	assert.False(valueBool)
	assert.False(found)

	valueString, found = GetString("not-exist")
	assert.Empty(valueString)
	assert.False(found)

	valueBytes, found = GetBytes("not-exist")
	assert.Nil(valueBytes)
	assert.False(found)
}

func TestCache(t *testing.T) {
	assert := assert.New(t)
	Reset()
	Set(HIGH, "key-1", "1")
	Set(HIGH, "key-2", "2")

	// wait for value to pass through buffers
	time.Sleep(10 * time.Millisecond)

	value, found := Get("key-1")
	assert.True(found)
	assert.Equal("1", value)
	value, found = Get("key-2")
	assert.True(found)
	assert.Equal("2", value)

	Delete("key-1")
	value, found = Get("key-1")
	assert.False(found)
	assert.Nil(value)

	value, found = Get("not exist")
	assert.False(found)
	assert.Nil(value)
}

func TestExpireCache(t *testing.T) {
	assert := assert.New(t)
	Reset()
	set("key", "1", 50*time.Millisecond)
	value, found := Get("key")
	assert.True(found)
	assert.Equal("1", value)

	// wait for value to pass through buffers
	time.Sleep(51 * time.Millisecond)
	value, found = Get("key")
	assert.False(found)
	assert.Nil(value)
}

func TestCachePurges(t *testing.T) {
	assert := assert.New(t)
	configCache(50*time.Millisecond, 50*time.Millisecond)
	defer configCache(10*time.Minute, 3*time.Minute)

	set("key", "1", 1500*time.Millisecond)
	value, found := Get("key")
	assert.True(found)
	assert.Equal("1", value)

	set("key2", "2", 50*time.Millisecond)
	value2, found2 := Get("key2")
	assert.True(found2)
	assert.Equal("2", value2)

	// wait for value to pass through buffers
	time.Sleep(51 * time.Millisecond)

	value, found = Get("key")
	assert.True(found)
	assert.Equal("1", value)

	value2, found2 = Get("key2")
	assert.False(found2)
	assert.Nil(value2)
}

func TestIncrement(t *testing.T) {
	assert := assert.New(t)
	Reset()
	key := "key"
	Increment(LOW, key, 2)

	value, found := Get(key)
	assert.True(found)
	assert.Equal(2, value)

	Increment(LOW, key, -1)

	value, found = Get("key")
	assert.True(found)
	assert.Equal(1, value)
}

func BenchmarkGoCache(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for x := 0; x < 100; x++ {
			Set(HIGH, "key-"+strconv.Itoa(x), i)
		}
	}

	// wait for value to pass through buffers
	time.Sleep(10 * time.Millisecond)

	for i := 0; i < b.N; i++ {
		for x := 0; x < 100; x++ {
			Get("key-" + strconv.Itoa(x))
		}
	}
}
