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
func MakeServer(db *Database) (*Server, error) {
	log.Println("Server: Set up book service...")
	return &Server{db: db}, nil
}

// GetAuthorsByLetter возвращает список музыкантов и количество их исполнителей через поиск по первой букве
func (s *Server) GetAuthorsByLetter(in *protocol.SearchString, p protocol.Tabs_GetAuthorsByLetterServer) error {
	log.Println("New request for searching musicians by letter", in.GetSearch())
	result, err := s.db.getMusiciansByLetter(in.GetSearch())
	if err != nil {
		log.Println("Can't get musicians by letter", in.GetSearch(), "from database.", err)
		return err
	}

	for _, res := range result {
		mwc := &protocol.MusicianWithCount{
			Name:  res.Name,
			Count: res.Count,
		}
		if err := p.Send(mwc); err == nil {
			log.Println("Can't send info about musician", res.Name)
			return err
		}
	}
	log.Println("Request for searcing musicians by letter", in.GetSearch(), "processes succsesfuly")
	return nil
}

// GetAuthorsByName возвращает список музыкантов и количество их исполнителей через поиск по подстроке
func (s *Server) GetAuthorsByName(in *protocol.SearchString, p protocol.Tabs_GetAuthorsByNameServer) error {
	log.Println("New request for searching musicians by substing", in.GetSearch())
	result, err := s.db.getMusicians(in.GetSearch())
	if err != nil {
		log.Println("Can't get musicians by substring", in.GetSearch(), "from database.", err)
		return err
	}

	for _, res := range result {
		mwc := &protocol.MusicianWithCount{
			Name:  res.Name,
			Count: res.Count,
		}
		if err := p.Send(mwc); err == nil {
			log.Println("Can't send info about musician", res.Name, err)
			return err
		}
	}
	log.Println("Request for searcing musicians by letter", in.GetSearch(), "processes succsesfuly")
	return nil
}

// FindTabsByName возвращает список табулатур и количество их исполнителей через поиск по подстроке
func (s *Server) FindTabsByName(in *protocol.SearchString, p protocol.Tabs_FindTabsByNameServer) error {
	log.Println("New request for searching tabs by substring", in.GetSearch())

	return nil
}
