package main

import (
	"gtp-tabs/GtpStorage/protocol"
	"log"
	"net"

	"google.golang.org/grpc"
)

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
