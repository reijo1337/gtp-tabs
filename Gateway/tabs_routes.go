package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gtp-tabs/Gateway/clients"
)

func (ch *clientHolder) getAuthorsByLetter(c *gin.Context) {
	code := c.Param("code")
	result, err := ch.storage.GetAuthorsByLetter(code)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (ch *clientHolder) getAuthorsByName(c *gin.Context) {
	search := c.Param("search")
	if search == "" {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "search required",
			},
		)
		return
	}
	result, err := ch.storage.GetAuthorsByName(search)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (ch *clientHolder) getTabsByName(c *gin.Context) {
	search := c.Param("search")
	if search == "" {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "search required",
			},
		)
		return
	}
	result, err := ch.storage.FindTabsByName(search)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (ch *clientHolder) getMusiciansByGategory(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "category name required",
			},
		)
		return
	}
	result, err := ch.storage.GetAuthorsByCategory(name)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (ch *clientHolder) uploadFile(c *gin.Context) {
	var upload *clients.FileUploadRequest
	if err := c.BindJSON(upload); err != nil {
		log.Printf("binding request body: %v", err)
		c.Status(http.StatusBadRequest)
		return
	}
	err := ch.storage.UploadFile(upload)
	if err != nil {
		log.Printf("sending file: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "can't send file"})
		return
	}
	post := &clients.Post{
		SongName: upload.Song,
	}
	ch.post.SetPost(post)
	c.Status(http.StatusOK)
}

func (ch *clientHolder) downloadFile(c *gin.Context) {
	ret, err := ch.storage.DownloadFile(c.Query("name"))
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}
	c.DataFromReader(http.StatusOK, ret.ContentLength, ret.ContentType, ret.FileContent, ret.ExtraHeaders)
}
