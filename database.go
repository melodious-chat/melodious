package main

import (
	"database/sql"

	_ "github.com/lib/pq"
)

// Database - used to access the database
type Database struct {
	mel *Melodious
	db  *sql.DB
}

// RegisterUser - adds a new user to the database
func (db *Database) RegisterUser(name string, passhash string) error {
	//sum := sha256.Sum256([]byte(passhash))
	//sumstr := string(sum[:32])

	return nil
}

// NewDatabase - creates a new Database instance
func NewDatabase(mel *Melodious, addr string) (*Database, error) {
	db, err := sql.Open("postgres", addr)

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS accounts (
			id serial PRIMARY KEY,
			username varchar(32) UNIQUE,
			passhash varchar(64)
		);`)
	if err != nil {
		return nil, err
	}

	//if err != nil {
	//	return nil, err
	//}

	return &Database{
		mel: mel,
		db:  db,
	}, nil
}
