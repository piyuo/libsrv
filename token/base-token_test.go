package token

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestExpiredToken(t *testing.T) {
	Convey("should check expired date", t, func() {
		t := time.Now().UTC().Add(60 * time.Second)
		txt := t.Format(expiredFormat)

		expired := isExpired(txt)
		So(expired, ShouldBeFalse)

		expired = isExpired("200001010101")
		So(expired, ShouldBeTrue)

		expired = isExpired("300001010101")
		So(expired, ShouldBeFalse)
	})
}

func TestTokenGetSetDelete(t *testing.T) {
	Convey("should get/set token", t, func() {
		token := NewToken()

		value := token.Get("a")
		So(value, ShouldEqual, "")

		token.Set("a", "1")
		value = token.Get("a")
		So(value, ShouldEqual, "1")

		token.Delete("a")
		value = token.Get("a")
		So(value, ShouldEqual, "")
	})
}

func TestTokenFromToString(t *testing.T) {
	Convey("should get/set token", t, func() {
		token := NewToken()
		token.Set("a", "1")

		crypted, err := token.ToString(60 * time.Second)
		So(err, ShouldBeNil)
		So(crypted, ShouldNotBeEmpty)

		token2, expired, err := FromString(crypted)
		So(err, ShouldBeNil)
		So(expired, ShouldBeFalse)

		value := token2.Get("a")
		So(value, ShouldEqual, "1")
	})
}

func TestTokenExpired(t *testing.T) {
	Convey("should expired", t, func() {
		token := NewToken()
		token.Set("a", "1")

		crypted, err := token.ToString(-60 * time.Second)
		So(err, ShouldBeNil)
		So(crypted, ShouldNotBeEmpty)

		token2, expired, err := FromString(crypted)
		So(err, ShouldBeNil)
		So(expired, ShouldBeTrue)
		So(token2, ShouldBeNil)
	})
}

func TestInvalidToken(t *testing.T) {
	Convey("should return error", t, func() {
		token, expired, err := FromString("")
		So(err, ShouldNotBeNil)
		So(expired, ShouldBeFalse)
		So(token, ShouldBeNil)

		token, expired, err = FromString("123213123")
		So(err, ShouldNotBeNil)
		So(expired, ShouldBeFalse)
		So(token, ShouldBeNil)
	})
}

func TestIsExpired(t *testing.T) {
	Convey("should expired when str is invalid", t, func() {
		result := isExpired("a")
		So(result, ShouldBeTrue)

	})
}
