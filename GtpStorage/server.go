package main

import (
	"gtp-tabs/GtpStorage/protocol"
	"log"
)

// Server структура для grpc сервера
type Server struct {
	db *Database
}

// MakeServer возвращает новый объект Server, который представляет определения для grpc
func MakeServer() (*Server, error) {
	log.Println("Server: Set up book service...")
	db, err := SetUpDatabase()
	if err != nil {
		return nil, err
	}
	return &Server{db: db}, nil
}

// GetAuthorsByLetter возвращает список музыкантов и количество их исполнителей через поиск по первой букве
func (s *Server) GetAuthorsByLetter(in *protocol.SearchString, p protocol.Tabs_GetAuthorsByLetterServer) error {
	return nil
}

// GetAuthorsByName возвращает список музыкантов и количество их исполнителей через поиск по подстроке
func (s *Server) GetAuthorsByName(in *protocol.SearchString, p protocol.Tabs_GetAuthorsByNameServer) error {
	return nil
}

// FindTabsByName возвращает список табулатур и количество их исполнителей через поиск по подстроке
func (s *Server) FindTabsByName(in *protocol.SearchString, p protocol.Tabs_FindTabsByNameServer) error {
	return nil
}
