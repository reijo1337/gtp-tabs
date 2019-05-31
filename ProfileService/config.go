package main

import (
	"github.com/kelseyhightower/envconfig"
)

type config struct {
	DB struct {
		User     string `default:"test" envconfig:"DB_USER"`
		Password string `default:"test" envconfig:"DB_PASSWORD"`
		Host     string `default:"127.0.0.1" envconfig:"DB_HOST"`
		Name     string `default:"profiles" envconfig:"DB_NAME"`
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
