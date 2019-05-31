package main

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type config struct {
	DB struct {
		User     string `default:"postgres" envconfig:"DB_USER"`
		Password string `default:"postgres" envconfig:"DB_PASSWORD"`
		Host     string `default:"127.0.0.1" envconfig:"DB_HOST"`
		Name     string `default:"auth" envconfig:"DB_NAME"`
	}
	Token struct {
		PrivateKey        string        `default:"keys/private.pem" envconfig:"PRIVATE_KEY"`
		AccessExpiration  time.Duration `default:"30m" envconfig:"TOKEN_ACCESS_EXPIRATION"`
		RefreshExpiration time.Duration `default:"300m" envconfig:"TOKEN_REFRESH_EXPIRATION"`
	}
	Port string `default:"9090"`
}

func parseConfig(app string) (cfg config, err error) {
	if err := envconfig.Process(app, &cfg); err != nil {
		envconfig.Usage(app, &cfg)
		return cfg, err
	}
	return cfg, nil
}
