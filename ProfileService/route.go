package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type service struct {
	db *database
}

func makeService(db *database) (*service, error) {
	return &service{db: db}, nil
}

func setUpRouter() (*gin.Engine, error) {
	r := gin.Default()
	db, err := setUpDatabase()
	if err != nil {
		return nil, fmt.Errorf("database setup: %v", err)
	}
	s, err := makeService(db)
	if err != nil {
		return nil, fmt.Errorf("make service: %v", err)
	}
	r.GET("/profile/:id", s.getProfile)
	// r.GET("/letter/musicians/:code", s.GetAuthorsByLetter)
	// r.GET("/musicians/:search", s.GetAuthorsByName)
	// r.GET("/tabs/:search", s.FindTabsByName)
	// r.GET("/category/:name", s.GetAuthorsByCategory)
	// r.PUT("/file", s.Upload)
	// r.GET("/file", s.Download)
	return r, nil
}

func (s *service) getProfile(c *gin.Context) {

}
