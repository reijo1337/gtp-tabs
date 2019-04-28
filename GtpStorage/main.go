package main

import (
	"gtp-tabs/GtpStorage/protocol"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
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
		DatabaseName = "postgres"
	}
	ServerPort = os.Getenv("SERVER_PORT")
	if ServerPort == "" {
		ServerPort = "8081"
	}
	DatabaseHost = os.Getenv("DB_HOST")
	if DatabaseHost == "" {
		DatabaseHost = "127.0.0.1"
	}
}

func main() {
	log.SetFlags(log.LstdFlags)
	lis, err := net.Listen("tcp", ":"+ServerPort)
	if err != nil {
		log.Printf("Main: Can't start server on port %s: %s", ServerPort, err)
	}

	db, err := SetUpDatabase()
	if err != nil {
		log.Println("Main: Can't  setup database", err)
	}

	serv, err := MakeServer(db)
	if err != nil {
		log.Println("Main: Can't  start server", err)
	}

	server := grpc.NewServer()

	protocol.RegisterTabsServer(server, serv)
	log.Println("Main: Starting server at port", ServerPort)
	server.Serve(lis)
}
