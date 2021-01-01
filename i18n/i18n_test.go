package i18n

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestIsPredefined(t *testing.T) {
	Convey("should return predefined locale", t, func() {

		exist, predefine := IsPredefined("en-us")
		So(exist, ShouldBeTrue)
		So(predefine, ShouldEqual, "en_US")

		exist, predefine = IsPredefined("zh-tw")
		So(exist, ShouldBeTrue)
		So(predefine, ShouldEqual, "zh_TW")

		exist, predefine = IsPredefined("en")
		So(exist, ShouldBeFalse)
		So(predefine, ShouldBeEmpty)

	})
}
