package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

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
	r.POST("/profile", s.setNewProfile)
	return r, nil
}

func (s *service) getProfile(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println("Can't get user id:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad id"})
		return
	}
	user, err := s.db.getUser(userID)
	if err != nil {
		log.Printf("getting user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "can't get user"})
		return
	}
	instruments, err := s.db.getInstruments(user.ID)
	if err != nil {
		log.Printf("getting instruments: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "can't get instruments"})
		return
	}
	user.Instruments = instruments
	c.JSON(http.StatusOK, user)
}

func (s *service) setNewProfile(c *gin.Context) {
	var newUser userInfo
	err := c.BindJSON(&newUser)
	if err != nil {
		log.Println("Can't get body", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Can't get body"})
		return
	}
	if err := s.db.setNewUser(&newUser); err != nil {
		log.Println("saving user to db:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Can't save user"})
		return
	}
	c.Status(http.StatusOK)
}
