package main

type feedback struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Message  string `json:"message"`
}

type response struct {
	Feedback feedback `json:"feedback"`
	Response string   `json:"message"`
}
