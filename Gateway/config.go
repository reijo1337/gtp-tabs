package main

import (
	"github.com/kelseyhightower/envconfig"
)

type config struct {
	Port string `default:":9090" envconfig:"PORT"`
	URL  struct {
		Storage string
		Auth    string
		Post    string
		Profile string
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
