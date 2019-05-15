package main

import (
	"crypto/rsa"
	"fmt"
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type auth struct {
	publicKey *rsa.PublicKey
}

func newAuth(publicKey []byte) (*auth, error) {
	key, err := jwt.ParseRSAPublicKeyFromPEM(publicKey)
	if err != nil {
		return nil, err
	}
	return &auth{
		publicKey: key,
	}, nil
}

func (a *auth) verifyToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Query("access_token")
		if tokenString == "" {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{
					"error": "Unauthorized",
				},
			)
		}
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			return a.publicKey, nil
		})

		if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Next()
		} else {
			log.Println("Gateway: Authorization failed: ", err.Error())
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{
					"error": "Unauthorized",
				},
			)
		}
		c.Next()
	}
}
