package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gtp-tabs/Gateway/clients"
)

func (ch *clientHolder) getPost(c *gin.Context) {
	tabID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("parsing id: %v", err)
		c.Status(http.StatusBadRequest)
		return
	}
	post, err := ch.post.GetPost(tabID)
	if err != nil {
		log.Printf("getting post: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	tab, err := ch.storage.GetTab(post.TabID)
	if err != nil {
		log.Printf("getting tab: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, postRenderInfo{
		Post: *post,
		Tab:  *tab,
	})
}

func (ch *clientHolder) updateRating(c *gin.Context) {
	var request postUpdateRating
	if err := c.BindJSON(&request); err != nil {
		log.Printf("bind json: %v", err)
		c.Status(http.StatusBadRequest)
		return
	}
	if err := ch.post.UpdateRating(request.PostID, request.Rating); err != nil {
		log.Printf("update rating: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusOK)
}

func (ch *clientHolder) addComment(c *gin.Context) {
	var request clients.Comment
	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("parsing id: %v", err)
		c.Status(http.StatusBadRequest)
		return
	}
	if err := c.BindJSON(&request); err != nil {
		log.Printf("bind json: %v", err)
		c.Status(http.StatusBadRequest)
		return
	}
	if err := ch.post.AddComment(postID, &request); err != nil {
		log.Printf("update rating: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, request)
}
