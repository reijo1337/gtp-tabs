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
	CREATE TABLE IF NOT EXISTS post (
		id SERIAL NOT NULL PRIMARY KEY,
		song_name VARCHAR(50) NOT NULL,
		tab_id INT NOT NULL,
		author_id INT NOT NULL,
		rating numeric DEFAULT 0.00,
		rating_count int DEFAULT 0
	)`); err != nil {
		return err
	}
	if _, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS comments (
		id SERIAL NOT NULL PRIMARY KEY,
		author_id INT NOT NULL,
		contents TEXT NOT NULL
	)`); err != nil {
		return err
	}
	if _, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS post_comments (
		id SERIAL NOT NULL PRIMARY KEY,
		post_id INT NOT NULL,
		comment_id INT NOT NULL
	)`); err != nil {
		return err
	}
	return nil
}

func (d *database) getPost(tabID int) (post, error) {
	var postInfo post
	var ratingCount int
	if err := d.QueryRow("select * from post where tab_id = $1", tabID).Scan(
		&postInfo.ID, &postInfo.SongName, &postInfo.TabID, &postInfo.AuthorID, &postInfo.Rating, &ratingCount); err != nil {
		return post{}, fmt.Errorf("getting post by tab id: %v", err)
	}
	coms := make([]comment, 0)
	rows, err := d.Query("select comment_id from post_comments where post_id = $1", postInfo.ID)
	if err != nil {
		return post{}, fmt.Errorf("getting post comments id: %v", err)
	}
	var commentID int
	for rows.Next() {
		var com comment
		if err := rows.Scan(&commentID); err != nil {
			return post{}, fmt.Errorf("scanning comment id: %v", err)
		}
		if err := d.QueryRow("select * from comments where id = $1", commentID).Scan(
			&com.ID, &com.AuthorID, &com.Content); err != nil {
			return post{}, fmt.Errorf("getting comment: %v", err)
		}
		coms = append(coms, com)
	}
	postInfo.Comments = coms
	return postInfo, nil
}

func (d *database) makePost(postInfo *post) error {
	if postInfo == nil {
		return fmt.Errorf("nil input")
	}
	if err := d.QueryRow("insert into post (song_name, tab_id, author_id) values ($1, $2, $3) returning id, rating",
		postInfo.SongName, postInfo.TabID, postInfo.AuthorID).Scan(&(postInfo.ID), &(postInfo.Rating)); err != nil {
		return fmt.Errorf("inserting post info: %v", err)
	}
	return nil
}

func (d *database) makeComment(postID int, com *comment) error {
	if com == nil {
		return fmt.Errorf("nil imput")
	}
	if err := d.QueryRow("insert into comments (author_id, contents) values ($1, $2) returning id",
		com.AuthorID, com.Content).Scan(&(com.ID)); err != nil {
		return fmt.Errorf("inserting comment: %v", err)
	}
	if _, err := d.Exec("insert into post_comments (post_id, comment_id) values ($1, $2)", postID, com.ID); err != nil {
		return fmt.Errorf("inserting post comment: %v", err)
	}
	return nil
}

func (d database) changeRating(postID int, userRating int) error {
	var ID, ratingCount int
	var rating float32
	if err := d.QueryRow("select id, rating, rating_count from post where id = $1", postID).Scan(
		&ID, &rating, &ratingCount); err != nil {
		return fmt.Errorf("getting post by id: %v", err)
	}
	newCount := ratingCount + 1
	newRating := (rating + float32(userRating)) / float32(newCount)
	if _, err := d.Exec("update post set rating = $1, rating_count = $2 where id = $3", newRating, newCount, postID); err != nil {
		return fmt.Errorf("updating rating: %v", err)
	}
	return nil
}
