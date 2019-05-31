package main

import (
	"github.com/kelseyhightower/envconfig"
)

type config struct {
	DB struct {
		Source string `default:"user=test password=test dbname=feedback host=127.0.0.1 sslmode=disable"`
	}
	SMTP struct {
		URL      string `default:"smtp.yandex.ru"`
		Port     int    `default:"465"`
		Login    string `default:"tantsevov"`
		Password string
	}
	Port string `default:":9090" envconfig:"PORT"`
}

func parseConfig(app string) (cfg config, err error) {
	if err := envconfig.Process(app, &cfg); err != nil {
		_ = envconfig.Usage(app, &cfg)
		return cfg, err
	}
	return cfg, nil
}
