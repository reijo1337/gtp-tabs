package main

import (
	"github.com/kelseyhightower/envconfig"
)

type config struct {
	DB struct {
		Source string `default:"user=test password=test dbname=posts host=127.0.0.1 sslmode=disable"`
	}
	Port string `default:"8084" envconfig:"PORT"`
}

func parseConfig(app string) (cfg config, err error) {
	if err := envconfig.Process(app, &cfg); err != nil {
		_ = envconfig.Usage(app, &cfg)
		return cfg, err
	}
	return cfg, nil
}
