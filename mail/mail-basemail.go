package mail

import (
	"strings"
)

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
//	mail.SetFrom("service","service@somedomain.com")
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

// ReplaceContent replace string in mail text and html content
//
//	mail.ReplaceContnet("%1","hello")
//
func (c *BaseMail) ReplaceContent(replaceFrom, replaceTo string) *BaseMail {
	c.Text = strings.ReplaceAll(c.Text, replaceFrom, replaceTo)
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
//	mail.AddTo("user","user@somedomain.com")
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
