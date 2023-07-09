package mailgun

import (
	"context"
	"fmt"
	"github.com/mailgun/mailgun-go/v4"
	"github.com/svilenkomitov/notifire/internal/domain"
	"github.com/svilenkomitov/notifire/internal/integration"
	"net/mail"
	"time"
)

type service struct {
	client *mailgun.MailgunImpl
}

func New(config *Config) integration.Service {
	return &service{
		client: mailgun.NewMailgun(config.Domain, config.ApiKey),
	}
}

func (s service) Send(notification domain.Notification) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := s.Validate(notification); err != nil {
		return err
	}

	_, _, err := s.client.Send(ctx, s.client.NewMessage(
		notification.Sender,
		notification.Subject,
		notification.Body,
		notification.Recipient))
	return err
}

func (s service) Validate(notification domain.Notification) *integration.ValidationError {
	if notification.Channel.ToLower() != domain.EMAIL {
		return integration.NewValidationError(fmt.Sprintf("invalid channel [%s]", notification.Channel))
	}

	_, err := mail.ParseAddress(notification.Sender)
	if err != nil {
		return integration.NewValidationError(fmt.Sprintf("sender [%s] is not valid",
			notification.Sender))
	}

	_, err = mail.ParseAddress(notification.Recipient)
	if err != nil {
		return integration.NewValidationError(fmt.Sprintf("recipient [%s] is not valid",
			notification.Sender))
	}

	if notification.Subject == "" {
		return integration.NewValidationError("subject cannot be empty")
	}

	if notification.Body == "" {
		return integration.NewValidationError("body cannot be empty")
	}
	return nil
}
