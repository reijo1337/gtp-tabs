package main

type Role struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

type user struct {
	ID       int32
	Login    string `json:"login"`
	Password string `json:"password"`
	Role     Role   `json:"role"`
}

type tokens struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
