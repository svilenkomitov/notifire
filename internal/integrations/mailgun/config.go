package mailgun

import (
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
)

const (
	envPrefix string = "MAILGUN"
)

type Config struct {
	Domain string `envconfig:"DOMAIN"`
	ApiKey string `envconfig:"API_KEY"`
}

func LoadConfig() *Config {
	var config Config
	if err := envconfig.Process(envPrefix, &config); err != nil {
		log.Fatal(err)
	}
	return &config
}
