package main

import "time"

type userInfo struct {
	ID          int          `json:"id"`
	AccountID   int          `json:"account_id"`
	Name        string       `json:"name"`
	Registered  time.Time    `json:"registered"`
	Birthday    time.Time    `json:"birhday"`
	Instruments []instrument `json:"instruments"`
}

type instrument struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
