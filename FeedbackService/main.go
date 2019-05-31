package main

import "log"

func main() {
	log.SetFlags(log.LstdFlags)
	config, err := parseConfig("")
	if err != nil {
		log.Panicln("Can't read config:", err)
	}
	r, err := setUpRouter(config.DB.Source, config.SMTP.URL, config.SMTP.Login, config.SMTP.Password, config.SMTP.Port)
	if err != nil {
		log.Panicln("Can't set up server:", err)
	}
	log.Println("Starting server on port ", config.Port)
	r.Run(config.Port)
}
