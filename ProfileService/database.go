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
	if _, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS user_instruments (
		id SERIAL NOT NULL PRIMARY KEY,
		user_id INT NOT NULL,
		instrument_id INT NOT NULL
	)`); err != nil {
		return err
	}
	return nil
}

func (db *database) getUser(id int) (userInfo, error) {
	var user userInfo
	if err := db.QueryRow("select * from user_info where id = $1", id).Scan(&user.ID, &user.AccountID,
		&user.Name, &user.Registered, &user.Birthday); err != nil {
		return userInfo{}, fmt.Errorf("getting user by id: %v", err)
	}
	return user, nil
}

func (db *database) getInstruments(id int) ([]instrument, error) {
	instruments := make([]instrument, 0)
	rows, err := db.Query("select instrument_id from user_instruments where user_id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("getting instrumets ids: %v", err)
	}
	for rows.Next() {
		var instrumentID int32
		var inst instrument
		if err := rows.Scan(&instrumentID); err != nil {
			return nil, fmt.Errorf("scannin instrument id: %v", err)
		}
		if err := db.QueryRow("select id, name from instruments where id = $1", instrumentID).Scan(
			&inst.ID, &inst.Name); err != nil {
			return nil, fmt.Errorf("scannin instrument id: %v", err)
		}
		instruments = append(instruments, inst)
	}
	return instruments, nil
}

func (db *database) setNewUser(user *userInfo) error {
	if err := db.QueryRow("insert into user_info (account_id, name, registered, bitrhday) vales ($1, $2, now(), $3) returning id, registered",
		user.AccountID, user.Name, user.Birthday).Scan(&(user.ID), &(user.Registered)); err != nil {
		return fmt.Errorf("inserting new user: %v", err)
	}
	return nil
}

func (db *database) setNewUserInstruments(user userInfo) error {
	for _, inst := range user.Instruments {
		if _, err := db.Exec("insert into user_instruments (user_id, instrument_id) values ($1, $2)", user.ID, inst.ID); err != nil {
			return fmt.Errorf("inserting new user instrument: %v", err)
		}
	}
	return nil
}

func (db *database) getProfileByUser(userID int) (userInfo, error) {
	var user userInfo
	if err := db.QueryRow("select * from user_info where account_id = $1", userID).Scan(&user.ID, &user.AccountID,
		&user.Name, &user.Registered, &user.Birthday); err != nil {
		return userInfo{}, fmt.Errorf("getting user by id: %v", err)
	}
	return user, nil
}
