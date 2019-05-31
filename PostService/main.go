package main

import "log"

func main() {
	log.SetFlags(log.LstdFlags)
	config, err := parseConfig("")
	if err != nil {
		log.Panicln("Can't read config:", err)
	}
	r, err := setUpRouter(config.DB.Source)
	if err != nil {
		log.Panicln("Can't set up router:", err)
	}
	r.Run(config.Port)
}
