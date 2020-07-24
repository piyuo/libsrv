package mail

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMail(t *testing.T) {
	Convey("should set/get method", t, func() {
		mail, err := NewMail("verify", "en-US")
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
		mail, err = NewMail("verify", "en-US")
		So(err, ShouldBeNil)
		So(mail, ShouldNotBeNil)
		So(mail.GetSubject(), ShouldEqual, backupSubject)

		mail, err = NewMail("not exist", "en-US")
		So(err, ShouldNotBeNil)
		So(mail, ShouldBeNil)
	})
}

func TestSendMail(t *testing.T) {
	Convey("should send mail", t, func() {
		ctx := context.Background()
		mail, err := NewMail("verify", "en-US")
		So(err, ShouldBeNil)
		mail.AddTo("piyuo", "piyuo.com@gmail.com")
		mail.ReplaceText("%1", "1234")
		mail.ReplaceHTML("%1", "1234")
		err = mail.Send(ctx)
		So(err, ShouldBeNil)
	})
}
