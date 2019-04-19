package main

import (
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
)

func TestFindMusiciansByLetter(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	ddb := &Database{DB: db}

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE products").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO product_viewers").WithArgs(2, 3).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	server, err := MakeServer(ddb)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when creating grpc server", err)
	}
}
