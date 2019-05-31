package main

import (
	"io/ioutil"
	"log"
)

func main() {
	log.SetFlags(log.LstdFlags)
	config, err := parseConfig("")
	if err != nil {
		log.Panicln("Can't read config:", err)
	}
	authPublicKey, err := ioutil.ReadFile(config.PublicKeyLoc)
	if err != nil {
		log.Panicln("reading auth token verification public key:", err)
	}
	r, err := setUpRouter(authPublicKey)
	if err != nil {
		log.Panicln("Can't set up router:", err)
	}
	r.Run(config.Port)
}
