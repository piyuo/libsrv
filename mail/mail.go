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
	SetSubject(subject string) *SMTPMail

	// GetText return mail text content
	//
	//	text := mail.GetText()
	//
	GetText() string

	// SetText set mail text content
	//
	//	mail.SetText("text body")
	//
	SetText(text string) *SMTPMail

	// GetHTML return mail html content
	//
	//	html := mail.GetHTML()
	//
	GetHTML() string

	// SetHTML set mail html content
	//
	//	mail.SetHTML("html body")
	//
	SetHTML(html string) *SMTPMail

	// GetFrom return from email address
	//
	//	name,address := mail.GetFrom()
	//
	GetFrom() (string, string)

	// SetFrom set from email address
	//
	//	mail.SetFrom("service","service@piyuo.com")
	//
	SetFrom(emailName, emailAddress string) *SMTPMail

	// ReplaceSubject replace string in mail subject
	//
	//	mail.ReplaceSubject("%1","hello")
	//
	ReplaceSubject(replaceFrom, replaceTo string) *SMTPMail

	// ReplaceText replace string in mail text content
	//
	//	mail.ReplaceText("%1","hello")
	//
	ReplaceText(replaceFrom, replaceTo string) *SMTPMail

	// GetTo get email to
	//
	//	to := mail.GetTo()
	//
	ReplaceHTML(replaceFrom, replaceTo string) *SMTPMail

	// GetTo get email to
	//
	//	to := mail.GetTo()
	//
	GetTo() []*Email

	// AddTo add email to
	//
	//	mail.AddTo("user","user@piyuo.com")
	//
	AddTo(emailName, emailAddress string) *SMTPMail

	// ResetTo reset to empty
	//
	//	mail.ResetTo()
	//
	ResetTo() *SMTPMail

	// Send mail
	//
	//	mail, err := NewMail("verify", "en-US")
	//	mail.AddTo("piyuo", "piyuo.com@gmail.com")
	//	mail.ReplaceText("%1", "1234")
	//	mail.ReplaceHTML("%1", "1234")
	//	mail.Send(ctx)
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

// SMTPMail implement basic property of mail
//
type SMTPMail struct {
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
func (c *SMTPMail) GetSubject() string {
	return c.Subject
}

// SetSubject set mail subject
//
//	mail.SetSubject("subject")
//
func (c *SMTPMail) SetSubject(subject string) *SMTPMail {
	c.Subject = subject
	return c
}

// GetText return mail text content
//
//	text := mail.GetText()
//
func (c *SMTPMail) GetText() string {
	return c.Text
}

// SetText set mail text content
//
//	mail.SetText("text body")
//
func (c *SMTPMail) SetText(text string) *SMTPMail {
	c.Text = text
	return c
}

// GetHTML return mail html content
//
//	html := mail.GetHTML()
//
func (c *SMTPMail) GetHTML() string {
	return c.HTML
}

// SetHTML set mail html content
//
//	mail.SetHTML("html body")
//
func (c *SMTPMail) SetHTML(html string) *SMTPMail {
	c.HTML = html
	return c
}

// GetFrom return from email address
//
//	name,address := mail.GetFrom()
//
func (c *SMTPMail) GetFrom() (string, string) {
	return c.FromName, c.FromAddress
}

// SetFrom set from email address
//
//	mail.SetFrom("service","service@piyuo.com")
//
func (c *SMTPMail) SetFrom(name, address string) *SMTPMail {
	c.FromName = name
	c.FromAddress = address
	return c
}

// ReplaceSubject replace string in mail subject
//
//	mail.ReplaceSubject("%1","hello")
//
func (c *SMTPMail) ReplaceSubject(replaceFrom, replaceTo string) *SMTPMail {
	c.Subject = strings.ReplaceAll(c.Subject, replaceFrom, replaceTo)
	return c
}

// ReplaceText replace string in mail text content
//
//	mail.ReplaceText("%1","hello")
//
func (c *SMTPMail) ReplaceText(replaceFrom, replaceTo string) *SMTPMail {
	c.Text = strings.ReplaceAll(c.Text, replaceFrom, replaceTo)
	return c
}

// ReplaceHTML replace string in mail html content
//
//	mail.ReplaceHTML("%1","hello")
//
func (c *SMTPMail) ReplaceHTML(replaceFrom, replaceTo string) *SMTPMail {
	c.HTML = strings.ReplaceAll(c.HTML, replaceFrom, replaceTo)
	return c
}

// GetTo get email to
//
//	to := mail.GetTo()
//
func (c *SMTPMail) GetTo() []*Email {
	return c.To
}

// AddTo add email to
//
//	mail.AddTo("user","user@piyuo.com")
//
func (c *SMTPMail) AddTo(emailName, emailAddress string) *SMTPMail {
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
func (c *SMTPMail) ResetTo() *SMTPMail {
	c.To = nil
	return c
}

// NewMail return Mail instance
//
//	mail, err := NewMail("verify", "en-US")
//	mail.AddTo("piyuo", "piyuo.com@gmail.com")
//	mail.ReplaceText("%1", "1234")
//	mail.ReplaceHTML("%1", "1234")
//	mail.Send(ctx)
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
	cache.Set(keyname, template, 10*time.Minute) // key never expire, cause we always need it
	return template, nil
}
