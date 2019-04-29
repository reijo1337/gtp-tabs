package main

type user struct {
	ID       int32
	Login    string `json:"login"`
	Password string `json:"password"`
}

type tokens struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}