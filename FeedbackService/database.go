package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type database struct {
	*sql.DB
}

func setUpDatabase(source string) (*database, error) {
	db, err := sql.Open("postgres", source)
	if err != nil {
		return nil, fmt.Errorf("opening connection: %v", err)
	}
	if err := createSchema(db); err != nil {
		return nil, err
	}
	return &database{DB: db}, nil
}

func createSchema(db *sql.DB) error {
	if _, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS feedback (
		id SERIAL NOT NULL PRIMARY KEY,
		username VARCHAR(50) NOT NULL,
		user_email VARCHAR(50) NOT NULL,
		massage text NOT NULL
	)`); err != nil {
		return err
	}
	return nil
}

func (d *database) getFeedbacks(limit int, page int) ([]feedback, error) {
	ret := make([]feedback, 0)
	rows, err := d.Query("select * from feedback limit $1 offset $2", limit, page)
	if err != nil {
		return nil, fmt.Errorf("can't get feedbacks: %v", err)
	}
	var fb feedback
	for rows.Next() {
		if err := rows.Scan(&fb.ID, &fb.Username, &fb.Email, &fb.Message); err != nil {
			return nil, fmt.Errorf("can't scan feedback: %v", err)
		}
		ret = append(ret, fb)
	}
	return ret, nil
}

func (d *database) addFeedback(fb feedback) error {
	if err := d.QueryRow("insert into feedback (username, user_email, message) values ($1, $2, $3)",
		fb.Username, fb.Email, fb.Message); err != nil {
		return fmt.Errorf("can't insert feedback: %v", err)
	}
	return nil
}

func (d *database) removeFeedback(feedbackID int) error {
	if _, err := d.Exec("delete from feedback where id = $1", feedbackID); err != nil {
		return fmt.Errorf("can't delete feedback: %v", err)
	}
	return nil
}
