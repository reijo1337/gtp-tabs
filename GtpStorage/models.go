package main

// MusiciansWithCount пресдтавление инфорамии об авторах и количеству табулатур
type MusiciansWithCount struct {
	ID    int32
	Name  string
	Count int32
}

// TabWithSize представление информации о табулатуре с ее размером
type TabWithSize struct {
	Musician string
	Name     string
	Size     int32
}
