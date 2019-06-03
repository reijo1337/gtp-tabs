package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gtp-tabs/Gateway/clients"
)

type clientHolder struct {
	storage clients.StorageClientInterface
	auth    clients.AuthClientInterface
	profile clients.ProfileClientInterface
	post    clients.PostClientInterface
}

func setUpClientHolder() (*clientHolder, error) {
	config, err := parseConfig("")
	if err != nil {
		return nil, err
	}

	return &clientHolder{
		storage: clients.MakeStorageClient(config.URL.Storage),
		auth:    clients.MakeAuthClient(config.URL.Auth),
		profile: clients.MakeProfileClient(config.URL.Profile),
		post:    clients.MakePostClient(config.URL.Post),
	}, nil
}

func setUpRouter(publicKey []byte) (*gin.Engine, error) {
	r := gin.Default()
	ch, err := setUpClientHolder()
	if err != nil {
		return nil, err
	}
	auth, err := newAuth(publicKey)
	if err != nil {
		return nil, err
	}
	authorized := r.Group("/", auth.verifyToken())
	authorized.POST("/file", ch.uploadFile)
	authorized.OPTIONS("/file", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	r.GET("/profile/:id", ch.getProfile)
	authorized.POST("/rating", ch.updateRating)
	authorized.OPTIONS("/rating", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	authorized.POST("/post/:id", ch.addComment)
	authorized.OPTIONS("/post/:id", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	r.GET("/alph/:code", ch.getAuthorsByLetter)
	r.GET("/musicians/:search", ch.getAuthorsByName)
	r.GET("/tabs/:search", ch.getTabsByName)
	r.GET("/category/:name", ch.getMusiciansByGategory)
	r.GET("/file", ch.downloadFile)
	r.GET("/musician/:id", ch.getTabsByMusicianID)

	r.GET("/post/:id", ch.getPost)

	r.POST("/", ch.getToken)
	r.POST("/vk", ch.getTokenVK)
	r.OPTIONS("/vk", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	r.GET("/", ch.refreshToken)

	reg := r.Group("/register")
	reg.POST("/", ch.register)
	reg.OPTIONS("/", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	reg.POST("/vk", ch.registerVk)
	reg.OPTIONS("/vk", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	r.POST("/auth", ch.login)
	r.OPTIONS("/auth", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	r.GET("/auth/vk", ch.authUser)
	r.GET("/auth", ch.refresh)

	return r, nil
}
