package server

import (
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
)

const (
	envPrefix string = "SERVER"
)

type Config struct {
	Port uint16 `default:"8080"  envconfig:"PORT"`
}

func LoadConfig() *Config {
	var config Config
	if err := envconfig.Process(envPrefix, &config); err != nil {
		log.Fatal(err)
	}
	return &config
}
