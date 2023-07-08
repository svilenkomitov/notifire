package slack

import (
	"github.com/slack-go/slack"
	"github.com/svilenkomitov/notifire/internal/domain"
	"github.com/svilenkomitov/notifire/internal/integration"
)

type service struct {
	client *slack.Client
}

func New(config *Config) integration.Service {
	return &service{
		client: slack.New(config.ApiKey),
	}
}

func (s service) Send(notification domain.Notification) error {
	_, _, err := s.client.PostMessage(notification.Message.Recipient,
		slack.MsgOptionAttachments(slack.Attachment{
			Pretext: notification.Message.Subject,
			Text:    notification.Message.Body,
		}))
	return err
}

func (s service) Validate(notification domain.Notification) *integration.ValidationError {
	if notification.Message.Recipient == "" {
		return integration.NewValidationError("recipient cannot be empty")
	}

	if notification.Message.Subject == "" {
		return integration.NewValidationError("subject cannot be empty")
	}

	if notification.Message.Body == "" {
		return integration.NewValidationError("body cannot be empty")
	}

	return nil
}
