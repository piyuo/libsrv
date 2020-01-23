package libsrv

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTokenContext(t *testing.T) {
	token := NewToken("userId", "storeId", "locationId", "permission", time.Now(), time.Now())
	Convey("should save to and load from string'", t, func() {

		everything := token.ToString()
		So(everything, ShouldNotBeEmpty)
		So(strings.Contains(everything, "userId"), ShouldBeTrue)
		So(strings.Contains(everything, "storeId"), ShouldBeTrue)
		So(strings.Contains(everything, "locationId"), ShouldBeTrue)

		tokenNew, err := TokenFromString(everything)
		everythingNew := tokenNew.ToString()
		So(err, ShouldBeNil)
		So(everythingNew, ShouldNotBeEmpty)
		So(strings.Contains(everythingNew, "userId"), ShouldBeTrue)
		So(strings.Contains(everythingNew, "storeId"), ShouldBeTrue)
		So(strings.Contains(everythingNew, "locationId"), ShouldBeTrue)
	})

	Convey("should save to and load from context'", t, func() {
		ctx := context.Background()
		ctx2 := token.ToContext(ctx)
		So(ctx, ShouldNotBeNil)

		token2, err := TokenFromContext(ctx2)
		So(err, ShouldBeNil)
		So(token2.UserID(), ShouldEqual, "userId")
		So(token2.StoreID(), ShouldEqual, "storeId")
		So(token2.LocationID(), ShouldEqual, "locationId")
	})
}

func saveCookie(w http.ResponseWriter, r *http.Request) {
	token := NewToken("userId", "storeId", "locationId", "permission", time.Now(), time.Now())
	err := token.ToCookie(w)
	So(err, ShouldBeNil)
}

func loadCookie(w http.ResponseWriter, r *http.Request) {
	token, err := TokenFromCookie(r)
	So(err, ShouldBeNil)
	So(token.UserID(), ShouldEqual, "userId")

}

func TestCookie(t *testing.T) {
	cookieValueBackup := ""
	Convey("should save to and load from cookie'", t, func() {
		request, _ := http.NewRequest("GET", "/", bytes.NewReader([]byte("")))
		responseRecord := httptest.NewRecorder()
		handler := http.HandlerFunc(saveCookie)
		handler.ServeHTTP(responseRecord, request)
		response := responseRecord.Result()
		So(response, ShouldNotBeNil)
		cookie := response.Cookies()[0]
		cookieValueBackup = cookie.Value
		So(cookie.Name, ShouldEqual, "piyuo")
		So(cookie.Value, ShouldNotBeEmpty)
	})

	Convey("should save to and load from cookie'", t, func() {
		request, _ := http.NewRequest("GET", "/", bytes.NewReader([]byte("")))
		request.AddCookie(&http.Cookie{Name: string(CookieTokenName), Value: cookieValueBackup})
		responseRecord := httptest.NewRecorder()
		handler := http.HandlerFunc(loadCookie)
		handler.ServeHTTP(responseRecord, request)
		response := responseRecord.Result()
		So(response, ShouldNotBeNil)
	})

}

func TestTime(t *testing.T) {
	Convey("should expire", t, func() {
		token := &token{
			created: time.Now(),
		}
		So(token.Expired(), ShouldBeFalse)
		token.created = time.Now().Add(-61 * time.Minute) //create by 61 minutes ago
		So(token.Expired(), ShouldBeTrue)
	})

	Convey("should revive", t, func() {
		begin := time.Now()
		token := &token{
			created: begin,
		}
		So(token.Revive(), ShouldBeFalse)
		token.created = time.Now().Add(-9 * time.Minute) //create by 9 minutes ago
		So(token.Revive(), ShouldBeFalse)
		token.created = time.Now().Add(-15 * time.Minute) //create by 15 minutes ago
		So(token.Revive(), ShouldBeTrue)
		So(token.created.Equal(begin) || token.created.After(begin), ShouldBeTrue)
		token.created = time.Now().Add(-61 * time.Minute) //create by 15 minutes ago
		So(token.Revive(), ShouldBeFalse)
	})

	Convey("should find is user just login", t, func() {
		token := &token{}
		So(token.IsUserJustLogin(), ShouldBeFalse)

		token.login = time.Now().Add(-4 * time.Minute) //login by  4 minutes ago
		So(token.IsUserJustLogin(), ShouldBeFalse)

		token.login = time.Now().Add(-2 * time.Minute) //login by  4 minutes ago
		So(token.IsUserJustLogin(), ShouldBeTrue)
	})

}
