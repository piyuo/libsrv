package cache

import (
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestCacheLimit(t *testing.T) {
	convey.Convey("should only cache 1000 low frequency item", t, func() {
		Reset()
		for i := 0; i <= 1002; i++ {
			Set(LOW, "l"+strconv.Itoa(i), "1")
		}
		convey.So(Count(), convey.ShouldEqual, 1000)
		Reset()
	})
	convey.Convey("should only cache 2000 medium frequency item", t, func() {
		Reset()
		for i := 0; i <= 2002; i++ {
			Set(MEDIUM, "m"+strconv.Itoa(i), "1")
		}
		convey.So(Count(), convey.ShouldEqual, 2000)
		Reset()
	})
	convey.Convey("should only cache 2500 high frequency item", t, func() {
		Reset()
		for i := 0; i <= 2530; i++ {
			Set(HIGH, "m"+strconv.Itoa(i), "1")
		}
		convey.So(Count(), convey.ShouldEqual, 2500)
		Reset()
	})
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

func TestGetMethod(t *testing.T) {
	convey.Convey("should set and get", t, func() {
		Reset()
		Set(HIGH, "key", 1)
		valueInt, found := GetInt("key")
		convey.So(found, convey.ShouldBeTrue)
		convey.So(valueInt, convey.ShouldEqual, 1)

		Set(HIGH, "key", int64(2))
		valueInt64, found := GetInt64("key")
		convey.So(found, convey.ShouldBeTrue)
		convey.So(valueInt64, convey.ShouldEqual, 2)

		Set(HIGH, "key", true)
		valueBool, found := GetBool("key")
		convey.So(found, convey.ShouldBeTrue)
		convey.So(valueBool, convey.ShouldEqual, true)

		Set(HIGH, "key", "hi")
		valueString, found := GetString("key")
		convey.So(found, convey.ShouldBeTrue)
		convey.So(valueString, convey.ShouldEqual, "hi")

		Set(HIGH, "key", []byte("hi"))
		valueBytes, found := GetBytes("key")
		convey.So(found, convey.ShouldBeTrue)
		convey.So(valueBytes, convey.ShouldNotBeNil)

		//test not exist
		valueInt, found = GetInt("not-exist")
		convey.So(valueInt, convey.ShouldEqual, 0)
		convey.So(found, convey.ShouldBeFalse)

		valueInt64, found = GetInt64("not-exist")
		convey.So(valueInt64, convey.ShouldEqual, 0)
		convey.So(found, convey.ShouldBeFalse)

		valueBool, found = GetBool("not-exist")
		convey.So(valueBool, convey.ShouldBeFalse)
		convey.So(found, convey.ShouldBeFalse)

		valueString, found = GetString("not-exist")
		convey.So(valueString, convey.ShouldBeEmpty)
		convey.So(found, convey.ShouldBeFalse)

		valueBytes, found = GetBytes("not-exist")
		convey.So(valueBytes, convey.ShouldBeNil)
		convey.So(found, convey.ShouldBeFalse)
	})
}

func TestCache(t *testing.T) {
	convey.Convey("should set and get", t, func() {
		Reset()
		Set(HIGH, "key-1", "1")
		Set(HIGH, "key-2", "2")

		// wait for value to pass through buffers
		time.Sleep(10 * time.Millisecond)

		value, found := Get("key-1")
		convey.So(found, convey.ShouldBeTrue)
		convey.So(value, convey.ShouldEqual, "1")
		value, found = Get("key-2")
		convey.So(found, convey.ShouldBeTrue)
		convey.So(value, convey.ShouldEqual, "2")
		if found {
		}

		Delete("key-1")
		value, found = Get("key-1")
		convey.So(found, convey.ShouldBeFalse)
		convey.So(value, convey.ShouldBeNil)

		value, found = Get("not exist")
		convey.So(found, convey.ShouldBeFalse)
		convey.So(value, convey.ShouldBeNil)
	})
}

func TestExpireCache(t *testing.T) {
	convey.Convey("should expire", t, func() {
		Reset()
		set("key", "1", 50*time.Millisecond)
		value, found := Get("key")
		convey.So(found, convey.ShouldBeTrue)
		convey.So(value, convey.ShouldEqual, "1")

		// wait for value to pass through buffers
		time.Sleep(51 * time.Millisecond)
		value, found = Get("key")
		convey.So(found, convey.ShouldBeFalse)
		convey.So(value, convey.ShouldBeNil)

	})
}

func TestCachePurges(t *testing.T) {
	convey.Convey("should purge expired item", t, func() {
		configCache(50*time.Millisecond, 50*time.Millisecond)
		defer configCache(10*time.Minute, 3*time.Minute)

		set("key", "1", 1500*time.Millisecond)
		value, found := Get("key")
		convey.So(found, convey.ShouldBeTrue)
		convey.So(value, convey.ShouldEqual, "1")

		set("key2", "2", 50*time.Millisecond)
		value2, found2 := Get("key2")
		convey.So(found2, convey.ShouldBeTrue)
		convey.So(value2, convey.ShouldEqual, "2")

		// wait for value to pass through buffers
		time.Sleep(51 * time.Millisecond)

		value, found = Get("key")
		convey.So(found, convey.ShouldBeTrue)
		convey.So(value, convey.ShouldEqual, "1")

		value2, found2 = Get("key2")
		convey.So(found2, convey.ShouldBeFalse)
		convey.So(value2, convey.ShouldBeNil)

	})
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
