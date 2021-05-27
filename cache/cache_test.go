package cache

import (
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/piyuo/libsrv/identifier"
	"github.com/stretchr/testify/assert"
)

func TestConcurrentCache(t *testing.T) {
	t.Parallel()
	var concurrent = 5
	var wg sync.WaitGroup
	wg.Add(concurrent)
	runCache := func() {
		for i := 0; i < 10; i++ {
			key := "concurrent-" + identifier.RandomNumber(6) + "-"
			Set(key+strconv.Itoa(i), []byte(strconv.Itoa(i)), 0)
			found, value, err := Get(key + strconv.Itoa(i))
			if err != nil {
				t.Fatal(key + strconv.Itoa(i) + " failed to get value")
				return
			}
			if !found {
				t.Fatal(key + strconv.Itoa(i) + " already set to cache, but not found in cache")
				return
			}
			s, err := strconv.Atoi(string(value))
			if err != nil {
				t.Fatal(key + strconv.Itoa(i) + " failed to convert value to int")
				return
			}
			if s != i {
				t.Fatal(key + strconv.Itoa(i) + " get value is not equal to set value")
				return
			}
		}
		wg.Done()
	}
	for i := 0; i < concurrent; i++ {
		go runCache()
	}
	wg.Wait()
}

func TestSet(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	key := "set-" + identifier.RandomNumber(6)
	err := Set(key, []byte("hi"), 0)
	assert.Nil(err)
	found, valueBytes, err := Get(key)
	assert.Nil(err)
	value := string(valueBytes)
	assert.True(found)
	assert.Equal("hi", value)

	err = SetString(key, "hi str", 0)
	assert.Nil(err)
	found, str, err := GetString(key)
	assert.Nil(err)
	assert.True(found)
	assert.Equal("hi str", str)

	err = SetInt64(key, 65536, 0)
	assert.Nil(err)
	found, valueInt64, err := GetInt64(key)
	assert.Nil(err)
	assert.True(found)
	assert.Equal(int64(65536), valueInt64)

	found, valueInt64, err = GetInt64("not-exists")
	assert.Nil(err)
	assert.False(found)
	assert.Equal(int64(-1), valueInt64)
}

func TestSetEmptyBytes(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	key := "empty-bytes-" + identifier.RandomNumber(6)
	err := Set(key, []byte{}, 0)
	assert.Nil(err)
	//set empty bytes will return nil
	found, got, err := Get(key)
	assert.Nil(err)
	assert.True(found)
	assert.Nil(got)
}

func TestNoDuration(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	key := "no-duration-" + identifier.RandomNumber(6)
	err := Set(key, []byte("hi"), -1)
	assert.Nil(err)
	found, valueBytes, err := Get(key)
	assert.Nil(err)
	assert.False(found)
	assert.Nil(valueBytes)

	err = ISet(123, []byte("hi"), -1)
	assert.Nil(err)
	found, valueBytes, err = IGet(123)
	assert.Nil(err)
	assert.False(found)
	assert.Nil(valueBytes)
}

func TestInt(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	err := ISet(123, []byte("hi"), 0)
	assert.Nil(err)
	found, valueBytes, err := IGet(123)
	assert.Nil(err)
	value := string(valueBytes)
	assert.True(found)
	assert.Equal("hi", value)
	IDelete(123)
	found, valueBytes, err = IGet(123)
	assert.Nil(err)
	assert.Nil(valueBytes)
	assert.False(found)
}

func TestDelete(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	key := "del-" + identifier.RandomNumber(6) + "-"
	count := Count()
	err := SetString(key+"1", "1", 0)
	assert.Nil(err)
	err = SetString(key+"2", "2", 0)
	assert.Nil(err)
	count2 := Count()
	assert.True(count2 > count)

	// wait for value to pass through buffers
	time.Sleep(10 * time.Millisecond)

	found, value, err := GetString(key + "1")
	assert.Nil(err)
	assert.True(found)
	assert.Equal("1", value)
	found, value, err = GetString(key + "2")
	assert.Nil(err)
	assert.True(found)
	assert.Equal("2", value)

	Delete(key + "1")
	found, value, err = GetString(key + "1")
	assert.Nil(err)
	assert.False(found)
	assert.Empty(value)

	Delete(key + "2")
	found, value, err = GetString(key + "2")
	assert.Nil(err)
	assert.False(found)
	assert.Empty(value)
}

func TestExpired(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	key := "expired-" + identifier.RandomNumber(6)
	err := SetString(key, "1", 1*time.Second)
	assert.Nil(err)

	found, value, err := GetString(key)
	assert.Nil(err)
	assert.True(found)
	assert.Equal("1", value)

	// wait for value to pass through buffers
	time.Sleep(2 * time.Second)
	found, value, err = GetString(key)
	assert.Nil(err)
	assert.False(found)
	assert.Empty(value)
}

func TestGzip(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	key := "zip-key-" + identifier.RandomNumber(6)
	value := "zip-value-" + identifier.RandomNumber(90)
	err := GzipSet(key, []byte(value), 0)
	assert.Nil(err)

	found, valueBytes, err := GzipGet(key)
	assert.Nil(err)
	value2 := string(valueBytes)
	assert.True(found)
	assert.Equal(value, value2)
}

func TestGzipEmptyArray(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	key := "zip-empty-" + identifier.RandomNumber(6)
	err := GzipSet(key, []byte{}, 0)
	assert.Nil(err)

	found, valueBytes, err := GzipGet(key)
	assert.Nil(err)
	assert.True(found)
	assert.Empty(valueBytes)
}

func BenchmarkCache(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for x := 0; x < 100; x++ {
			SetInt64("benchmark-"+strconv.Itoa(x), int64(i), 0)
		}
	}

	// wait for value to pass through buffers
	time.Sleep(10 * time.Millisecond)

	for i := 0; i < b.N; i++ {
		for x := 0; x < 100; x++ {
			GetInt64("benchmark-" + strconv.Itoa(x))
		}
	}
}
