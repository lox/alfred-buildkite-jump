package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var respositorySql = `
CREATE TABLE IF NOT EXISTS pipeline (
    id varchar(255) PRIMARY KEY,
    url varchar(255) NOT NULL,
    org varchar(255) NOT NULL,
    name varchar(255) NOT NULL,
    description text,
    created_at timestamp,
    updated_at timestamp
)`

func OpenDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./sqlite.db")
	if err != nil {
		return nil, err
	}

	if _, err = db.Exec(respositorySql); err != nil {
		return nil, err
	}

	return db, nil
}
