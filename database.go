package main

import (
	"crypto/sha256"
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
	return db.RegisterUserAdmin(name, passhash, false)
}

// RegisterUserAdmin - adds a new user to the database, possibly admin
func (db *Database) RegisterUserAdmin(name string, passhash string, admin bool) error {
	sum := sha256.Sum256([]byte(passhash))
	sumstr := string(sum[:32])

	_, err := db.db.Exec(`
		INSERT INTO accounts (username, passhash, admin) VALUES (?, ?, ?);
	`, name, sumstr, admin)

	if err != nil {
		return err
	}

	return nil
}

// UserExists - checks if user with given name exists.
// Always check if returned error is not-nil, as it returns false on errors
func (db *Database) UserExists(name string) (bool, error) {
	row := db.db.QueryRow(`
		SELECT id FROM accounts WHERE EXISTS(SELECT id FROM accounts WHERE username=?);
	`, name)

	var id int // this is unused though
	err := row.Scan(&id)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

// CheckUserPassword - checks if there's a user with the given password.
// Always check if returned error is not-nil, as it returns false on errors
func (db *Database) CheckUserPassword(name string, passhash string) (bool, error) {
	sum := sha256.Sum256([]byte(passhash))
	sumstr := string(sum[:32])

	row := db.db.QueryRow(`
		SELECT id FROM accounts WHERE EXISTS(SELECT id FROM accounts WHERE username=? AND passhash=?);
	`, name, sumstr)

	var id int // this is unused though
	err := row.Scan(&id)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

// IsUserAdmin - checks if user with given name is an admin
func (db *Database) IsUserAdmin(name string) (bool, error) {
	row := db.db.QueryRow(`
		SELECT id FROM accounts WHERE EXISTS(SELECT id FROM accounts WHERE username=? AND admin=TRUE);
	`, name)

	var id int // this is unused though
	err := row.Scan(&id)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

// SetUserAdmin - sets users admin status
func (db *Database) SetUserAdmin(name string, admin bool) error {
	_, err := db.db.Exec(`
		UPDATE accounts SET admin=? WHERE username=?;
	`, admin, name)

	if err != nil {
		return err
	}

	return nil
}

// NewDatabase - creates a new Database instance
func NewDatabase(mel *Melodious, addr string) (*Database, error) {
	db, err := sql.Open("postgres", addr)

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS accounts (
			id serial PRIMARY KEY,
			username varchar(32) UNIQUE,
			passhash varchar(64),
			admin BOOLEAN
		);`)
	if err != nil {
		return nil, err
	}

	return &Database{
		mel: mel,
		db:  db,
	}, nil
}
