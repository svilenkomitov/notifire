package redis

import (
	"fmt"
	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
	"github.com/svilenkomitov/notifire/internal/domain"
	"github.com/svilenkomitov/notifire/internal/integration"
	"github.com/svilenkomitov/notifire/internal/integration/mailgun"
	"github.com/svilenkomitov/notifire/internal/mq"
)

type service struct {
	client         *redis.Client
	mailgunService integration.Service
}

func New(config *Config) mq.Service {
	redisClient := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", config.Host, config.Port),
	})
	return service{
		client:         redisClient,
		mailgunService: mailgun.New(mailgun.LoadConfig()),
	}.Subscribe()
}

func (s service) Publish(notification domain.Notification) error {
	return s.client.Publish(string(notification.Channel.ToLower()), notification).Err()
}

func (s service) Subscribe() mq.Service {
	mailChannel := s.client.Subscribe(string(domain.EMAIL.ToLower())).Channel()
	go consumeMessages(mailChannel, s.mailgunService.Send)
	return s
}

type handler func(notification domain.Notification) (domain.Notification, error)

func consumeMessages(ch <-chan *redis.Message, handler handler) {
	for msg := range ch {
		notification, err := domain.Unmarshal([]byte(msg.Payload))
		if err != nil {
			log.Error(err)
			return
		}
		_, err = handler(notification)
		if err != nil {
			log.Error(err)
		}
	}
}
