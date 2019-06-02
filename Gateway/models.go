package main

import "github.com/gtp-tabs/Gateway/clients"

type postRenderInfo struct {
	Post clients.Post    `json:"post"`
	Tab  clients.TabInfo `json:"tab"`
}

type postUpdateRating struct {
	PostID int `json:"post_id"`
	Rating int `json:"rating"`
}

type vkCredentials struct {
	AccessToken string `json:"access_token"`
	UserID      int    `json:"user_id"`
}
