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

func (s service) Send(n domain.Notification) (domain.Notification, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	n.Status = domain.Pending

	if err := s.Validate(n); err != nil {
		n.Status = domain.Failed
		return n, err
	}

	message := s.client.NewMessage(
		n.Message.Sender,
		n.Message.Subject,
		n.Message.Body,
		n.Message.Recipient)

	_, _, err := s.client.Send(ctx, message)
	if err != nil {
		n.Status = domain.Failed
	}

	n.Status = domain.Sent
	return n, err
}

func (s service) Validate(n domain.Notification) *integration.ValidationError {
	if n.Channel.ToLower() != domain.EMAIL {
		return integration.NewValidationError(fmt.Sprintf("invalid channel [%s]", n.Channel))
	}

	_, err := mail.ParseAddress(n.Message.Sender)
	if err != nil {
		return integration.NewValidationError(fmt.Sprintf("sender [%s] is not valid",
			n.Message.Sender))
	}

	_, err = mail.ParseAddress(n.Message.Recipient)
	if err != nil {
		return integration.NewValidationError(fmt.Sprintf("recipient [%s] is not valid",
			n.Message.Sender))
	}

	if n.Message.Subject == "" {
		return integration.NewValidationError("subject cannot be empty")
	}

	if n.Message.Body == "" {
		return integration.NewValidationError("body cannot be empty")
	}
	return nil
}
