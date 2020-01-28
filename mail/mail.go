package mail

import (
	"context"
)

//Send email
//
//	mail.Send('from','to','subject','text')
func Send(ctx context.Context, from, to, subject, text string) {
	SendGrid(ctx, from, to, subject, text)
}
