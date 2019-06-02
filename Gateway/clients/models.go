package clients

import (
	"io"
	"time"
)

// MusiciansWithCount пресдтавление инфорамии об авторах и количеству табулатур
type MusiciansWithCount struct {
	ID    int32  `json:"id"`
	Name  string `json:"name"`
	Count int32  `json:"count"`
}

// TabWithSize представление информации о табулатуре с ее размером
type TabWithSize struct {
	ID       int     `json:"id"`
	Musician string  `json:"musician"`
	Name     string  `json:"name"`
	Size     float64 `json:"size"`
}

// ErrorResponse если результат прошел неудачно
type ErrorResponse struct {
	Error string `json:"error"`
}

type FileDownloadResponse struct {
	FileContent   io.Reader
	ContentLength int64
	ContentType   string
	ExtraHeaders  map[string]string
}

type Role struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type User struct {
	ID       int
	Login    string `json:"login"`
	Password string `json:"password"`
	Role     Role   `json:"role"`
}

type VkUser struct {
	ID     int
	UserID int  `json:"user_id"`
	Role   Role `json:"role"`
}

type Tokens struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type ProfileInfo struct {
	ID          int          `json:"id"`
	AccountID   int          `json:"account_id"`
	Name        string       `json:"name"`
	Registered  time.Time    `json:"registered"`
	Birthday    time.Time    `json:"birhday"`
	Instruments []Instrument `json:"instruments"`
}

type Instrument struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type RegisterResponse struct {
	Profile ProfileInfo `json:"profile"`
	User    User        `json:"user"`
	Tokens  Tokens      `json:"tokens"`
}

type RegisterVkResponse struct {
	Profile ProfileInfo `json:"profile"`
	User    VkUser      `json:"user"`
	Tokens  Tokens      `json:"tokens"`
}

type LoginResponse struct {
	ProfileID int    `json:"profile_id"`
	Tokens    Tokens `json:"tokens"`
}

type LoginWithUserResponse struct {
	UserID int    `json:"user_id"`
	Tokens Tokens `json:"tokens"`
}

type Post struct {
	ID       int       `json:"id"`
	SongName string    `json:"song_name"`
	TabID    int       `json:"tab_id"`
	AuthorID int       `json:"author_id"`
	Rating   float32   `json:"rating"`
	Comments []Comment `json:"comments"`
}

type Comment struct {
	ID       int    `json:"id"`
	AuthorID int    `json:"author_id"`
	Content  string `json:"content"`
}

type UpdateRatingRequest struct {
	PostID int `json:"post_id"`
	Rating int `json:"rating"`
}

type FileUploadRequest struct {
	Filename string `json:"filename"`
	Song     string `json:"song"`
	Musician string `json:"musician"`
	Category string `json:"category"`
	Content  string `json:"content"`
}

type TabInfo struct {
	ID     int      `json:"id"`
	Author Musician `json:"musician"`
	Name   string   `json:"name"`
	Size   int64    `json:"size"`
	Cat    Category `json:"category"`
}

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Musician struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
