package mail

import (
	"context"
	"net/http"
	"testing"

	"github.com/piyuo/libsrv/session"
	. "github.com/smartystreets/goconvey/convey"
)

func TestMail(t *testing.T) {
	Convey("should set/get method", t, func() {
		req, _ := http.NewRequest("GET", "/", nil)
		req.Header.Add("Accept-Language", "en_US")
		ctx := context.WithValue(context.Background(), session.KeyRequest, req)

		mail, err := NewMail(ctx, "mock-mail")
		So(err, ShouldBeNil)
		So(mail.GetSubject(), ShouldNotBeEmpty)
		backupSubject := mail.GetSubject()

		mail.SetSubject("ok")
		So(mail.GetSubject(), ShouldEqual, "ok")
		mail.ReplaceSubject("ok", "1")
		So(mail.GetSubject(), ShouldEqual, "1")

		So(mail.GetText(), ShouldNotBeEmpty)
		mail.SetText("ok")
		So(mail.GetText(), ShouldEqual, "ok")
		mail.ReplaceText("ok", "1")
		So(mail.GetText(), ShouldEqual, "1")

		So(mail.GetHTML(), ShouldNotBeEmpty)
		mail.SetHTML("ok")
		So(mail.GetHTML(), ShouldEqual, "ok")
		mail.ReplaceHTML("ok", "1")
		So(mail.GetHTML(), ShouldEqual, "1")

		name, address := mail.GetFrom()
		So(name, ShouldNotBeEmpty)
		So(address, ShouldNotBeEmpty)
		mail.SetFrom("1", "2")
		name, address = mail.GetFrom()
		So(name, ShouldEqual, "1")
		So(address, ShouldEqual, "2")

		So(mail.GetTo(), ShouldBeNil)
		mail.AddTo("name", "address")
		So(mail.GetTo(), ShouldNotBeNil)
		So(len(mail.GetTo()), ShouldEqual, 1)
		mail.AddTo("name1", "address1")
		So(len(mail.GetTo()), ShouldEqual, 2)
		mail.ResetTo()
		So(mail.GetTo(), ShouldBeNil)

		//should from cache
		mail, err = NewMail(ctx, "mock-mail")
		So(err, ShouldBeNil)
		So(mail, ShouldNotBeNil)
		So(mail.GetSubject(), ShouldEqual, backupSubject)

		mail, err = NewMail(ctx, "notExist")
		So(err, ShouldNotBeNil)
		So(mail, ShouldBeNil)
	})
}

func TestSendMail(t *testing.T) {
	Convey("should send mail", t, func() {
		req, _ := http.NewRequest("GET", "/", nil)
		req.Header.Add("Accept-Language", "en_US")
		ctx := context.WithValue(context.Background(), session.KeyRequest, req)
		mail, err := NewMail(ctx, "mock-mail")
		So(err, ShouldBeNil)
		mail.AddTo("p", "a@b.c")
		mail.ReplaceText("%1", "1234")
		mail.ReplaceHTML("%1", "1234")
		err = mail.Send(ctx)
		So(err, ShouldBeNil)
	})
}

func TestMockSendMail(t *testing.T) {
	Convey("should mock send mail", t, func() {
		So(mockResult, ShouldBeNil)
		Mock(true)
		req, _ := http.NewRequest("GET", "/", nil)
		req.Header.Add("Accept-Language", "en_US")
		ctx := context.WithValue(context.Background(), session.KeyRequest, req)
		mail, err := NewMail(ctx, "mock-mail")
		mail.AddTo("p", "a@b.c")
		mail.ReplaceText("%1", "1234")
		So(err, ShouldBeNil)
		err = mail.Send(ctx)
		So(err, ShouldBeNil)
		So(MockResult(), ShouldNotBeNil)
		So(mockResult.GetTo()[0].Name, ShouldEqual, "p")
		So(mockResult.GetTo()[0].Address, ShouldEqual, "a@b.c")
	})
}
