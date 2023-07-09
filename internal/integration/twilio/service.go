package twilio

import (
	"fmt"
	"github.com/svilenkomitov/notifire/internal/domain"
	"github.com/svilenkomitov/notifire/internal/integration"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
	"regexp"
)

const (
	validatePhoneRegex = "((\\+|\\(|0)?\\d{1,3})?((\\s|\\)|\\-))?(\\d{10})$"
)

type service struct {
	client *twilio.RestClient
}

func New(config *Config) integration.Service {
	return &service{
		client: twilio.NewRestClientWithParams(twilio.ClientParams{
			Username: config.Domain,
			Password: config.ApiKey,
		}),
	}
}

func (s service) Send(notification domain.Notification) error {
	if err := s.Validate(notification); err != nil {
		return err
	}

	_, err := s.client.Api.CreateMessage(&openapi.CreateMessageParams{
		From: &notification.Sender,
		To:   &notification.Recipient,
		Body: &notification.Body,
	})
	return err
}

func (s service) Validate(notification domain.Notification) *integration.ValidationError {
	if notification.Channel.ToLower() != domain.SMS {
		return integration.NewValidationError(fmt.Sprintf("invalid channel [%s]", notification.Channel))
	}

	if isValid, _ := regexp.MatchString(validatePhoneRegex, notification.Sender); !isValid {
		return integration.NewValidationError(fmt.Sprintf("sender [%s] is not valid",
			notification.Sender))
	}

	if isValid, _ := regexp.MatchString(validatePhoneRegex, notification.Recipient); !isValid {
		return integration.NewValidationError(fmt.Sprintf("recipient [%s] is not valid",
			notification.Sender))
	}

	if notification.Body == "" {
		return integration.NewValidationError("body cannot be empty")
	}

	return nil
}
