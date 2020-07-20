package mail

import (
	"context"
)

// Mail use template to generate mail content and send
//
type Mail interface {
	Send(ctx context.Context, to string) error
}

// NewMail return Mail instance
//
//	inTx := conn.IsInTransaction()
//
func NewMail(template, language string) Mail {

	return &Sendgrid{
		template: template,
		language: language,
	}
}

func getJSON(template, language string) error {
	path:=
}

//Sendgrid using SendGrid to implement mail
//
type Sendgrid struct {
	Mail

	// template indicate which template to be use
	//
	template string

	// language indicate which language to be use
	//
	language string
}

//Send email
//
//	mail.Send('from','to','subject','text')
func Send(ctx context.Context, to string) {
	//	SendGrid(ctx, from, to, subject, text)
}
