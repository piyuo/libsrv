package file

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCache(t *testing.T) {

	Convey("should set and get'", t, func() {
		Set("hi-1", "1", 0)
		Set("hi-2", "2", 0)

		// wait for value to pass through buffers
		time.Sleep(10 * time.Millisecond)

		value, found := Get("hi-1")
		So(found, ShouldBeTrue)
		if found {
			So(value, ShouldEqual, "1")
		}
		value, found = Get("hi-2")
		So(found, ShouldBeTrue)
		So(value, ShouldEqual, "2")
		if found {
		}

		value, found = Get("not exist")
		So(found, ShouldBeFalse)
		So(value, ShouldBeNil)
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
