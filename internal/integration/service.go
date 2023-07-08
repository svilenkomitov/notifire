package integration

import (
	"github.com/svilenkomitov/notifire/internal/domain"
)

type Service interface {
	Send(notification domain.Notification) (domain.Notification, error)
	Validate(notification domain.Notification) *ValidationError
}
