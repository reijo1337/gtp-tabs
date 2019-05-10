package main

import (
	// "fmt"
	// jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"

	"github.com/gtp-tabs/Gateway/clients"
)

type clientHolder struct {
	storage clients.StorageClientInterface
}

func setUpClientHolder() (*clientHolder, error) {
	config, err := parseConfig("GATEWAY")
	if err != nil {
		return nil, err
	}

	storage := clients.MakeStorageClient(config.Storage.Host, config.Storage.Port)
	return &clientHolder{
		storage: storage,
	}, nil
}

func authRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// log.Println("Gateway: New authorized request")
		// tokenString := c.Query("access_token")
		// if tokenString == "" {
		// 	c.AbortWithStatusJSON(
		// 		http.StatusUnauthorized,
		// 		gin.H{
		// 			"error": "Unauthorized",
		// 		},
		// 	)
		// }
		// token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 	// hmacSampleSecret := os.Getenv("SECRET")
		// 	hmacSampleSecret := []byte("secc")
		// 	// Don't forget to validate the alg is what you expect:
		// 	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		// 		return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		// 	}

		// 	// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		// 	return hmacSampleSecret, nil
		// })

		// if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// 	c.Next()
		// } else {
		// 	log.Println("Gateway: Authorization failed: ", err.Error())
		// 	c.AbortWithStatusJSON(
		// 		http.StatusUnauthorized,
		// 		gin.H{
		// 			"error": "Unauthorized",
		// 		},
		// 	)
		// }
		c.Next()
	}
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

func setUpRouter() (*gin.Engine, error) {
	r := gin.Default()
	ch, err := setUpClientHolder()
	if err != nil {
		return nil, err
	}
	authorized := r.Group("/", authRequired())
	authorized.PUT("/", ch.uploadFile)
	r.GET("/alph/:code", ch.getAuthorsByLetter)
	r.GET("/musicians/:search", ch.getAuthorsByName)
	r.GET("/tabs/:search", ch.getTabsByName)
	r.GET("/category/:name", ch.getMusiciansByGategory)
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

func main() {
	log.SetFlags(log.LstdFlags)
	config, err := parseConfig("GATEWAY")
	if err != nil {
		log.Panicln("Can't read config:", err)
	}
	r, err := setUpRouter()
	if err != nil {
		log.Panicln("Can't set up router:", err)
	}
	r.Run(":" + config.Port)
}
