package redis

import (
	"encoding/json"
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

	nBytes, err := notification.MarshalBinary()
	if err != nil {
		return err
	}

	var data map[string]interface{}
	err = json.Unmarshal(nBytes, &data)
	if err != nil {
		return err
	}

	_, err = s.client.XAdd(&redis.XAddArgs{
		Stream: string(notification.Channel.ToLower()),
		Values: data,
	}).Result()
	return err
}

func (s service) Subscribe() mq.Service {
	go s.consumeMessages(string(domain.EMAIL.ToLower()), s.mailgunService.Send)
	go s.consumeMessages(string(domain.SMS.ToLower()), s.twilioService.Send)
	go s.consumeMessages(string(domain.SLACK.ToLower()), s.slackService.Send)
	return s
}

type handler func(notification domain.Notification) error

func (s service) consumeMessages(channel string, handler handler) {
	for {
		streams, err := s.client.XRead(&redis.XReadArgs{
			Streams: []string{channel, "$"},
			Count:   1,
			Block:   0,
		}).Result()
		if err != nil {
			log.Error(err)
		}

		for _, messages := range streams {
			for _, message := range messages.Messages {
				jsonData, err := json.Marshal(message.Values)
				if err != nil {
					log.Error(err)
					continue
				}

				var notification domain.Notification
				err = json.Unmarshal(jsonData, &notification)
				if err != nil {
					log.Error(err)
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
	}
}
