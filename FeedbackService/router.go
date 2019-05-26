package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	gomail "gopkg.in/gomail.v2"
)

type service struct {
	db       *database
	url      string
	port     int
	login    string
	password string
}

func makeService(db *database, url, login, password string, port int) (*service, error) {
	return &service{
		db:       db,
		url:      url,
		port:     port,
		login:    login,
		password: password,
	}, nil
}

func setUpRouter(dbSource, url, login, password string, port int) (*gin.Engine, error) {
	r := gin.Default()
	db, err := setUpDatabase(dbSource)
	if err != nil {
		return nil, fmt.Errorf("database setup: %v", err)
	}
	s, err := makeService(db, url, login, password, port)
	if err != nil {
		return nil, fmt.Errorf("make service: %v", err)
	}
	r.GET("/feedback", s.getFeedbacks)
	r.POST("/feedback", s.addFeedback)
	r.POST("/response", s.sendResponse)
	return r, nil
}

func (s *service) getFeedbacks(c *gin.Context) {
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		log.Printf("parsing limit: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit"})
		return
	}
	offset, err := strconv.Atoi(c.Query("offset"))
	if err != nil {
		log.Printf("parsing offset: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid offset"})
		return
	}
	fbs, err := s.db.getFeedbacks(limit, offset)
	if err != nil {
		log.Printf("parsing offset: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "can't get feedbck"})
		return
	}
	c.JSON(http.StatusOK, fbs)
}

func (s *service) addFeedback(c *gin.Context) {
	var fb feedback
	if err := c.BindJSON(&fb); err != nil {
		log.Printf("parsing body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	if err := s.db.addFeedback(fb); err != nil {
		log.Printf("making feedback: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "can't make feedback"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": "ok"})
}

func (s *service) sendResponse(c *gin.Context) {
	var fb response
	if err := c.BindJSON(&fb); err != nil {
		log.Printf("parsing body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	m := gomail.NewMessage()
	m.SetHeader("From", "tantsevov@yandex.ru")
	m.SetHeader("To", fb.Feedback.Email)
	m.SetHeader("Subject", fb.Feedback.Username+". Вам ответ от админки")
	m.SetBody("text/html", fb.Response)

	d := gomail.NewDialer(s.url, s.port, s.login, s.password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		log.Printf("can't send feedback response: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "can't send feedback response"})
	}

	if err := s.db.removeFeedback(fb.Feedback.ID); err != nil {
		log.Printf("can't remove feedback: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "can't remove feedback"})
	}
	c.Status(http.StatusOK)
}
