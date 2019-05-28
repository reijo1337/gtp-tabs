package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gtp-tabs/Gateway/clients"
)

type clientHolder struct {
	storage clients.StorageClientInterface
	auth    clients.AuthClientInterface
	profile clients.ProfileClientInterface
}

func setUpClientHolder() (*clientHolder, error) {
	config, err := parseConfig("GATEWAY")
	if err != nil {
		return nil, err
	}

	return &clientHolder{
		storage: clients.MakeStorageClient(config.Storage.Host, config.Storage.Port),
		auth:    clients.MakeAuthClient(config.Storage.Host, config.Storage.Port),
		profile: clients.MakeProfileClient(config.Storage.Host, config.Storage.Port),
	}, nil
}

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

func (ch *clientHolder) getToken(c *gin.Context) {
	var user clients.User
	if err := c.BindJSON(&user); err != nil {
		log.Printf("parsing request: %v", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	res, err := ch.auth.GenToken(&user)
	if err != nil {
		log.Println("getting tokens: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "can't get token"})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (ch *clientHolder) getTokenVK(c *gin.Context) {
	var user clients.VkUser
	if err := c.BindJSON(&user); err != nil {
		log.Printf("parsing request: %v", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	res, err := ch.auth.GenTokenVk(&user)
	if err != nil {
		log.Println("getting tokens: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "can't get token"})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (ch *clientHolder) refreshToken(c *gin.Context) {
	tokenString := c.Query("refresh_token")
	if tokenString == "" {
		c.AbortWithStatus(http.StatusBadRequest)
	}
	tokens, err := ch.auth.RefreshToken(tokenString)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	c.JSON(http.StatusOK, tokens)
}

func (ch *clientHolder) register(c *gin.Context) {
	var user clients.User
	if err := c.BindJSON(&user); err != nil {
		log.Printf("parsing request: %v", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	res, err := ch.auth.Register(&user)
	if err != nil {
		log.Println("getting tokens: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "can't get token"})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (ch *clientHolder) registerVk(c *gin.Context) {
	var user clients.VkUser
	if err := c.BindJSON(&user); err != nil {
		log.Printf("parsing request: %v", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	res, err := ch.auth.RegisterVk(&user)
	if err != nil {
		log.Println("getting tokens: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "can't get token"})
		return
	}
	c.JSON(http.StatusOK, res)
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
	authorized.PUT("/", ch.uploadFile)

	r.GET("/alph/:code", ch.getAuthorsByLetter)
	r.GET("/musicians/:search", ch.getAuthorsByName)
	r.GET("/tabs/:search", ch.getTabsByName)
	r.GET("/category/:name", ch.getMusiciansByGategory)
	r.GET("/file", ch.downloadFile)

	r.POST("/", ch.getToken)
	r.POST("/vk", ch.getTokenVK)
	r.GET("/", ch.refreshToken)
	reg := r.Group("/register")
	reg.POST("/", ch.register)
	reg.POST("/vk", ch.registerVk)
	// authorized.GET("/getUserArrears", getUserArrears)
	// authorized.POST("/arrear", newArear)
	// authorized.DELETE("/arrear", closeArrear)
	// authorized.GET("/freeBooks", freeBooks)
	// authorized.OPTIONS("/arrear", func(c *gin.Context) {
	// c.JSON(http.StatusOK, "")
	// })
	// authorized.OPTIONS("/freeBooks", func(c *gin.Context) {
	// c.JSON(http.StatusOK, "")
	// })

	// r.POST("/auth", Login)
	// r.OPTIONS("/auth", func(c *gin.Context) {
	// c.JSON(http.StatusOK, "")
	// })
	// r.GET("/auth", Refresh)

	return r, nil
}
