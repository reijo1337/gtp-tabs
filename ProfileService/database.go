package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type database struct {
	*sql.DB
}

func setUpDatabase() (*database, error) {
	config, err := parseConfig("PROFILE")
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

	ddb := &database{DB: db}

	return ddb, nil
}

func createSchema(db *sql.DB) error {
	if _, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS instruments (
		id SERIAL NOT NULL PRIMARY KEY,
		name VARCHAR(50) NOT NULL UNIQUE,
		pic_path VARCHAR(50) NOT NULL
	)`); err != nil {
		return err
	}

	if _, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS user_info (
		id SERIAL NOT NULL PRIMARY KEY,
		account_id INT NOT NULL,
		name VARCHAR(50) NOT NULL UNIQUE,
		profile_pic_path VARCHAR(50) NOT NULL,
		registered DATE,
		birthday DATE
	)`); err != nil {
		return err
	}
	if _, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS user_tabs (
		id SERIAL NOT NULL PRIMARY KEY,
		user_id INT NOT NULL,
		tab_id INT NOT NULL
	)`); err != nil {
		return err
	}
	return nil
}
