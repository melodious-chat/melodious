package main

import (
	"crypto/sha256"
	"database/sql"
	"errors"

	_ "github.com/lib/pq"
)

// Database - used to access the database
type Database struct {
	mel *Melodious
	db  *sql.DB
}

// RegisterUser - adds a new user to the database
func (db *Database) RegisterUser(name string, passhash string) error {
	return db.RegisterUserOwner(name, passhash, false)
}

// RegisterUserOwner - adds a new user to the database, possibly owner
func (db *Database) RegisterUserOwner(name string, passhash string, owner bool) error {
	sum := sha256.Sum256([]byte(passhash))
	sumstr := string(sum[:32])

	_, err := db.db.Exec(`
		INSERT INTO accounts (username, passhash, owner) VALUES (?, ?, ?);
	`, name, sumstr, owner)

	if err != nil {
		return err
	}

	return nil
}

// GetUserID - gets id of user with the given username.
// Always check if returned error is not nil, as it returns -1 on errors
func (db *Database) GetUserID(name string) (int, error) {
	row := db.db.QueryRow(`
		SELECT id FROM accounts WHERE username=?;
	`, name)

	var id int
	err := row.Scan(&id)
	if err == sql.ErrNoRows {
		return -1, errors.New("no such user")
	} else if err != nil {
		return -1, err
	}

	return 0, err
}

// UserExists - checks if user with given name exists.
// Always check if returned error is not-nil, as it returns false on errors
func (db *Database) UserExists(name string) (bool, error) {
	row := db.db.QueryRow(`
		SELECT id FROM accounts WHERE username=?;
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

// UserExistsID - checks if user with given id exists.
// Always check if returned error is not-nil, as it returns false on errors
func (db *Database) UserExistsID(id int) (bool, error) {
	row := db.db.QueryRow(`
		SELECT id FROM accounts WHERE id=?;
	`, id)

	var _id int // this is unused though
	err := row.Scan(&_id)
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
		SELECT id FROM accounts WHERE username=? AND passhash=?;
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

// CheckUserPasswordID - checks if there's a user with the given password.
// Always check if returned error is not-nil, as it returns false on errors
func (db *Database) CheckUserPasswordID(id int, passhash string) (bool, error) {
	sum := sha256.Sum256([]byte(passhash))
	sumstr := string(sum[:32])

	row := db.db.QueryRow(`
		SELECT id FROM accounts WHERE id=? AND passhash=?;
	`, id, sumstr)

	var _id int // this is unused though
	err := row.Scan(&_id)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

// IsUserOwner - checks if user with given name is an owner
func (db *Database) IsUserOwner(name string) (bool, error) {
	row := db.db.QueryRow(`
		SELECT id FROM accounts WHERE username=? AND owner=TRUE;
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

// IsUserOwnerID - checks if user with given id is an owner
func (db *Database) IsUserOwnerID(id int) (bool, error) {
	row := db.db.QueryRow(`
		SELECT id FROM accounts WHERE id=? AND owner=TRUE;
	`, id)

	var _id int // this is unused though
	err := row.Scan(&_id)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

// SetUserOwner - sets users owner status
func (db *Database) SetUserOwner(name string, owner bool) error {
	_, err := db.db.Exec(`
		UPDATE accounts SET owner=? WHERE username=?;
	`, owner, name)

	if err != nil {
		return err
	}

	return nil
}

// SetUserOwnerID - sets users owner status
func (db *Database) SetUserOwnerID(id int, owner bool) error {
	_, err := db.db.Exec(`
		UPDATE accounts SET owner=? WHERE id=?;
	`, owner, id)

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
			owner BOOLEAN
		);`)
	if err != nil {
		return nil, err
	}

	return &Database{
		mel: mel,
		db:  db,
	}, nil
}
