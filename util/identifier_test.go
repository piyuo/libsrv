package util

import (
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

func TestOrderNumber(t *testing.T) {
	Convey("should generate Order Number", t, func() {
		id := OrderNumber()
		So(id, ShouldNotBeEmpty)
		So(len(id), ShouldEqual, 19)
		//fmt.Printf("%v, %v\n", id, len(id))
	})
}

func BenchmarkOrderNumber(b *testing.B) {
	for i := 0; i < 10000; i++ {
		OrderNumber()
	}
}
