package main

import (
	"log"
	"os"
)

// Инициализирует глобальные переменные на основе системных
func init() {
	DatabaseUserName = os.Getenv("DB_USERNAME")
	if DatabaseUserName == "" {
		DatabaseUserName = "postgres"
	}
	DatabasePassword = os.Getenv("DB_PASSWORD")
	if DatabasePassword == "" {
		DatabasePassword = "postgres"
	}
	DatabaseName = os.Getenv("DB_NAME")
	if DatabaseName == "" {
		DatabaseName = "auth"
	}
	ServerPort = os.Getenv("SERVER_PORT")
	if ServerPort == "" {
		ServerPort = "8082"
	}
	DatabaseHost = os.Getenv("DB_HOST")
	if DatabaseHost == "" {
		DatabaseHost = "127.0.0.1"
	}
}

func main() {
	log.SetFlags(log.LstdFlags)
	r, err := SetUpRouter()
	if err != nil {
		log.Panicln("Can't set up server:", err)
	}
	log.Println("Starting server on port ", ServerPort)
	r.Run(":" + ServerPort)
}
