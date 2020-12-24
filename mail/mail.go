package mail

import (
	"context"
	"errors"

	cache "github.com/piyuo/libsrv/cache"
	file "github.com/piyuo/libsrv/file"
)

var mockMailService = false

// MockMailService fake send mail by always send success
//
func MockMailService(mock bool) {
	mockMailService = mock
}

type template struct {
	subject string

	text string

	html string

	fromName string

	fromAddress string
}

// Mail use template to generate mail content and send
//
type Mail interface {

	// GetSubject return mail subject
	//
	//	subject := mail.GetSubject()
	//
	GetSubject() string

	// SetSubject set mail subject
	//
	//	mail.SetSubject("subject")
	//
	SetSubject(subject string) *BaseMail

	// GetText return mail text content
	//
	//	text := mail.GetText()
	//
	GetText() string

	// SetText set mail text content
	//
	//	mail.SetText("text body")
	//
	SetText(text string) *BaseMail

	// GetHTML return mail html content
	//
	//	html := mail.GetHTML()
	//
	GetHTML() string

	// SetHTML set mail html content
	//
	//	mail.SetHTML("html body")
	//
	SetHTML(html string) *BaseMail

	// GetFrom return from email address
	//
	//	name,address := mail.GetFrom()
	//
	GetFrom() (string, string)

	// SetFrom set from email address
	//
	//	mail.SetFrom("service","service@piyuo.com")
	//
	SetFrom(emailName, emailAddress string) *BaseMail

	// ReplaceSubject replace string in mail subject
	//
	//	mail.ReplaceSubject("%1","hello")
	//
	ReplaceSubject(replaceFrom, replaceTo string) *BaseMail

	// ReplaceText replace string in mail text content
	//
	//	mail.ReplaceText("%1","hello")
	//
	ReplaceText(replaceFrom, replaceTo string) *BaseMail

	// ReplaceHTML replace string in html comtent
	//
	//	to := mail.GetTo()
	//
	ReplaceHTML(replaceFrom, replaceTo string) *BaseMail

	// GetTo get email to
	//
	//	to := mail.GetTo()
	//
	GetTo() []*Email

	// AddTo add email to
	//
	//	mail.AddTo("user","user@piyuo.com")
	//
	AddTo(emailName, emailAddress string) *BaseMail

	// ResetTo reset to empty
	//
	//	mail.ResetTo()
	//
	ResetTo() *BaseMail

	// Send mail
	//
	//	m, err := mail.NewMail("verify", "en-US")
	//	m.AddTo("piyuo", "piyuo.com@gmail.com")
	//	m.ReplaceText("%1", "1234")
	//	m.ReplaceHTML("%1", "1234")
	//	err := m.Send(ctx)
	//
	Send(ctx context.Context) error
}

// Email is a single email address
//
type Email struct {

	// Name is email name
	//
	Name string

	// Address is email address
	//
	Address string
}

// NewMail return Mail instance, require template name and locale to find template
//
//	m, err := mail.NewMail("verify", "en-us")
//	m.AddTo("piyuo", "piyuo.com@gmail.com")
//	m.ReplaceText("%1", "1234")
//	m.ReplaceHTML("%1", "1234")
//	m.Send(ctx)
//
func NewMail(templateName, locale string) (Mail, error) {
	template, err := getTemplate(templateName, locale)
	if err != nil {
		return nil, err
	}
	return newSendgridMail(template)
}

// getTemplate get mail template
//
//	template, err := getTemplate("verify", "en-us")
//
func getTemplate(templateName, locale string) (*template, error) {
	filename := templateName + "_" + locale + ".json"
	keyname := "MAIL" + filename
	value, found := cache.Get(keyname)
	if found {
		return value.(*template), nil
	}

	filepath, found := file.Find("assets/mail/" + filename)
	if !found {
		return nil, errors.New("mail template " + filename + " not found")
	}

	json, err := file.ReadJSON(filepath)
	if err != nil {
		return nil, err
	}
	template := &template{
		subject:     json["subject"].(string),
		text:        json["text"].(string),
		html:        json["html"].(string),
		fromName:    json["fromName"].(string),
		fromAddress: json["fromAddress"].(string),
	}
	cache.Set(cache.MEDIUM, keyname, template)
	return template, nil
}
