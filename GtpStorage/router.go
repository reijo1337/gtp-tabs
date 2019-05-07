package main

import (
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
	r.GET("/musicians", s.GetAuthorsByName)
	r.GET("/tabs", s.FindTabsByName)
	r.GET("/category/:name", s.GetAuthorsByCategory)
	r.PUT("/file", s.Upload)
	r.GET("/file", s.Download)
	return r, nil
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
	searchString := c.Query("search")
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
	searchString := c.Query("search")
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
	filename := c.GetHeader("Filename")
	fileBody, err := c.GetRawData()
	if err != nil {
		log.Println("Can't get file", err)
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "Can't get file",
			},
		)
		return
	}
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0755)
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
	_, err = file.Write(fileBody)
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
	c.JSON(
		http.StatusOK,
		gin.H{
			"ok": "ok",
		},
	)
}

// Download скачивание файла
func (s *service) Download(c *gin.Context) {
	filename := c.GetHeader("Filename")
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
