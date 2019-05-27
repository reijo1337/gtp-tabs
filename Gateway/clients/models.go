package clients

import "io"

// MusiciansWithCount пресдтавление инфорамии об авторах и количеству табулатур
type MusiciansWithCount struct {
	ID    int32  `json:"id"`
	Name  string `json:"name"`
	Count int32  `json:"count"`
}

// TabWithSize представление информации о табулатуре с ее размером
type TabWithSize struct {
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
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

type User struct {
	ID       int32
	Login    string `json:"login"`
	Password string `json:"password"`
	Role     Role   `json:"role"`
}

type VkUser struct {
	ID     int64
	UserID int64 `json:"user_id"`
	Role   Role  `json:"role"`
}

type Tokens struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
