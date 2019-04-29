package main

import (
	"database/sql"
	"fmt"
	"log"
	"crypto/sha256"

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
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL NOT NULL PRIMARY KEY,
			login VARCHAR(20) UNIQUE NOT NULL,
			passhash VARCHAR(70) NOT NULL
		)`); err != nil {
		return err
	}

	return nil
}

func (db *Database) insertNewUser(login string, password string) (*user, error) {
	log.Println("DB: Inserting new user ", login)
	passHash := sha256.New()
	passHash.Write([]byte(password))
	pass := passHash.Sum(nil)
	passStr := fmt.Sprintf("%x\n", pass)

	rows, err := db.Query("SELECT id FROM users WHERE login = $1", login)

	if err != nil {
		return nil, err
	}

	var ID int32
	for rows.Next() {
		rows.Scan(&ID)
	}

	if ID > 0 {
		return &user{ID: ID, Login: login, Password: password}, nil
	}

	row := db.QueryRow("INSERT INTO users (login, passhash) VALUES ($1, $2) RETURNING id",
		login, passStr)
	if err := row.Scan(&ID); err != nil {
		return nil, err
	}
	log.Println("DB: user inserted succesfully")
	return &user{
		ID:       ID,
		Login:    login,
		Password: password,
	}, nil
}

func (db *Database) isAuthorized(user *user) bool {
	log.Println("DB: Check user is in DB")
	var (
		passhash string
	)

	err := db.QueryRow("SELECT passhash FROM users WHERE login = $1", user.Login).Scan(
		&passhash)
	if err != nil {
		return false
	}
	passHash := sha256.New()
	passHash.Write([]byte(user.Password))
	pass := passHash.Sum(nil)
	passStr := fmt.Sprintf("%x\n", pass)
	return passhash == passStr
}