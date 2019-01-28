package main

import (
	"database/sql"

	_ "github.com/lib/pq"
)

// Database - used to access the database
type Database struct {
	db *sql.DB
}

// NewDatabase - creates a new Database instance
func NewDatabase(addr string) (*Database, error) {
	db, err := sql.Open("postgres", addr)

	if err != nil {
		return nil, err
	}

	return &Database{
		db: db,
	}, nil
}
