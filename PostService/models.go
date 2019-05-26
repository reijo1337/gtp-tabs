package main

type post struct {
	ID       int       `json:"id"`
	SongName string    `json:"song_name"`
	TabID    int       `json:"tab_id"`
	AuthorID int       `json:"author_id"`
	Rating   float32   `json:"rating"`
	Comments []comment `json:"comments"`
}

type comment struct {
	ID       int    `json:"id"`
	AuthorID int    `json:"author_id"`
	Content  string `json:"content"`
}

type updateRatingRequest struct {
	PostID int `json:"post_id"`
	Rating int `json:"rating"`
}
