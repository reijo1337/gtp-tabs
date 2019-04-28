package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type service struct {
	db *Database
}

func makeService(db *Database) (*service, error) {
	log.Println("Server: Set up book service...")
	return &service{db: db}, nil
}

// SetUpRouter утсановка методов на прослушку
func SetUpRouter() (*gin.Engine, error) {
	r := gin.Default()
	db, err := SetUpDatabase()
	if err != nil {
		return nil, err
	}
	s, err := makeService(db)
	if err != nil {
		return nil, err
	}
	r.GET("/authors_by_letter", s.GetAuthorsByLetter)
	r.GET("/author_by_name", s.GetAuthorsByName)
	return r, nil
}

// GetAuthorsByLetter возвращает список музыкантов и количество их исполнителей через поиск по первой букве
func (s *service) GetAuthorsByLetter(c *gin.Context) {
	letter := c.Query("letter")
	log.Println("New request for searching musicians by letter", letter)
	result, err := s.db.getMusiciansByLetter(letter)
	if err != nil {
		log.Println("Can't get musicians by letter", letter, "from database.", err)
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "Can't get musicians by letter",
			},
		)
		return
	}
	c.JSON(http.StatusOK, result)
}

// GetAuthorsByName возвращает список музыкантов и количество их исполнителей через поиск по подстроке
func (s *service) GetAuthorsByName(c *gin.Context) {
	searchString := c.Query("search")
	log.Println("New request for searching musicians by substing", searchString)
	result, err := s.db.getMusicians(searchString)
	if err != nil {
		log.Println("Can't get musicians by substring", searchString, "from database.", err)
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "Can't get musicians by substring",
			},
		)
		return
	}
	c.JSON(http.StatusOK, result)
}
