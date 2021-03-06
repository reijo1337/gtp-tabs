package main

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// Database стуктура для взаимодействия с БД
type Database struct {
	*sql.DB
}

// SetUpDatabase устанавливает соединение с бд и разворачивает схему, если ее нет
func SetUpDatabase() (*Database, error) {
	config, err := parseConfig("AUTH")
	if err != nil {
		return nil, err
	}
	log.Println("DB: Connecting to", config.DB.Name, "database")
	db, err := sql.Open("postgres", fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=disable",
		config.DB.User, config.DB.Password, config.DB.Name, config.DB.Host))
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(100)

	log.Println("Creating schema")
	if err := createSchema(db); err != nil {
		return nil, err
	}

	if err := populateDB(db); err != nil {
		return nil, fmt.Errorf("populating db: %v", err)
	}

	ddb := &Database{DB: db}

	log.Println("DB: succesful setup")
	return ddb, nil
}

func createSchema(db *sql.DB) error {
	if _, err := db.Exec(`
		CREATE SEQUENCE IF NOT EXISTS all_users_id_seq;
		`); err != nil {
		return err
	}
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS roles (
			id SERIAL NOT NULL PRIMARY KEY,
			name VARCHAR(20) UNIQUE NOT NULL
		)`); err != nil {
		return err
	}
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INT DEFAULT nextval('all_users_id_seq') NOT NULL,
			login VARCHAR(20) UNIQUE NOT NULL,
			passhash VARCHAR(70) NOT NULL,
			role INT NOT NULL
		)`); err != nil {
		return err
	}
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS vk_users (
			id INT DEFAULT nextval('all_users_id_seq') NOT NULL,
			vk_id BIGINT NOT NULL,
			role INT NOT NULL
		)`); err != nil {
		return err
	}

	return nil
}

func populateDB(db *sql.DB) error {
	var roleID int32
	err := db.QueryRow("SELECT id FROM roles WHERE name = 'user'").Scan(&roleID)
	if err != nil {
		if _, err := db.Exec("insert into roles (name) values ('user')"); err != nil {
			return fmt.Errorf("can't insert role user: %v", err)
		}
	}
	err = db.QueryRow("SELECT id FROM roles WHERE name = 'amdin'").Scan(&roleID)
	if err != nil {
		if _, err := db.Exec("insert into roles (name) values ('amdin')"); err != nil {
			return fmt.Errorf("can't insert role amdin: %v", err)
		}
	}
	return nil
}

func (db *Database) insertNewUser(login string, password string, role string) (*user, error) {
	log.Println("DB: Inserting new user", login, "with role", role)
	var roleID int32
	err := db.QueryRow("SELECT id FROM roles WHERE name = $1", role).Scan(&roleID)
	if err != nil {
		return nil, err
	}
	passHash := sha256.New()
	passHash.Write([]byte(password))
	pass := passHash.Sum(nil)
	passStr := fmt.Sprintf("%x\n", pass)

	var ID int32
	if err = db.QueryRow("SELECT id FROM users WHERE login = $1 and role = $2", login, roleID).Scan(&ID); err == nil {
		return &user{ID: ID, Login: login, Password: password}, nil
	}

	row := db.QueryRow("INSERT INTO users (login, passhash, role) VALUES ($1, $2, $3) RETURNING id",
		login, passStr, roleID)
	if err := row.Scan(&ID); err != nil {
		return nil, err
	}
	log.Println("DB: user inserted succesfully")
	return &user{
		ID:       ID,
		Login:    login,
		Password: password,
		Role: Role{
			ID:   roleID,
			Name: role,
		},
	}, nil
}

func (db *Database) isAuthorized(user *user) bool {
	log.Println("DB: Check user is in DB")
	var (
		passhash string
	)

	err := db.QueryRow("SELECT passhash, id FROM users WHERE login = $1", user.Login).Scan(
		&passhash, &(user.ID))
	if err != nil {
		return false
	}
	passHash := sha256.New()
	passHash.Write([]byte(user.Password))
	pass := passHash.Sum(nil)
	passStr := fmt.Sprintf("%x\n", pass)
	return passhash == passStr
}

func (db *Database) isAuthorizedVK(user *vkUser) bool {
	log.Println("DB: Check user is in DB")
	err := db.QueryRow("SELECT id FROM vk_users WHERE vk_id = $1", user.UserID).Scan(
		&(user.ID))
	if err != nil {
		return false
	}
	return true
}

func (db *Database) credentialsInUse(user *user) bool {
	err := db.QueryRow("SELECT id FROM users WHERE login = $1", user.Login).Scan(&(user.ID))
	if err != nil {
		return false
	}
	return true
}

func (db *Database) insertNewVkUser(userID int64, role string) (*vkUser, error) {
	var roleID int32
	var ID int32
	err := db.QueryRow("SELECT id FROM roles WHERE name = $1", role).Scan(&roleID)
	if err != nil {
		return nil, err
	}
	if err := db.QueryRow("INSERT INTO vk_users (vk_id, role) VALUES ($1, $2) RETURNING id", userID, roleID).Scan(&ID); err != nil {
		return nil, err
	}
	return &vkUser{
		ID:     ID,
		UserID: userID,
		Role: Role{
			ID:   roleID,
			Name: role,
		},
	}, nil
}
