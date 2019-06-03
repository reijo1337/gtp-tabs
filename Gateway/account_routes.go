package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gtp-tabs/Gateway/clients"
)

func (ch *clientHolder) getToken(c *gin.Context) {
	var user clients.User
	if err := c.BindJSON(&user); err != nil {
		log.Printf("parsing request: %v", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	res, err := ch.auth.GenToken(&user)
	if err != nil {
		log.Printf("getting tokens: %v", err)
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
		log.Printf("getting tokens: %v", err)
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
		log.Printf("getting tokens: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "can't get token"})
		return
	}
	profile := &clients.ProfileInfo{
		AccountID:  res.User.ID,
		Name:       res.User.Login,
		Registered: time.Now(),
	}
	if err := ch.profile.SetProfile(profile); err != nil {
		log.Printf("making profile: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "can't make profile"})
		return
	}
	res.Profile = *profile
	c.JSON(http.StatusOK, res)
}

func (ch *clientHolder) getProfile(c *gin.Context) {
	profileID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("parsing id: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad profile id"})
		return
	}
	res, err := ch.profile.GetProfile(profileID)
	if err != nil {
		log.Printf("getting profile: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "bad profile id"})
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
		log.Printf("getting tokens: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "can't get token"})
		return
	}
	profile := &clients.ProfileInfo{
		AccountID:  res.User.ID,
		Name:       "",
		Registered: time.Now(),
	}
	if err := ch.profile.SetProfile(profile); err != nil {
		log.Printf("making profile: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "can't make profile"})
		return
	}
	res.Profile = *profile
	c.JSON(http.StatusOK, res)
}

func (ch *clientHolder) getUser(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("parsing user id: %v", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	userProfile, err := ch.profile.GetProfile(userID)
	if err != nil {
		log.Printf("getting user profile: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, userProfile)
}

func (ch *clientHolder) login(c *gin.Context) {
	var user clients.User
	if err := c.BindJSON(&user); err != nil {
		log.Printf("parsing request body: %v", err)
		c.Status(http.StatusBadRequest)
		return
	}
	tokens, err := ch.auth.GenToken(&user)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	profile, err := ch.profile.GetProfileByAcc(tokens.UserID)
	if err != nil {
		log.Printf("getting user profile: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	log.Println(profile)
	c.JSON(http.StatusOK, clients.LoginResponse{
		ProfileID: profile.ID,
		Tokens:    tokens.Tokens})
}

func (ch *clientHolder) loginVk(c *gin.Context) {
	var user clients.VkUser
	if err := c.BindJSON(&user); err != nil {
		log.Printf("parsing request body: %v", err)
		c.Status(http.StatusBadRequest)
		return
	}
	tokens, err := ch.auth.GenTokenVk(&user)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	profile, err := ch.profile.GetProfileByAcc(tokens.UserID)
	if err != nil {
		log.Printf("getting user profile: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	log.Println(profile)
	c.JSON(http.StatusOK, clients.LoginResponse{
		ProfileID: profile.ID,
		Tokens:    tokens.Tokens})
}

func (ch *clientHolder) refresh(c *gin.Context) {
	tokens, err := ch.auth.RefreshToken(c.Query("refresh_token"))
	if err != nil {
		log.Printf("refresh token: %v", err)
		c.Status(http.StatusBadRequest)
		return
	}
	c.JSON(
		http.StatusOK,
		tokens,
	)
}

func (ch *clientHolder) authUser(c *gin.Context) {
	code := c.Query("code")
	req, err := http.NewRequest("GET", "https://oauth.vk.com/access_token", nil)
	if err != nil {
		log.Print(err)
		return
	}

	q := req.URL.Query()
	q.Add("client_id", "6978271")
	q.Add("client_secret", "DKdWDto5gJU4ViGrJW4d")
	q.Add("redirect_uri", "http://localhost:3000/vk")
	q.Add("code", code)

	req.URL.RawQuery = q.Encode()
	resp, err := http.Get(req.URL.String())
	if err != nil {
		log.Println(fmt.Errorf("send vk auth request: %v", err))
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(fmt.Errorf("reading response: %v", err))
		return
	}
	var respVk vkCredentials
	if err = json.Unmarshal(body, &respVk); err != nil {
		log.Println(fmt.Errorf("parsing response: %v", err))
		return
	}

	vkU := clients.VkUser{
		UserID: respVk.UserID,
		Role: clients.Role{
			Name: "user",
		},
	}

	vkAu, err := ch.auth.GenTokenVk(&vkU)
	if err != nil {
		res, err := ch.auth.RegisterVk(&vkU)
		if err != nil {
			log.Println("register: ", err)
			return
		}
		profile := &clients.ProfileInfo{
			AccountID:  res.User.ID,
			Name:       "VK" + strconv.Itoa(res.User.ID),
			Registered: time.Now(),
		}
		if err := ch.profile.SetProfile(profile); err != nil {
			log.Printf("making profile: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "can't make profile"})
			return
		}
		res.Profile = *profile
		log.Println(*res)
		c.JSON(http.StatusOK, clients.LoginResponse{
			ProfileID: res.Profile.ID,
			Tokens:    res.Tokens})
		return
	}
	profile, err := ch.profile.GetProfileByAcc(vkAu.UserID)
	if err != nil {
		log.Printf("getting user profile: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, clients.LoginResponse{
		ProfileID: profile.ID,
		Tokens:    vkAu.Tokens})
}
