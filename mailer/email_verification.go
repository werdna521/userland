package mailer

import (
	"context"
	"fmt"
)

func SendEmailVerificationMail(
	ctx context.Context,
	m Mailer,
	to Email,
	link string,
) error {
	mo := &MailOptions{
		To:          []Email{to},
		Subject:     "Verify your email",
		HTMLContent: fmt.Sprintf(emailVerificationTemplate, link),
		TextContent: "Hi Userlanders, please verify your email",
	}

	return m.SendMail(ctx, mo)
}
