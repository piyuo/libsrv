package sms

import (
	"context"
	"errors"
	"fmt"

	key "github.com/piyuo/libsrv/key"
	"github.com/sfreiberg/gotwilio"
)

// TwilioSMS using twilio to implement SMS
//
type TwilioSMS struct {
	BaseSMS

	// twilio client
	//
	client *gotwilio.Twilio

	// sender mobile number, this number normally assigned by SMS delivery service, please don't use non-verify number
	//
	sender string
}

// newTwilioSMS return SMS instance
//
//	sms := NewSendgridMail("%1 is your verification code.")
//
func newTwilioSMS(template string) (SMS, error) {
	json, err := key.JSON("twilio.json")
	if err != nil {
		return nil, err
	}
	sid := json["sid"].(string)
	if sid == "" {
		return nil, errors.New("sid can not be empty")
	}
	token := json["token"].(string)
	if token == "" {
		return nil, errors.New("token can not be empty")
	}
	sender := json["sender"].(string)
	if sender == "" {
		return nil, errors.New("sender can not be empty,please be aware sender mobile number need verify by SMS delivery service")
	}

	client := gotwilio.NewTwilioClient(sid, token)
	return &TwilioSMS{
		BaseSMS: BaseSMS{
			Text: template,
		},
		client: client,
		sender: sender,
	}, nil
}

//Send using sendgrid to send email
//
//	err := mail.Send(ctx,"+11234567890")
//
func (c *TwilioSMS) Send(ctx context.Context, receiver string) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	_, exception, err := c.client.SendSMS(c.sender, receiver, c.Text, "", "")
	if err != nil {
		return err
	}

	if exception != nil {
		msg := fmt.Sprintf("twillio error, response=%v, message=%v", exception.Status, exception.Message)
		return errors.New(msg)
	}
	return nil
}
