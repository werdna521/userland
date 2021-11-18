package mailer

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
)

const sendinblueAPIURL = "https://api.sendinblue.com/v3/smtp/email"

type Mailer interface {
	SendMail(ctx context.Context, mo *MailOptions) error
}

type BaseMailer struct {
	Sender Email
	APIKey string
}

func NewBaseMailer(config Config) *BaseMailer {
	return &BaseMailer{
		Sender: Email{
			Name:  config.SenderName,
			Email: config.SenderEmail,
		},
		APIKey: config.APIKey,
	}
}

type Config struct {
	SenderName  string
	SenderEmail string
	APIKey      string
}

type MailerBody struct {
	Sender      Email   `json:"sender"`
	To          []Email `json:"to"`
	HTMLContent string  `json:"htmlContent"`
	TextContent string  `json:"textContent"`
	Subject     string  `json:"subject"`
}

type Email struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type MailOptions struct {
	To          []Email
	HTMLContent string
	TextContent string
	Subject     string
}

func (m *BaseMailer) SendMail(ctx context.Context, mo *MailOptions) error {
	body := MailerBody{
		Sender: Email{
			Name:  m.Sender.Name,
			Email: m.Sender.Email,
		},
		To:          mo.To,
		HTMLContent: mo.HTMLContent,
		TextContent: mo.TextContent,
		Subject:     mo.Subject,
	}
	log.Info().Msg("stringify-ing request body")
	bodyStr, err := json.Marshal(body)
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal body")
		return err
	}

	// TODO: create a simple custom http client
	log.Info().Msg("creating http request to send email")
	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		sendinblueAPIURL,
		bytes.NewBuffer(bodyStr),
	)
	if err != nil {
		log.Error().Err(err).Msg("failed to create request")
		return err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-key", m.APIKey)

	_, err = http.DefaultClient.Do(req)
	return err
}
