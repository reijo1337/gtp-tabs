package main

import (
	"encoding/base64"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

type service struct {
	db *Database
}

func makeService(db *Database) (*service, error) {
	log.Println("Server: Set up storage service...")
	return &service{db: db}, nil
}

// SetUpRouter утсановка методов на прослушку
func SetUpRouter() (*gin.Engine, error) {
	r := gin.Default()
	db, err := SetUpDatabase()
	if err != nil {
		return nil, err
	}
	s, err := makeService(db)
	if err != nil {
		return nil, err
	}
	r.GET("/letter/musicians/:code", s.GetAuthorsByLetter)
	r.GET("/musicians/:search", s.GetAuthorsByName)
	r.GET("/tabs/:search", s.FindTabsByName)
	r.GET("/category/:name", s.GetAuthorsByCategory)
	r.GET("/tab/:id", s.getTabByID)
	r.POST("/file", s.Upload)
	r.GET("/file", s.Download)
	return r, nil
}

func (s *service) getTabByID(c *gin.Context) {
	tabID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("getting tab id: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect id"})
		return
	}
	tab, err := s.db.tabByID(tabID)
	if err != nil {
		log.Printf("getting tab by id: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "can't get tab"})
		return
	}
	c.JSON(http.StatusOK, tab)
}

// GetAuthorsByLetter возвращает список музыкантов и количество их исполнителей через поиск по первой букве
func (s *service) GetAuthorsByLetter(c *gin.Context) {
	var (
		err    error
		result []MusiciansWithCount
	)
	letter := c.Param("code")
	log.Println("New request for searching musicians by letter code", letter)
	code, err := strconv.Atoi(letter)
	if err != nil {
		log.Println("Can't get code.", err)
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "Can't get code",
			},
		)
		return
	}
	if code == 0 {
		result, err = s.db.getMusiciansByNumber()
	} else if !(code > 64 && code < 91) && !(code > 1039 && code < 1072) {
		log.Println("Wrong code ", code)
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "Wrong code",
			},
		)
		return
	} else {
		result, err = s.db.getMusiciansByLetter(string(code))
	}
	if err != nil {
		log.Println("Can't get musicians by letter", letter, "from database.", err)
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "Can't get musicians by letter",
			},
		)
		return
	}
	c.JSON(http.StatusOK, result)
}

// GetAuthorsByName возвращает список музыкантов и количество их исполнителей через поиск по подстроке
func (s *service) GetAuthorsByName(c *gin.Context) {
	searchString := c.Param("search")
	log.Println("New request for searching musicians by substing", searchString)
	result, err := s.db.getMusicians(searchString)
	if err != nil {
		log.Println("Can't get musicians by substring", searchString, "from database.", err)
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "Can't get musicians by substring",
			},
		)
		return
	}
	c.JSON(http.StatusOK, result)
}

// FindTabsByName возвращает список табулатур и количество их исполнителей через поиск по подстроке
func (s *service) FindTabsByName(c *gin.Context) {
	searchString := c.Param("search")
	log.Println("New request for searching tabs by substring", searchString)
	results, err := s.db.getTabsByName(searchString)
	if err != nil {
		log.Println("Can't get tabs by substring", searchString, "from database.", err)
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "Can't get tabs by substring",
			},
		)
		return
	}
	c.JSON(http.StatusOK, results)
}

// GetAuthorsByCategory поиск по категориям
func (s *service) GetAuthorsByCategory(c *gin.Context) {
	category := c.Param("name")
	log.Println("New request for searching by category", category)
	results, err := s.db.getMusiciansByCategory(category)
	if err != nil {
		log.Println("Can't get musicians by category", category, "from database.", err)
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "Can't get musicians by category",
			},
		)
		return
	}
	c.JSON(http.StatusOK, results)
}

// Upload загрузка файла на сервер
func (s *service) Upload(c *gin.Context) {
	var req fileUploadRequest
	err := c.BindJSON(&req)
	if err != nil {
		log.Println("Can't get body", err)
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "Can't get body",
			},
		)
		return
	}

	mus, err := s.db.getOrCreateMusician(req.Musician)
	if err != nil {
		log.Println("Can't get musician", err)
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": "Can't get musician",
			},
		)
		return
	}

	cat, err := s.db.getOrCreateCategory(req.Category)
	if err != nil {
		log.Println("Can't get category", err)
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": "Can't get category",
			},
		)
		return
	}

	file, err := os.OpenFile(req.Filename, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Println("Can't save file", err)
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": "Can't save file",
			},
		)
		return
	}
	byteContent, err := base64.StdEncoding.DecodeString(req.Content)
	if err != nil {
		log.Println("Can't get file", err)
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "Can't save file",
			},
		)
		return
	}
	_, err = file.Write(byteContent)
	if err != nil {
		log.Println("Can't save file", err)
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": "Can't save file",
			},
		)
		return
	}
	fi, err := file.Stat()
	if err != nil {
		log.Println("Can't get file size", err)
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": "Can't get file size",
			},
		)
		return
	}
	tabID, err := s.db.createSong(mus.ID, cat.ID, req.Song, fi.Size())
	if err != nil {
		log.Println("Can't save info about file", err)
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": "Can't save info about file",
			},
		)
		return
	}
	tab := &tabInfo{
		ID:     tabID,
		Author: mus,
		Name:   req.Song,
		Cat:    cat,
		Size:   fi.Size(),
	}
	c.JSON(http.StatusOK, tab)
}

// Download скачивание файла
func (s *service) Download(c *gin.Context) {
	filename := c.Query("name")
	f, err := os.Open(filepath.Join(filename))
	if err != nil {
		log.Println("Can't send file", filename, err)
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "Can't send file",
			},
		)
		return
	}
	defer f.Close()
	fileStats, err := f.Stat()
	if err != nil {
		log.Println("Can't send file", filename, err)
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": "Can't get info about file",
			},
		)
		return
	}
	contentLength := fileStats.Size()
	contentType := "application/octet-stream"
	extraHeaders := map[string]string{
		"Content-Disposition": `attachment; filename="` + filename + `"`,
	}

	c.DataFromReader(http.StatusOK, contentLength, contentType, f, extraHeaders)
}
