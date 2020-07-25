package util

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestOrderNumber(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	Convey("should generate Order Number string", t, func() {
		num := OrderNumber()
		str := OrderNumberToString(num)
		valid := OrderNumberIsValid(str)
		So(valid, ShouldBeTrue)
		So(str, ShouldNotBeEmpty)
		retNum, err := OrderNumberFromString(str)
		So(err, ShouldBeNil)
		So(retNum, ShouldEqual, num)
		fmt.Printf("%v\n", str)

		retNum, err = OrderNumberFromString("a")
		So(err, ShouldNotBeNil)
		So(retNum, ShouldEqual, 0)
		valid = OrderNumberIsValid("aaaa")
		So(valid, ShouldBeFalse)
		valid = OrderNumberIsValid("0725-1726-4071-2412")
		So(valid, ShouldBeFalse)
	})

	Convey("should generate Order Number", t, func() {
		num := OrderNumber()
		So(num > 0, ShouldBeTrue)
		//fmt.Printf("%v\n", num)
	})

}

func BenchmarkOrderNumber(b *testing.B) {
	for i := 0; i < 10000; i++ {
		OrderNumber()
	}
}
