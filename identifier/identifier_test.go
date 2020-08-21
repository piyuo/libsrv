package identifier

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUUID(t *testing.T) {
	Convey("should generate uuid", t, func() {
		id := UUID()
		So(id, ShouldNotBeEmpty)
		//fmt.Printf("%v, %v\n", id, len(id))
	})
}

func TestSerialID16(t *testing.T) {
	Convey("should generate base58 id from uint16", t, func() {
		id := SerialID16(uint16(42))
		So(id, ShouldNotBeEmpty)
		//fmt.Printf("%v, %v\n", id, len(id))
	})
}

func TestSerialID32(t *testing.T) {
	Convey("should generate base58 id from uint32", t, func() {
		//		for i := 0; i < 10000; i++ {
		id := SerialID32(uint32(42))
		So(id, ShouldNotBeEmpty)
		//		fmt.Printf("%v, %v\n", id, len(id))
		//		}
	})
}

func TestSerialID64(t *testing.T) {
	Convey("should generate base58 id from uint64", t, func() {
		//for i := 0; i < 10000; i++ {
		id := SerialID64(uint64(42))
		So(id, ShouldNotBeEmpty)
		//fmt.Printf("%v, %v\n", id, len(id))
		//}
	})
}

func TestRandomNumber(t *testing.T) {
	Convey("should generate 6 digit random number string", t, func() {
		id := RandomNumber(6)
		So(id, ShouldNotBeEmpty)
		So(len(id), ShouldEqual, 6)
		fmt.Printf("%v, %v\n", id, len(id))
	})
}

func TestIdenticalNumberString(t *testing.T) {
	Convey("should be Identical", t, func() {
		So(IsNumberStringIdentical("111111"), ShouldBeTrue)
		So(IsNumberStringIdentical("111122"), ShouldBeTrue)
	})
	Convey("should Not be Identical", t, func() {
		So(IsNumberStringIdentical("111124"), ShouldBeFalse)
		So(IsNumberStringIdentical("123456"), ShouldBeFalse)
		So(IsNumberStringIdentical("177756"), ShouldBeFalse)
		So(IsNumberStringIdentical("211311"), ShouldBeFalse)
		So(IsNumberStringIdentical("111311"), ShouldBeFalse)
	})
	Convey("should Not Identical random", t, func() {
		So(NotIdenticalRandomNumber(6), ShouldNotBeEmpty)
	})
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
	Convey("should generate id for map", t, func() {
		m := map[string]string{}
		id, err := MapID(m)
		So(err, ShouldBeNil)
		So(id, ShouldEqual, "1")
		m[id] = "a"

		id, err = MapID(m)
		So(err, ShouldBeNil)
		So(id, ShouldEqual, "2")
		m[id] = "b"

		id, err = MapID(m)
		So(err, ShouldBeNil)
		So(id, ShouldEqual, "3")
		m[id] = "c"
	})
}

func BenchmarkMapID(b *testing.B) {
	m := map[string]string{}
	for i := 0; i < 100; i++ {
		id, _ := MapID(m)
		m[id] = id
	}
}
