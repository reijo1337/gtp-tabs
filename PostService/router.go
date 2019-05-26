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

func setUpRouter(dbSource string) (*gin.Engine, error) {
	r := gin.Default()
	db, err := setUpDatabase(dbSource)
	if err != nil {
		return nil, fmt.Errorf("database setup: %v", err)
	}
	s, err := makeService(db)
	if err != nil {
		return nil, fmt.Errorf("make service: %v", err)
	}
	r.GET("/post/:id", s.getPost)
	r.POST("/post", s.setPost)
	r.PUT("/post", s.updateRating)
	r.POST("/comment/:id", s.makeComment)
	return r, nil
}

func (s *service) getPost(c *gin.Context) {
	tabID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("parsing id: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	postInfo, err := s.db.getPost(tabID)
	if err != nil {
		log.Printf("getting post: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "can't get post"})
		return
	}
	c.JSON(http.StatusOK, postInfo)
}

func (s *service) setPost(c *gin.Context) {
	var postInfo post
	if err := c.BindJSON(&postInfo); err != nil {
		log.Printf("parsing body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	if err := s.db.makePost(&postInfo); err != nil {
		log.Printf("making post: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "can't make post"})
		return
	}
	c.JSON(http.StatusOK, postInfo)
}

func (s *service) updateRating(c *gin.Context) {
	var request updateRatingRequest
	if err := c.BindJSON(&request); err != nil {
		log.Printf("parsing body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	if err := s.db.changeRating(request.PostID, request.Rating); err != nil {
		log.Printf("updating rating: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "can't update rating"})
		return
	}
}

func (s *service) makeComment(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("parsing id: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var com comment
	if err := s.db.makeComment(postID, &com); err != nil {
		log.Printf("making comment: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "can't make comment"})
		return
	}
	c.JSON(http.StatusOK, com)
}
