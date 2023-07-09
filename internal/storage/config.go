package storage

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
)

const (
	envPrefix          string = "DB"
	dataSourceTemplate string = "host=%s port=%d user=%s password=%s dbname=%s"
)

type Config struct {
	Host     string `envconfig:"HOST"`
	Port     int    `envconfig:"PORT"`
	User     string `envconfig:"USER"`
	Password string `envconfig:"PASS"`
	DbName   string `envconfig:"NAME"`
	DbDriver string `default:"pgx" envconfig:"DRIVER"`
}

func LoadConfig() *Config {
	var config Config
	if err := envconfig.Process(envPrefix, &config); err != nil {
		log.Fatal(err)
	}
	return &config
}

func (config Config) GetDataSourceName() string {
	return fmt.Sprintf(dataSourceTemplate, config.Host, config.Port, config.User,
		config.Password, config.DbName)
}
