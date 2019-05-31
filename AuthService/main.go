package main

import (
	"log"
)

func main() {
	log.SetFlags(log.LstdFlags)
	config, err := parseConfig("")
	if err != nil {
		log.Panicln("Can't read config:", err)
	}
	r, err := SetUpRouter(config.Token.PrivateKey, config.Token.AccessExpiration, config.Token.RefreshExpiration)
	if err != nil {
		log.Panicln("Can't set up server:", err)
	}
	log.Println("Starting server on port ", config.Port)
	r.Run(config.Port)
}
