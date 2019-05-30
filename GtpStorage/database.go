package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

// Database стуктура для взаимодействия с БД
type Database struct {
	*sql.DB
}

// SetUpDatabase устанавливает соединение с бд и разворачивает схему, если ее нет
func SetUpDatabase() (*Database, error) {
	config, err := parseConfig("STORAGE")
	if err != nil {
		return nil, err
	}
	db, err := sql.Open("postgres", fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=disable",
		config.DB.User, config.DB.Password, config.DB.Name, config.DB.Host))
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(100)

	if err := createSchema(db); err != nil {
		return nil, err
	}

	ddb := &Database{DB: db}

	return ddb, nil
}

func createSchema(db *sql.DB) error {
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS musicians (
			id SERIAL NOT NULL PRIMARY KEY,
			name VARCHAR(50) NOT NULL UNIQUE
		)`); err != nil {
		return err
	}

	if _, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS categories (
		id SERIAL NOT NULL PRIMARY KEY,
		name VARCHAR(50) NOT NULL UNIQUE
	)
	`); err != nil {
		return err
	}

	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS tabs (
			id SERIAL NOT NULL PRIMARY KEY,
			author INT NOT NULL REFERENCES musicians (id),
			name VARCHAR(50) NOT NULL UNIQUE,
			category INT NOT NULL REFERENCES categories (id),
			size BIGINT DEFAULT 0
		)
	`); err != nil {
		return err
	}

	return nil
}

func (db *Database) getMusiciansByLetter(searchString string) ([]MusiciansWithCount, error) {
	lowerSearchString := strings.ToLower(searchString)
	rows, err := db.Query("SELECT id, name FROM musicians WHERE (lower(name) LIKE '" + lowerSearchString + "%')")
	if err != nil {
		return nil, err
	}
	return db.getMusiciansWithCount(rows)
}

func (db *Database) getMusiciansByNumber() ([]MusiciansWithCount, error) {
	rows, err := db.Query("SELECT id, name FROM musicians WHERE (lower(name) LIKE '[0-9]%')")
	if err != nil {
		return nil, err
	}
	return db.getMusiciansWithCount(rows)
}

func (db *Database) getMusiciansWithCount(rows *sql.Rows) ([]MusiciansWithCount, error) {
	result := make([]MusiciansWithCount, 0)
	for rows.Next() {
		var resMusician MusiciansWithCount
		err := rows.Scan(&resMusician.ID, &resMusician.Name)
		if err != nil {
			return nil, err
		}
		result = append(result, resMusician)
	}

	var count int32
	for _, musician := range result {
		err := db.QueryRow("SELECT count(*) FROM tabs WHERE author = $1", musician.ID).Scan(&count)
		if err != nil {
			return nil, err
		}
		musician.Count = count
	}
	return result, nil
}

func (db *Database) getMusicians(searchString string) ([]MusiciansWithCount, error) {
	lowerSearchString := strings.ToLower(searchString)
	rows, err := db.Query("SELECT id, name FROM musicians WHERE (lower(name) LIKE '%" + lowerSearchString + "%')")
	if err != nil {
		return nil, err
	}
	return db.getMusiciansWithCount(rows)
}

func (db *Database) getTabsByName(searchString string) ([]TabWithSize, error) {
	ret := make([]TabWithSize, 0)
	lowerSearchString := strings.ToLower(searchString)
	rows, err := db.Query("SELECT author, name, size FROM tabs WHERE (lower(name) LIKE '%" + lowerSearchString + "%')")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var tabInfo TabWithSize
		var musicianID int32
		_ = rows.Scan(&musicianID, &tabInfo.Name, &tabInfo.Size)
		if err := db.QueryRow("SELECT name FROM musicians WHERE id=$1", musicianID).Scan(&(tabInfo.Musician)); err != nil {
			return nil, err
		}
		ret = append(ret, tabInfo)
	}
	return ret, nil
}

func (db *Database) getMusiciansByCategory(category string) ([]MusiciansWithCount, error) {
	ret := make([]MusiciansWithCount, 0)
	var categoryID int32
	err := db.QueryRow("SELECT id FROM categories WHERE (lower(name) = $1)", strings.ToLower(category)).Scan(&categoryID)
	if err != nil {
		if err == sql.ErrNoRows {
			return ret, nil
		}
		return nil, err
	}
	rows, err := db.Query("SELECT count(*) as c, author FROM tabs WHERE category = $1 GROUP BY author", categoryID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var tabInfo MusiciansWithCount
		_ = rows.Scan(&tabInfo.Count, &tabInfo.ID)
		if err := db.QueryRow("SELECT name FROM musicians WHERE id=$1", tabInfo.ID).Scan(&(tabInfo.Name)); err != nil {
			return nil, err
		}
		ret = append(ret, tabInfo)
	}
	return ret, nil
}

func (db *Database) getOrCreateMusician(name string) (musician, error) {
	ret := musician{
		Name: name,
	}
	err := db.QueryRow("SELECT id FROM musicians WHERE (lower(name) = $1)", strings.ToLower(name)).Scan(&ret.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			err = db.QueryRow("INSERT INTO musicians (name) VALUES ($1) RETURNING id", name).Scan(&ret.ID)
			if err != nil {
				return musician{}, err
			}
			return ret, nil
		}
		return musician{}, err
	}
	return ret, nil
}

func (db *Database) getOrCreateCategory(name string) (category, error) {
	ret := category{
		Name: name,
	}
	err := db.QueryRow("SELECT id FROM categories WHERE (lower(name) = $1)", strings.ToLower(name)).Scan(&ret.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			err = db.QueryRow("INSERT INTO categories (name) VALUES ($1) RETURNING id", name).Scan(&ret.ID)
			if err != nil {
				return category{}, err
			}
			return ret, nil
		}
		return category{}, err
	}
	return ret, nil
}

func (db *Database) createSong(musicianID, categoryID int32, name string, size int64) (int, error) {
	var ID string
	err := db.QueryRow("SELECT name FROM tabs WHERE author = $1 AND category = $2 AND name LIKE '"+name+"%' ORDER BY name DESC LIMIT 1", musicianID, categoryID).Scan(&ID)
	if err != nil {
		if err != sql.ErrNoRows {
			return 0, err
		}
	} else {
		if ID[len(ID)-1] == ')' {
			nameN, err := strconv.Atoi(ID[len(ID)-2 : len(ID)-1])
			if err != nil {
				return 0, err
			}
			newNameN := strconv.Itoa(nameN + 1)
			name = ID[:len(ID)-2] + newNameN + ")"
		} else {
			name = name + " (1)"
		}
	}
	var tabID int
	err = db.QueryRow("INSERT INTO tabs (author, name, category, size) VALUES ($1, $2, $3, $4) returning id",
		musicianID, name, categoryID, size).Scan(&tabID)
	return tabID, err
}
