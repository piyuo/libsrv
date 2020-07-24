package sms

import (
	"testing"

	"github.com/piyuo/libsrv/util"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSMS(t *testing.T) {
	Convey("should set/get method", t, func() {
		sms, err := NewSMS("verify", "en-US")
		So(err, ShouldBeNil)
		backupText := sms.GetText()
		So(sms.GetText(), ShouldNotBeEmpty)
		sms.SetText("ok")
		So(sms.GetText(), ShouldEqual, "ok")
		sms.ReplaceText("ok", "1")
		So(sms.GetText(), ShouldEqual, "1")

		//should from cache
		sms, err = NewSMS("verify", "en-US")
		So(err, ShouldBeNil)
		So(sms, ShouldNotBeNil)
		So(sms.GetText(), ShouldEqual, backupText)

	})
}

func TestSMSError(t *testing.T) {
	Convey("should set/get method", t, func() {
		//test not exist template
		sms, err := NewSMS("not exist", "en-US")
		So(err, ShouldNotBeNil)
		So(sms, ShouldBeNil)

		sms, err = NewSMS("verify", "en-US")
		So(err, ShouldBeNil)

		//test canceled ctx
		canceledCtx := util.CanceledCtx()
		err = sms.Send(canceledCtx, "+19999999999")
		So(err, ShouldNotBeNil)
	})
}

func TestSendSMS(t *testing.T) {
	Convey("should send SMS", t, func() {
		sms, err := NewSMS("verify", "en-US")
		So(err, ShouldBeNil)
		sms.ReplaceText("%1", "1234")
		//err = sms.Send(context.Background(), "+19493026176")
		So(err, ShouldBeNil)
	})
}
