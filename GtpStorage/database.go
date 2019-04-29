package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/lib/pq"
)

// Database стуктура для взаимодействия с БД
type Database struct {
	*sql.DB
}

// SetUpDatabase устанавливает соединение с бд и разворачивает схему, если ее нет
func SetUpDatabase() (*Database, error) {
	log.Println("DB: Connecting to", DatabaseName, "database")
	db, err := sql.Open("postgres", fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=disable",
		DatabaseUserName, DatabasePassword, DatabaseName, DatabaseHost))
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(100)

	log.Println("Creating schema")
	if err := createSchema(db); err != nil {
		return nil, err
	}

	ddb := &Database{DB: db}

	log.Println("DB: succesful setup")
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
			size REAL DEFAULT 0
		)
	`); err != nil {
		return err
	}

	return nil
}

func (db *Database) getMusiciansByLetter(searchString string) ([]*MusiciansWithCount, error) {
	log.Printf("DB: Getting musicians by search request: %s\n", searchString)
	lowerSearchString := strings.ToLower(searchString)
	rows, err := db.Query("SELECT id, name FROM musicians WHERE (lower(name) LIKE '" + lowerSearchString + "%')")
	if err != nil {
		return nil, err
	}
	return db.getMusiciansWithCount(rows)
}

func (db *Database) getMusiciansWithCount(rows *sql.Rows) ([]*MusiciansWithCount, error) {
	result := make([]*MusiciansWithCount, 0)
	for rows.Next() {
		var resMusician *MusiciansWithCount
		rows.Scan(&(resMusician.ID), &(resMusician.Name))
		result = append(result, resMusician)
	}

	var count int32
	for _, musician := range result {
		err := db.QueryRow("SELECT count(*) FROM books WHERE author = $1", musician.ID).Scan(&count)
		if err != nil {
			return nil, err
		}
		musician.Count = count
	}
	return result, nil
}

func (db *Database) getMusicians(searchString string) ([]*MusiciansWithCount, error) {
	log.Printf("DB: Getting musicians by search request: %s\n", searchString)
	lowerSearchString := strings.ToLower(searchString)
	rows, err := db.Query("SELECT id, name FROM musicians WHERE (lower(name) LIKE '%" + lowerSearchString + "%')")
	if err != nil {
		return nil, err
	}
	return db.getMusiciansWithCount(rows)
}

func (db *Database) getTabsByName(searchString string) ([]*TabWithSize, error) {
	log.Printf("DB: Getting tabs with size by search request: %s\n", searchString)
	ret := make([]*TabWithSize, 0)
	lowerSearchString := strings.ToLower(searchString)
	rows, err := db.Query("SELECT author, name, size FROM tabs WHERE (lower(name) LIKE '%$1%')", lowerSearchString)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var tabInfo *TabWithSize
		var musicianID int32
		rows.Scan(&musicianID, &(tabInfo.Name), &(tabInfo.Size))
		if err := db.QueryRow("SELECT name FROM musicians WHERE id=$1", musicianID).Scan(&(tabInfo.Musician)); err != nil {
			return nil, err
		}
		ret = append(ret, tabInfo)
	}
	return ret, nil
}

func (db *Database) getMusiciansByCategory(category string) ([]*MusiciansWithCount, error) {
	log.Println("DB: Getting musicians by category search:", category)
	ret := make([]*MusiciansWithCount, 0)
	var categoryID int32
	err := db.QueryRow("SELECT id FROM categories WHERE (lower(name) = %1)", strings.ToLower(category)).Scan(&categoryID)
	if err != nil {
		return nil, err
	}
	rows, err := db.Query("SELECT count(*) as c, author FROM tabs WHERE category = %1 GROUP BY author", categoryID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var tabInfo *MusiciansWithCount
		rows.Scan(&(tabInfo.Count), &(tabInfo.ID))
		if err := db.QueryRow("SELECT name FROM musicians WHERE id=$1", tabInfo.ID).Scan(&(tabInfo.Name)); err != nil {
			return nil, err
		}
		ret = append(ret, tabInfo)
	}
	return ret, nil
}
