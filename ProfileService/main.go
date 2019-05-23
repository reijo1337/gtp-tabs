package main

import "log"

func main() {
	log.SetFlags(log.LstdFlags)
	config, err := parseConfig("GATEWAY")
	if err != nil {
		log.Panicln("Can't read config:", err)
	}
	r, err := setUpRouter()
	if err != nil {
		log.Panicln("Can't set up router:", err)
	}
	r.Run(":" + config.Port)
}
