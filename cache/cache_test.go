package file

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"

	convey "github.com/smartystreets/goconvey/convey"
)

func TestCache(t *testing.T) {

	convey.Convey("should set and get", t, func() {
		Reset()
		Set("key-1", "1", 0)
		Set("key-2", "2", 0)

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
		Set("key", "1", 50*time.Millisecond)
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

func TestDefaultExpire(t *testing.T) {
	convey.Convey("should expire", t, func() {
		configCache(100*time.Millisecond, 150*time.Millisecond)
		Set("never-expire", "1", -1)
		Set("default-expire", "2", 0)
		value, found := Get("never-expire")
		convey.So(found, convey.ShouldBeTrue)
		convey.So(value, convey.ShouldEqual, "1")
		value, found = Get("default-expire")
		convey.So(found, convey.ShouldBeTrue)
		convey.So(value, convey.ShouldEqual, "2")

		// wait for value to pass through buffers
		time.Sleep(151 * time.Millisecond)

		value, found = Get("never-expire")
		convey.So(found, convey.ShouldBeTrue)
		convey.So(value, convey.ShouldEqual, "1")
		value, found = Get("default-expire")
		convey.So(found, convey.ShouldBeFalse)
		convey.So(value, convey.ShouldBeNil)
	})
}

func TestConcurrentCache(t *testing.T) {
	var concurrent = 20
	var wg sync.WaitGroup
	wg.Add(concurrent)
	runCache := func() {
		for i := 0; i < 100; i++ {
			Set("key"+strconv.Itoa(i), i, 0)
			value, found := Get("key" + strconv.Itoa(i))
			if !found {
				t.Fatal("key" + strconv.Itoa(i) + " already set to cache, but can not get")
			}
			if value != i {
				t.Fatal("key" + strconv.Itoa(i) + " get value is not equal to set value")
			}
			fmt.Print(strconv.Itoa(i) + "\n")

		}
		wg.Done()
	}

	//create go routing to do counting
	for i := 0; i < concurrent; i++ {
		go runCache()
	}
	wg.Wait()
}

func BenchmarkGoCache(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for x := 0; x < 100; x++ {
			Set("key-"+strconv.Itoa(x), i, 0)
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
