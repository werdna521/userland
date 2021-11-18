package mailer

import (
	"context"
	"fmt"
)

func SendPasswordResetMail(
	ctx context.Context,
	m Mailer,
	to Email,
	token string,
) error {
	mo := &MailOptions{
		To:          []Email{to},
		Subject:     "Reset Password",
		HTMLContent: fmt.Sprintf(passwordResetTemplate, token),
		TextContent: "Hi Userlanders, use this token to reset your password",
	}

	return m.SendMail(ctx, mo)
}
