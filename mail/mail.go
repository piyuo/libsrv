package mail

import (
	"context"
	"errors"
	"strings"
	"time"

	cache "github.com/piyuo/libsrv/cache"
	file "github.com/piyuo/libsrv/file"
)

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
	//	m.Send(ctx)
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

// BaseMail implement basic property of mail
//
type BaseMail struct {
	Mail

	// Subject is mail subject
	//
	Subject string

	// Text is mail text body
	//
	Text string

	// HTML is mail html body
	//
	HTML string

	// FromName is email from name
	//
	FromName string

	// FromAddress is email from address
	//
	FromAddress string

	// To is collection of email address use in email send to
	//
	To []*Email
}

// GetSubject return mail subject
//
//	subject := mail.GetSubject()
//
func (c *BaseMail) GetSubject() string {
	return c.Subject
}

// SetSubject set mail subject
//
//	mail.SetSubject("subject")
//
func (c *BaseMail) SetSubject(subject string) *BaseMail {
	c.Subject = subject
	return c
}

// GetText return mail text content
//
//	text := mail.GetText()
//
func (c *BaseMail) GetText() string {
	return c.Text
}

// SetText set mail text content
//
//	mail.SetText("text body")
//
func (c *BaseMail) SetText(text string) *BaseMail {
	c.Text = text
	return c
}

// GetHTML return mail html content
//
//	html := mail.GetHTML()
//
func (c *BaseMail) GetHTML() string {
	return c.HTML
}

// SetHTML set mail html content
//
//	mail.SetHTML("html body")
//
func (c *BaseMail) SetHTML(html string) *BaseMail {
	c.HTML = html
	return c
}

// GetFrom return from email address
//
//	name,address := mail.GetFrom()
//
func (c *BaseMail) GetFrom() (string, string) {
	return c.FromName, c.FromAddress
}

// SetFrom set from email address
//
//	mail.SetFrom("service","service@piyuo.com")
//
func (c *BaseMail) SetFrom(name, address string) *BaseMail {
	c.FromName = name
	c.FromAddress = address
	return c
}

// ReplaceSubject replace string in mail subject
//
//	mail.ReplaceSubject("%1","hello")
//
func (c *BaseMail) ReplaceSubject(replaceFrom, replaceTo string) *BaseMail {
	c.Subject = strings.ReplaceAll(c.Subject, replaceFrom, replaceTo)
	return c
}

// ReplaceText replace string in mail text content
//
//	mail.ReplaceText("%1","hello")
//
func (c *BaseMail) ReplaceText(replaceFrom, replaceTo string) *BaseMail {
	c.Text = strings.ReplaceAll(c.Text, replaceFrom, replaceTo)
	return c
}

// ReplaceHTML replace string in mail html content
//
//	mail.ReplaceHTML("%1","hello")
//
func (c *BaseMail) ReplaceHTML(replaceFrom, replaceTo string) *BaseMail {
	c.HTML = strings.ReplaceAll(c.HTML, replaceFrom, replaceTo)
	return c
}

// GetTo get email to
//
//	to := mail.GetTo()
//
func (c *BaseMail) GetTo() []*Email {
	return c.To
}

// AddTo add email to
//
//	mail.AddTo("user","user@piyuo.com")
//
func (c *BaseMail) AddTo(emailName, emailAddress string) *BaseMail {
	if c.To == nil {
		c.To = []*Email{}
	}
	email := &Email{
		Name:    emailName,
		Address: emailAddress,
	}
	c.To = append(c.To, email)
	return c
}

// ResetTo reset to empty
//
//	mail.ResetTo()
//
func (c *BaseMail) ResetTo() *BaseMail {
	c.To = nil
	return c
}

// NewMail return Mail instance
//
//	m, err := mail.NewMail("verify", "en-US")
//	m.AddTo("piyuo", "piyuo.com@gmail.com")
//	m.ReplaceText("%1", "1234")
//	m.ReplaceHTML("%1", "1234")
//	m.Send(ctx)
//
func NewMail(templateName, language string) (Mail, error) {
	template, err := getTemplate(templateName, language)
	if err != nil {
		return nil, err
	}
	return newSendgridMail(template)
}

// getTemplate get mail template
//
//	template, err := getTemplate(templateName, language)
//
func getTemplate(templateName, language string) (*template, error) {
	filename := templateName + "_" + language + ".json"
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
	cache.Set(keyname, template, 10*time.Minute) // mail template cache last for 10 min
	return template, nil
}
