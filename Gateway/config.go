package main

import (
	"github.com/kelseyhightower/envconfig"
)

type config struct {
	Port    string `default:"8080" envconfig:"PORT"`
	Storage struct {
		Host string `default:"127.0.0.1" envconfig:"STORAGE_HOST"`
		Port string `default:"8081" envconfig:"STORAGE_PORT"`
	}
	PublicKeyLoc string `envconfig:"PUBLIC_KEY_LOC"`
}

func parseConfig(app string) (cfg config, err error) {
	if err := envconfig.Process(app, &cfg); err != nil {
		envconfig.Usage(app, &cfg)
		return cfg, err
	}
	return cfg, nil
}
