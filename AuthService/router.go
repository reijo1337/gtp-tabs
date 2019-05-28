package main

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type service struct {
	db                *Database
	privateKey        *rsa.PrivateKey
	accessExpiration  time.Duration
	refreshExpiration time.Duration
}

func makeService(db *Database, privateKeyLoc string, accExp time.Duration, refExp time.Duration) (*service, error) {
	privateKeyBytes, err := ioutil.ReadFile(privateKeyLoc)
	if err != nil {
		return nil, err
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("parsing private key: %v", err)
	}
	return &service{
		db:                db,
		privateKey:        privateKey,
		accessExpiration:  accExp,
		refreshExpiration: refExp,
	}, nil
}

// SetUpRouter утсановка методов на прослушку
func SetUpRouter(privateKeyLoc string, accExp time.Duration, refExp time.Duration) (*gin.Engine, error) {
	r := gin.Default()
	db, err := SetUpDatabase()
	if err != nil {
		return nil, err
	}
	s, err := makeService(db, privateKeyLoc, accExp, refExp)
	if err != nil {
		return nil, err
	}
	r.POST("/", s.getToken)
	r.POST("/vk", s.getTokenVK)
	r.GET("/", s.refreshToken)
	reg := r.Group("/register")
	reg.POST("/", s.register)
	reg.POST("/vk", s.registerVk)
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
		token, err := s.genToken(req.Login)
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
			loginResponse{
				userID: req.ID,
				tokens: *token,
			},
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
		token, err := s.genToken(strconv.FormatInt(req.UserID, 10))
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
			loginResponse{
				userID: req.ID,
				tokens: *token,
			},
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
		newTokens, err := s.genToken(claims["login"].(string))
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

func (s *service) genToken(login string) (*tokens, error) {
	AccessTokenExp := time.Now().Add(s.accessExpiration).Unix()
	RefreshTokenExp := time.Now().Add(s.refreshExpiration).Unix()
	accesToken := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"login": login,
		"iss":   iss,
		"exp":   AccessTokenExp,
		"aud":   aud,
	})
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"login": login,
		"iss":   iss,
		"exp":   RefreshTokenExp,
		"aud":   aud,
	})

	accessTokenString, err := accesToken.SignedString(s.privateKey)
	if err != nil {
		return nil, err
	}
	refreshTokenString, err := refreshToken.SignedString(s.privateKey)

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
	newUser, err := s.db.insertNewUser(req.Login, req.Password, req.Role.Name)
	if err != nil {
		log.Println("Server: can't regiser user", err)
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": "Не удалось зарегистрировать пользователя",
			},
		)
	}
	token, err := s.genToken(req.Login)
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
		gin.H{
			"user":  *newUser,
			"token": token,
		},
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
	newUser, err := s.db.insertNewVkUser(req.UserID, req.Role.Name)
	if err != nil {
		log.Println("Server: can't regiser user", err)
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": "Не удалось зарегистрировать пользователя",
			},
		)
	}
	token, err := s.genToken(strconv.FormatInt(req.UserID, 10))
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
		gin.H{
			"user":  *newUser,
			"token": token,
		},
	)
}
