package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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
	c.JSON(http.StatusOK, post)
}
