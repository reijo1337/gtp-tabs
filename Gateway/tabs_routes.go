package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
	err := ch.storage.UploadFile(c.Request.Body)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}
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
