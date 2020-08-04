package sms

import (
	"context"
	"strings"

	"github.com/pkg/errors"

	"github.com/nyaruka/phonenumbers"
	cache "github.com/piyuo/libsrv/cache"
	file "github.com/piyuo/libsrv/file"
)

// SMS use template to generate SMS content
//
type SMS interface {

	// GetText return sms text content
	//
	//	text := sms.GetText()
	//
	GetText() string

	// SetText set sms text content
	//
	//	sms.SetText("text body")
	//
	SetText(text string) *BaseSMS

	// ReplaceText replace string in SMS text content
	//
	//	sms.ReplaceText("%1","hello")
	//
	ReplaceText(replaceFrom, replaceTo string) *BaseSMS

	// Send sms
	//
	//	sms, err := NewSMS("verify", "en-US")
	//	sms.ReplaceText("%1", "1234")
	//	sms.ReplaceHTML("%1", "1234")
	//	sms.Send(ctx,"+1123456789")
	//
	Send(ctx context.Context, receiver string) error
}

// BaseSMS implement basic struct of sms
//
type BaseSMS struct {
	SMS

	// Text is sms text body
	//
	Text string
}

// GetText return SMS text content
//
//	text := sms.GetText()
//
func (c *BaseSMS) GetText() string {
	return c.Text
}

// SetText set SMS text content
//
//	sms.SetText("text body")
//
func (c *BaseSMS) SetText(text string) *BaseSMS {
	c.Text = text
	return c
}

// ReplaceText replace string in sms text content
//
//	sms.ReplaceText("%1","hello")
//
func (c *BaseSMS) ReplaceText(replaceFrom, replaceTo string) *BaseSMS {
	c.Text = strings.ReplaceAll(c.Text, replaceFrom, replaceTo)
	return c
}

// NewSMS return sms instance
//
//	sms, err := NewSMS("verify", "en-US")
//	sms.ReplaceText("%1", "1234")
//	sms.ReplaceHTML("%1", "1234")
//	sms.Send(ctx,"+1123456789")
//
func NewSMS(templateName, language string) (SMS, error) {
	templateTxt, err := getTemplate(templateName, language)
	if err != nil {
		return nil, err
	}
	return newTwilioSMS(templateTxt)
}

// getTemplate get sms template
//
//	template, err := getTemplate(templateName, language)
//
func getTemplate(templateName, language string) (string, error) {
	filename := templateName + "_" + language + ".txt"
	keyname := "SMS" + filename
	value, found := cache.Get(keyname)
	if found {
		return value.(string), nil
	}

	filepath, found := file.Find("assets/sms/" + filename)
	if !found {
		return "", errors.New("SMS template " + filename + " not found")
	}

	txt, err := file.ReadText(filepath)
	if err != nil {
		return "", err
	}
	cache.Set(cache.MEDIUM, keyname, txt)
	return txt, nil
}

// E164 return E164 format international number
//
//	mobile, err := E164(9493026176, "US")
//
func E164(phoneNumber, countryCode string) (string, error) {
	num, err := phonenumbers.Parse(phoneNumber, countryCode)
	if err != nil {
		return "", errors.Wrap(err, "failed to parse phone number: "+phoneNumber+", country:"+countryCode)
	}

	if !phonenumbers.IsValidNumber(num) {
		return "", errors.New("number: " + phoneNumber + ", country:" + countryCode + " is not a valid number")
	}

	formatted := phonenumbers.Format(num, phonenumbers.E164)

	return formatted, nil
}
