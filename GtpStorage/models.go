package main

// MusiciansWithCount пресдтавление инфорамии об авторах и количеству табулатур
type MusiciansWithCount struct {
	ID    int32  `json:"id"`
	Name  string `json:"name"`
	Count int32  `json:"count"`
}

// TabWithSize представление информации о табулатуре с ее размером
type TabWithSize struct {
	Musician string `json:"musician"`
	Name     string `json:"name"`
	Size     int64  `json:"size"`
}

type musician struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

type fileUploadRequest struct {
	Filename string `json:"filename"`
	Song     string `json:"song"`
	Musician string `json:"musician"`
	Category string `json:"category"`
	Content  string `json:"content"`
}

type category struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}
