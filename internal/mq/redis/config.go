package redis

import (
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
)

const (
	envPrefix string = "REDIS"
)

type Config struct {
	Host string `envconfig:"HOST"`
	Port uint16 `default:"6379"  envconfig:"PORT"`
}

func LoadConfig() *Config {
	var config Config
	if err := envconfig.Process(envPrefix, &config); err != nil {
		log.Fatal(err)
	}
	return &config
}
