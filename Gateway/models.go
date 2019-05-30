package main

import "github.com/gtp-tabs/Gateway/clients"

type postRenderInfo struct {
	Post clients.Post    `json:"post"`
	Tab  clients.TabInfo `json:"tab"`
}
