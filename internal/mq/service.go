package mq

import "github.com/svilenkomitov/notifire/internal/domain"

type Service interface {
	Subscribe() Service
	Publish(notification domain.Notification) error
}
