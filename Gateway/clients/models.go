package clients

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
