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
)

type service struct {
	client         *redis.Client
	mailgunService integration.Service
	twilioService  integration.Service
	slackService   integration.Service
}

func New(config *Config) mq.Service {
	redisClient := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", config.Host, config.Port),
	})
	return service{
		client:         redisClient,
		mailgunService: mailgun.New(mailgun.LoadConfig()),
		twilioService:  twilio.New(twilio.LoadConfig()),
		slackService:   slack.New(slack.LoadConfig()),
	}.Subscribe()
}

func (s service) Publish(notification domain.Notification) error {
	return s.client.Publish(string(notification.Channel.ToLower()), notification).Err()
}

func (s service) Subscribe() mq.Service {
	mailChannel := s.client.Subscribe(string(domain.EMAIL.ToLower())).Channel()
	go consumeMessages(mailChannel, s.mailgunService.Send)

	smsChannel := s.client.Subscribe(string(domain.SMS.ToLower())).Channel()
	go consumeMessages(smsChannel, s.twilioService.Send)

	slackChannel := s.client.Subscribe(string(domain.SLACK.ToLower())).Channel()
	go consumeMessages(slackChannel, s.slackService.Send)

	return s
}

type handler func(notification domain.Notification) error

func consumeMessages(ch <-chan *redis.Message, handler handler) {
	for msg := range ch {
		notification, err := domain.Unmarshal([]byte(msg.Payload))
		if err != nil {
			log.Error(err)
			return
		}

		if err = handler(notification); err != nil {
			log.Error(err)
		}
	}
}
