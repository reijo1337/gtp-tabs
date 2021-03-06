package main

import (
	"log"
	"net/http"
	"strconv"

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

func (ch *clientHolder) getTabsByMusicianID(c *gin.Context) {
	search, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("getting id: %v", err)
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "invalid id",
			},
		)
		return
	}
	result, err := ch.storage.FindTabsByMusicianID(search)
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
	upload := &clients.FileUploadRequest{}
	if err := c.BindJSON(upload); err != nil {
		log.Printf("binding request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "can't get request"})
		return
	}
	tab, err := ch.storage.UploadFile(upload)
	if err != nil {
		log.Printf("sending file: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "can't send file"})
		return
	}
	post := &clients.Post{
		SongName: upload.Song,
		TabID:    tab.ID,
		AuthorID: tab.Author.ID,
	}
	if err := ch.post.SetPost(post); err != nil {
		log.Printf("making post: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Header("Access-Control-Allow-Origin", "*")
	c.JSON(http.StatusOK, &postRenderInfo{
		Tab:  *tab,
		Post: *post,
	})
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
