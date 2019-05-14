package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type service struct {
	db *Database
}

func makeService(db *Database) (*service, error) {
	log.Println("Server: Set up auth service...")
	return &service{db: db}, nil
}

// SetUpRouter утсановка методов на прослушку
func SetUpRouter() (*gin.Engine, error) {
	r := gin.Default()
	db, err := SetUpDatabase()
	if err != nil {
		return nil, err
	}
	s, err := makeService(db)
	if err != nil {
		return nil, err
	}
	r.POST("/", s.getToken)
	r.POST("/vk", s.getTokenVK)
	r.GET("/", s.refreshToken)
	return r, nil
}

func (s *service) getToken(c *gin.Context) {
	log.Println("Server: request for new token for local user")
	req := &user{}
	if err := c.BindJSON(req); err != nil {
		log.Println("Server: Can't parse request body:", err.Error())
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "Пробелмы с обработкой запроса",
			},
		)
		return
	}

	log.Println("Server: Checking login ", req.Login)
	if s.db.isAuthorized(req) {
		token, err := genToken(req.Login)
		if err != nil {
			log.Println("Server: Can't authorize this user: ", err.Error())
			c.JSON(
				http.StatusUnauthorized,
				gin.H{
					"error": "Неудачная авторизация",
				},
			)
			return
		}
		c.JSON(
			http.StatusOK,
			token,
		)
	} else {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"error": "Неудачная авторизация",
			},
		)
		return
	}
}

func (s *service) getTokenVK(c *gin.Context) {
	log.Println("Server: request for new token for vk user")
	req := &vkUser{}
	if err := c.BindJSON(req); err != nil {
		log.Println("Server: Can't parse request body:", err.Error())
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "Пробелмы с обработкой запроса",
			},
		)
		return
	}

	log.Println("Server: Checking login ", req.UserID)
	if s.db.isAuthorizedVK(req) {
		token, err := genToken(strconv.FormatInt(req.UserID, 10))
		if err != nil {
			log.Println("Server: Can't authorize this user: ", err.Error())
			c.JSON(
				http.StatusUnauthorized,
				gin.H{
					"error": "Неудачная авторизация",
				},
			)
			return
		}
		c.JSON(
			http.StatusOK,
			token,
		)
	} else {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"error": "Неудачная авторизация",
			},
		)
		return
	}
}

func (s *service) refreshToken(c *gin.Context) {
	log.Println("Server: Request refresh")
	tokenString := c.Query("refresh_token")
	if tokenString == "" {
		c.AbortWithStatus(http.StatusBadRequest)
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// hmacSampleSecret := os.Getenv("SECRET")
		hmacSampleSecret := []byte("secc")
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return hmacSampleSecret, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		newTokens, err := genToken(claims["login"].(string))
		if err != nil {
			log.Println("Server: Can't authorize this user: ", err.Error())
			c.JSON(
				http.StatusUnauthorized,
				gin.H{
					"error": "Неудачная авторизация",
				},
			)
			return
		}
		c.JSON(
			http.StatusOK,
			newTokens,
		)
	} else {
		log.Println("Gateway: Authorization failed: ", err.Error())
		c.AbortWithStatusJSON(
			http.StatusUnauthorized,
			gin.H{
				"error": "Неудачная авторизация",
			},
		)
	}
}

func genToken(login string) (*tokens, error) {
	log.Println("Server: Generating token")
	// hmacSampleSecret := os.Getenv("SECRET")
	hmacSampleSecret := []byte("secc")
	AccessTokenExp := time.Now().Add(time.Second * 30).Unix()
	RefreshTokenExp := time.Now().Add(time.Hour * 24).Unix()
	log.Println("Server: Gen access token")
	accesToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"login": login,
		"iss":   iss,
		"exp":   AccessTokenExp,
		"aud":   aud,
	})
	log.Println("Server: Gen refresh token")
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"login": login,
		"iss":   iss,
		"exp":   RefreshTokenExp,
		"aud":   aud,
	})

	log.Println("Server: Signing access token", accesToken, hmacSampleSecret)
	accessTokenString, err := accesToken.SignedString(hmacSampleSecret)
	if err != nil {
		return nil, err
	}
	log.Println("Server: Signing refresh token", refreshToken, hmacSampleSecret)
	refreshTokenString, err := refreshToken.SignedString(hmacSampleSecret)

	return &tokens{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
	}, nil
}

func (s *service) register(c *gin.Context) {
	log.Println("Server: request for registration for local user")
	req := &user{}
	if err := c.BindJSON(req); err != nil {
		log.Println("Server: Can't parse request body:", err.Error())
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "Пробелмы с обработкой запроса",
			},
		)
		return
	}
	log.Println("Server: Checking login ", req.Login)
	if s.db.credentialsInUse(req) {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "Пользователь с таким логином уже существует",
			},
		)
		return
	}
	_, err := s.db.insertNewUser(req.Login, req.Password, req.Role.Name)
	if err != nil {
		log.Println("Server: can't regiser user", err)
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": "Не удалось зарегистрировать пользователя",
			},
		)
	}
	token, err := genToken(req.Login)
	if err != nil {
		log.Println("Server: Can't authorize this user: ", err.Error())
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": "Неудачная авторизация",
			},
		)
		return
	}
	c.JSON(
		http.StatusOK,
		token,
	)
}

func (s *service) registerVk(c *gin.Context) {
	log.Println("Server: request for registration for vk user")
	req := &vkUser{}
	if err := c.BindJSON(req); err != nil {
		log.Println("Server: Can't parse request body:", err.Error())
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "Пробелмы с обработкой запроса",
			},
		)
		return
	}
	log.Println("Server: Checking login ", req.UserID)
	if s.db.isAuthorizedVK(req) {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "Пользователь с таким логином уже существует",
			},
		)
		return
	}
	_, err := s.db.insertNewVkUser(req.UserID, req.Role.Name)
	if err != nil {
		log.Println("Server: can't regiser user", err)
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": "Не удалось зарегистрировать пользователя",
			},
		)
	}
	token, err := genToken(strconv.FormatInt(req.UserID, 10))
	if err != nil {
		log.Println("Server: Can't authorize this user: ", err.Error())
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": "Неудачная авторизация",
			},
		)
		return
	}
	c.JSON(
		http.StatusOK,
		token,
	)
}
