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
		notification.Message.Sender,
		notification.Message.Subject,
		notification.Message.Body,
		notification.Message.Recipient))
	return err
}

func (s service) Validate(notification domain.Notification) *integration.ValidationError {
	if notification.Channel.ToLower() != domain.EMAIL {
		return integration.NewValidationError(fmt.Sprintf("invalid channel [%s]", notification.Channel))
	}

	_, err := mail.ParseAddress(notification.Message.Sender)
	if err != nil {
		return integration.NewValidationError(fmt.Sprintf("sender [%s] is not valid",
			notification.Message.Sender))
	}

	_, err = mail.ParseAddress(notification.Message.Recipient)
	if err != nil {
		return integration.NewValidationError(fmt.Sprintf("recipient [%s] is not valid",
			notification.Message.Sender))
	}

	if notification.Message.Subject == "" {
		return integration.NewValidationError("subject cannot be empty")
	}

	if notification.Message.Body == "" {
		return integration.NewValidationError("body cannot be empty")
	}
	return nil
}
