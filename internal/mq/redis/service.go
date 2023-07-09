package redis

import (
	"fmt"
	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
	"github.com/svilenkomitov/notifire/internal/domain"
	"github.com/svilenkomitov/notifire/internal/integration"
	"github.com/svilenkomitov/notifire/internal/integration/mailgun"
	"github.com/svilenkomitov/notifire/internal/integration/slack"
	"github.com/svilenkomitov/notifire/internal/integration/twilio"
	"github.com/svilenkomitov/notifire/internal/mq"
	"github.com/svilenkomitov/notifire/internal/notification"
)

type service struct {
	client         *redis.Client
	mailgunService integration.Service
	twilioService  integration.Service
	slackService   integration.Service
	repository     notification.Repository
}

func New(config *Config, repository notification.Repository) mq.Service {
	redisClient := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", config.Host, config.Port),
	})
	return service{
		client:         redisClient,
		mailgunService: mailgun.New(mailgun.LoadConfig()),
		twilioService:  twilio.New(twilio.LoadConfig()),
		slackService:   slack.New(slack.LoadConfig()),
		repository:     repository,
	}.Subscribe()
}

func (s service) Publish(notification domain.Notification) error {
	notification.Status = domain.Pending
	id, err := s.repository.Create(notification)
	if err != nil {
		return err
	}
	notification.ID = id
	return s.client.Publish(string(notification.Channel.ToLower()), notification).Err()
}

func (s service) Subscribe() mq.Service {
	mailChannel := s.client.Subscribe(string(domain.EMAIL.ToLower())).Channel()
	go s.consumeMessages(mailChannel, s.mailgunService.Send)

	smsChannel := s.client.Subscribe(string(domain.SMS.ToLower())).Channel()
	go s.consumeMessages(smsChannel, s.twilioService.Send)

	slackChannel := s.client.Subscribe(string(domain.SLACK.ToLower())).Channel()
	go s.consumeMessages(slackChannel, s.slackService.Send)

	return s
}

type handler func(notification domain.Notification) error

func (s service) consumeMessages(ch <-chan *redis.Message, handler handler) {
	for msg := range ch {
		notification, err := domain.Unmarshal([]byte(msg.Payload))
		if err != nil {
			log.Error(err)
			if err := s.repository.UpdateStatus(notification.ID, domain.Failed); err != nil {
				log.Error(err)
			}
			continue
		}

		if err = handler(notification); err != nil {
			log.Error(err)
			if err := s.repository.UpdateStatus(notification.ID, domain.Failed); err != nil {
				log.Error(err)
			}
			continue
		}

		if err := s.repository.UpdateStatus(notification.ID, domain.Sent); err != nil {
			log.Error(err)
		}
	}
}
