package config

import (
	"github.com/apex/log"
	"github.com/caarlos0/env"
)

// Config of the app
type Config struct {
	JbossHome string   `env:"JBOSS_HOME"`
	Args      string   `env:"GOBOSS_ARGS"`
	BuildArgs []string `env:"GOBOSS_BUILD_ARGS" envSeparator:","`
	LogLevel  string   `env:"LOG_LEVEL" envDefault:"debug"`
}

// MustGet returns the config
func MustGet() Config {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		log.WithError(err).Fatal("failed to load config")
	}
	return cfg
}
