package main

import (
	"github.com/kelseyhightower/envconfig"
)

type config struct {
	DB struct {
		User     string `default:"postgres" envconfig:"DB_USER"`
		Password string `default:"postgres" envconfig:"DB_PASSWORD"`
		Host     string `default:"127.0.0.1" envconfig:"DB_HOST"`
		Name     string `default:"auth" envconfig:"DB_NAME"`
	}
	PublicColFile string
	Port          string `default:"8081"`
}

func parseConfig(app string) (cfg config, err error) {
	if err := envconfig.Process(app, &cfg); err != nil {
		envconfig.Usage(app, &cfg)
		return cfg, err
	}
	return cfg, nil
}
