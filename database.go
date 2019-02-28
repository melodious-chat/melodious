package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/apex/log"
	"github.com/lib/pq"
)

// Database - used to access the database
type Database struct {
	mel *Melodious
	db  *sql.DB
}

// GetUsersList - gets users' data stored in the database
func (db *Database) GetUsersList() ([]*User, error) {
	rows, err := db.db.Query(`
		SELECT id, username, owner FROM melodious.accounts WHERE banned=false;
	`)
	if err != nil {
		return []*User{}, err
	}
	users := []*User{}
	for rows.Next() {
		user := &User{}
		err := rows.Scan(&(user.ID), &(user.Username), &(user.Owner))
		if err != nil {
			return []*User{}, err
		}
		users = append(users, user)
	}
	return users, nil
}

// HasUsers - checks if there are any users registered
func (db *Database) HasUsers() (bool, error) {
	row := db.db.QueryRow(`
		SELECT id FROM melodious.accounts LIMIT 1;
	`)

	var id int
	err := row.Scan(&id)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, err
}

// RegisterUser - adds a new user to the database
func (db *Database) RegisterUser(name string, passhash string, ip string) error {
	return db.RegisterUserOwner(name, passhash, false, ip)
}

// RegisterUserOwner - adds a new user to the database, possibly owner
func (db *Database) RegisterUserOwner(name string, passhash string, owner bool, ip string) error {
	sum := sha256.Sum256([]byte(passhash))
	sumstr := fmt.Sprintf("%x", sum[:32])

	_, err := db.db.Exec(`
		INSERT INTO melodious.accounts (username, passhash, owner) VALUES ($1, $2, $3);
	`, name, sumstr, owner)

	if err != nil {
		return err
	}

	return nil
}

// DeleteUser - deletes/unregisters a user
func (db *Database) DeleteUser(name string) error {
	_, err := db.db.Exec(`
		DELETE FROM melodious.accounts WHERE username=$1;
	`, name)
	if err != nil {
		return err
	}
	return nil
}

// DeleteUserID - deletes/unregisters a user by their id
func (db *Database) DeleteUserID(id int) error {
	_, err := db.db.Exec(`
		DELETE FROM melodious.accounts WHERE id=$1;
	`, id)
	if err != nil {
		return err
	}
	return nil
}

// GetUserID - gets id of user with the given username.
// Always check if returned error is not nil, as it returns -1 on errors
func (db *Database) GetUserID(name string) (int, error) {
	row := db.db.QueryRow(`
		SELECT id FROM melodious.accounts WHERE username=$1 LIMIT 1;
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
		SELECT id FROM melodious.accounts WHERE username=$1 LIMIT 1;
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
		SELECT id FROM melodious.accounts WHERE id=$1 LIMIT 1;
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
	sumstr := fmt.Sprintf("%x", sum[:32])

	row := db.db.QueryRow(`
		SELECT id FROM melodious.accounts WHERE username=$1 AND passhash=$2 LIMIT 1;
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
	sumstr := fmt.Sprintf("%x", sum[:32])

	row := db.db.QueryRow(`
		SELECT id FROM melodious.accounts WHERE id=$1 AND passhash=$2 LIMIT 1;
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
		SELECT id FROM melodious.accounts WHERE username=$1 AND owner=TRUE LIMIT 1;
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
		SELECT id FROM melodious.accounts WHERE id=$1 AND owner=TRUE LIMIT 1;
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
		UPDATE melodious.accounts SET owner=$1 WHERE username=$2;
	`, owner, name)

	if err != nil {
		return err
	}

	return nil
}

// SetUserOwnerID - sets users owner status
func (db *Database) SetUserOwnerID(id int, owner bool) error {
	_, err := db.db.Exec(`
		UPDATE melodious.accounts SET owner=$1 WHERE id=$2;
	`, owner, id)

	if err != nil {
		return err
	}

	return nil
}

// Ban - sets user's banned flag to true
func (db *Database) Ban(username string) error {
	_, err := db.db.Exec(`
		UPDATE melodious.accounts SET banned=true WHERE username=$1;
	`, username)
	if err != nil {
		return err
	}
	return nil
}

// BanID - sets user's banned flag to true by id
func (db *Database) BanID(id int) error {
	_, err := db.db.Exec(`
		UPDATE melodious.accounts SET banned=true WHERE id=$1;
	`, id)
	if err != nil {
		return err
	}
	return nil
}

// IsUserBanned - checks if the given user with ip is banned
func (db *Database) IsUserBanned(username string, ip string) (bool, error) {
	row := db.db.QueryRow(`
		SELECT banned FROM melodious.accounts WHERE username=$1 OR ip=$2;
	`, username, ip)
	var banned bool
	err := row.Scan(&banned)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return banned, nil
}

// IsUserBannedID - checks if the given user by id with ip is banned
func (db *Database) IsUserBannedID(id int, ip string) (bool, error) {
	row := db.db.QueryRow(`
		SELECT banned FROM melodious.accounts WHERE id=$1 OR ip=$2;
	`, id, ip)
	var banned bool
	err := row.Scan(&banned)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return banned, nil
}

// GetAccountCount - gets a count of all accounts that are on the same IP.
func (db *Database) GetAccountCount(ip string) (int, error) {
	row := db.db.QueryRow(`
		SELECT COUNT(*) FROM melodious.accounts WHERE ip=$1;
	`, ip)
	var count int
	err := row.Scan(&count)
	if err != nil {
		return -1, err
	}
	return count, nil
}

// NewChannel - creates a new channel
func (db *Database) NewChannel(name string, topic string) error {
	_, err := db.db.Exec(`
		INSERT INTO melodious.channels (name, topic) VALUES ($1, $2);
	`, name, topic)
	if err != nil {
		return err
	}
	return nil
}

// DeleteChannel - deletes given channel
func (db *Database) DeleteChannel(name string) error {
	_, err := db.db.Exec(`
		DELETE FROM melodious.channels WHERE name=$1;
	`, name)
	if err != nil {
		return err
	}
	return nil
}

// DeleteChannelID - deletes given channel by ID
func (db *Database) DeleteChannelID(id int) error {
	_, err := db.db.Exec(`
		DELETE FROM melodious.channels WHERE id=$1;
	`, id)
	if err != nil {
		return err
	}
	return nil
}

// ListChannels - puts all channel names into an array of Channel structs
func (db *Database) ListChannels() ([]*Channel, error) {
	rows, err := db.db.Query(`
		SELECT id, name, topic FROM melodious.channels;
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	m := []*Channel{}

	for rows.Next() {
		chnl := &Channel{}
		if err := rows.Scan(&(chnl.ID), &(chnl.Name), &(chnl.Topic)); err != nil {
			return nil, err
		}
		m = append(m, chnl)
	}

	return m, nil
}

// ChannelExists - checks if a channel exists
func (db *Database) ChannelExists(name string) (bool, error) {
	row := db.db.QueryRow(`
		SELECT EXISTS(SELECT * FROM melodious.channels WHERE name=$1);
	`, name)
	var exists bool
	err := row.Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// SetChannelTopic - sets topic of the channel
func (db *Database) SetChannelTopic(name string, topic string) error {
	_, err := db.db.Exec(`
		UPDATE melodious.channels SET topic=$2 WHERE name=$1;
	`, name, topic)
	if err != nil {
		return err
	}
	return nil
}

// DeleteOldMessages - deletes old messages from the database
func (db *Database) DeleteOldMessages(period string) error {
	_, err := db.db.Exec(`
		DELETE FROM melodious.messages WHERE dt < (NOW() - $1::INTERVAL);
	`, period)
	if err != nil {
		return err
	}
	return nil
}

// PostMessage - posts a new message
func (db *Database) PostMessage(chanName string, message string, pings []string, author string) (*ChatMessage, error) {
	row := db.db.QueryRow(`
		INSERT INTO melodious.messages
		(chan_id, message, dt, pings, author_id)
		VALUES (
			(SELECT id FROM melodious.channels WHERE name=$1 LIMIT 1),
			$2,
			NOW(),
			$3,
			(SELECT id FROM melodious.accounts WHERE username=$4 LIMIT 1)
		)
		RETURNING message, pings, id, dt, $4, author_id;
	`, chanName, message, pq.Array(pings), author)
	msg := &ChatMessage{}
	var cpings pq.StringArray
	err := row.Scan(&(msg.Message), &cpings, &(msg.ID), &(msg.Timestamp), &(msg.Author), &(msg.AuthorID))
	if err != nil {
		return nil, err
	}
	msg.Pings = []string(cpings)
	return msg, nil
}

// PostMessageChanID - posts a new message
func (db *Database) PostMessageChanID(chanID int, message string, pings []string, author string) error {
	_, err := db.db.Exec(`
		INSERT INTO melodious.messages
		(chan_id, message, dt, pings, author_id)
		VALUES (
			(SELECT id FROM melodious.channels WHERE id=$1 LIMIT 1),
			$2,
			NOW(),
			$3,
			(SELECT id FROM melodious.accounts WHERE username=$4 LIMIT 1)
		);
	`, chanID, message, pq.Array(pings), author)
	if err != nil {
		return err
	}
	return nil
}

// DeleteMessage - deletes a message by ID.
func (db *Database) DeleteMessage(id int) error {
	_, err := db.db.Exec(`
		DELETE FROM melodious.messages WHERE id=$1;
	`, id)
	if err != nil {
		return err
	}
	return nil
}

// EditMessage - edits a message by ID.
//func (db *Database) EditMessage(id int, content string) {
//todo
//}

// GetMessages - gets last n messages in a channel starting from an id from the database
func (db *Database) GetMessages(chanid int, msgid int, amount int) ([]*ChatMessage, error) {
	rows, err := db.db.Query(`
		SELECT
			m.id,
			m.message,
			m.dt,
			m.pings,
			a.username author,
			m.author_id
		FROM melodious.messages m
		INNER JOIN melodious.accounts a ON m.author_id = a.id
		WHERE m.chan_id=$1 AND m.id<$2
		ORDER BY m.id DESC
		LIMIT $3;
	`, chanid, msgid, amount)
	if err != nil {
		return []*ChatMessage{}, err
	}
	msgs := []*ChatMessage{}
	for rows.Next() {
		msg := &ChatMessage{}
		var pings pq.StringArray
		err := rows.Scan(&(msg.ID), &(msg.Message), &(msg.Timestamp), &pings, &(msg.Author), &(msg.AuthorID))
		if err != nil {
			return []*ChatMessage{}, err
		}
		msg.Pings = []string(pings)
		msgs = append(msgs, msg)
	}
	return msgs, nil
}

// GetMessageDetails - gets a message and the channel it's from by id
func (db *Database) GetMessageDetails(id int) (string, *ChatMessage, error) {
	row := db.db.QueryRow(`
		SELECT
			m.message,
			m.dt,
			m.pings,
			a.username author,
			m.author_id,
			c.name channel
		FROM melodious.messages m
		INNER JOIN melodious.accounts a ON m.author_id = a.id
		INNER JOIN melodious.channels c ON m.chan_id = c.id
		WHERE m.id=$1;
	`, id)
	var pings pq.StringArray
	var channel string
	msg := &ChatMessage{}
	err := row.Scan(&(msg.Message), &(msg.Timestamp), &pings, &(msg.Author), &(msg.AuthorID), &channel)
	if err != nil {
		return "", &ChatMessage{}, err
	}
	msg.Pings = []string(pings)
	msg.ID = id
	return channel, msg, nil
}

// AddGroup - adds a group
func (db *Database) AddGroup(name string) (int, error) {
	row := db.db.QueryRow(`
		INSERT INTO melodious.groups (name)	VALUES ($1) RETURNING id;
	`, name)
	var id int
	err := row.Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}

// DeleteGroup - deletes a group
func (db *Database) DeleteGroup(name string) error {
	_, err := db.db.Exec(`
		DELETE FROM melodious.groups WHERE name=$1;
	`, name)
	return err
}

// GroupExists - checks if a group exists
func (db *Database) GroupExists(name string) (bool, error) {
	row := db.db.QueryRow(`
		SELECT EXISTS(SELECT * FROM melodious.groups WHERE name=$1);
	`, name)
	var exists bool
	err := row.Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// SetFlag - sets a flag. Returns a flag id
func (db *Database) SetFlag(flag *Flag) (int, error) {
	data, err := json.Marshal(flag.Flag)
	if err != nil {
		return -1, err
	}

	row := db.db.QueryRow(`
		SELECT melodious.set_flag($1, $2, $3::JSON);
	`, flag.Group, flag.Name, string(data))
	var id int
	err = row.Scan(&id)
	if err != nil {
		return -1, err
	}

	return id, nil
}

// DeleteFlag - deletes a flag.
func (db *Database) DeleteFlag(flag *Flag) error {
	_, err := db.db.Exec(`
		CALL melodious.delete_flag($1, $2);
	`, flag.Group, flag.Name)
	return err
}

// AddGroupHolder - adds a group holder
func (db *Database) AddGroupHolder(gh *GroupHolder) (int, error) {
	// make sure if strings are empty we pass NULL to PostgreSQL
	var user interface{}
	var channel interface{}
	if gh.User != "" {
		user = gh.User
	} else {
		user = nil
	}
	if gh.Channel != "" {
		channel = gh.Channel
	} else {
		channel = nil
	}
	row := db.db.QueryRow(`
		SELECT melodious.insert_group_holder($1, $2, $3);
	`, gh.Group, user, channel)
	var id int
	err := row.Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}

// DeleteGroupHolder - deletes a group holder
func (db *Database) DeleteGroupHolder(id int) error {
	_, err := db.db.Exec(`
		DELETE FROM melodious.group_holders WHERE id=$1;
	`, id)
	return err
}

// DeleteGroupHolders - deletes group holders by template. Use empty strings where you usually would use a *
func (db *Database) DeleteGroupHolders(gh GroupHolder) error {
	_, err := db.db.Exec(`
		CALL melodious.delete_group_holders($1, $2, $3);
	`, gh.Group, gh.User, gh.Channel)
	return err
}

// QueryFlags -
// When checkflags:
//   true - queries if theres a flag available to the given channel/user/user-on-channel
//   false - queries flags using given pattern. Use empty strings where you usually would use a *
func (db *Database) QueryFlags(user string, channel string, group string, flag string, checkflags bool) ([]*FlagQueryResult, error) {
	rows, err := db.db.Query(`
		SELECT group_holders, flag_id, flag_name, flag FROM melodious.query_flags($1, $2, $3, $4, $5);
	`, user, channel, group, flag, checkflags)
	if err != nil {
		return nil, err
	}

	var s []*FlagQueryResult

	for rows.Next() {
		var sa pq.StringArray
		var fid int
		var fn string
		var fs string
		if err := rows.Scan(&sa, &fid, &fn, &fs); err != nil {
			return nil, err
		}
		fqr := &FlagQueryResult{}
		fqr.GroupHolders = []*GroupHolderFQResult{}
		fqr.FlagID = fid
		fqr.FlagName = fn
		err = json.Unmarshal([]byte(fs), &(fqr.Flag))
		if err != nil {
			return nil, err
		}
		for _, gh := range sa {
			r := &GroupHolderFQResult{}
			err := json.Unmarshal([]byte(gh), r)
			if err != nil {
				return nil, err
			}
			fqr.GroupHolders = append(fqr.GroupHolders, r)
		}
		s = append(s, fqr)
	}

	return s, nil
}

// HasFlag - checks if given user has a flag
func (db *Database) HasFlag(user string, channel string, flag string) (bool, error) {
	row := db.db.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM melodious.query_flags($1, $2, '', $3, true) LIMIT 1);
	`, user, channel, flag)

	var exists bool
	err := row.Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

// HasFlagChID - checks if given user has a flag
func (db *Database) HasFlagChID(user string, channel int, flag string) (bool, error) {
	row := db.db.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM melodious.query_flags($1, (SELECT name FROM melodious.channels WHERE id=$2 LIMIT 1), '', $3, true) LIMIT 1);
	`, user, channel, flag)

	var exists bool
	err := row.Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

// GetUser - gets a user's info by their ID
func (db *Database) GetUser(id int) (*User, error) {
	row := db.db.QueryRow(`
		SELECT id, username, owner FROM melodious.accounts WHERE id=$1;
	`, id)
	user := &User{}
	err := row.Scan(&(user.ID), &(user.Username), &(user.Owner))
	if err != nil {
		return &User{}, err
	}
	return user, nil
}

// GetGroupHolders - gets all group holders that exist
func (db *Database) GetGroupHolders() ([]*GroupHolder, error) {
	rows, err := db.db.Query(`
	SELECT 
		gh.id, 
		(SELECT name AS group FROM melodious.groups WHERE id=gh.group_id), 
		(SELECT username AS user FROM melodious.accounts WHERE id=gh.user_id), 
		(SELECT name AS channel FROM melodious.channels WHERE id=gh.channel_id) 
	FROM melodious.group_holders gh;
	`)
	if err != nil {
		return []*GroupHolder{}, err
	}
	ghs := []*GroupHolder{}
	for rows.Next() {
		gh := &GroupHolder{}
		// NULL handling
		var user interface{}
		var channel interface{}
		err = rows.Scan(&(gh.ID), &(gh.Group), &user, &channel)
		if err != nil {
			return []*GroupHolder{}, err
		}
		switch user.(type) {
		case string:
			gh.User = user.(string)
		}
		switch channel.(type) {
		case string:
			gh.Channel = channel.(string)
		}
		ghs = append(ghs, gh)
	}
	return ghs, nil
}

// NewDatabase - creates a new Database instance
func NewDatabase(mel *Melodious, addr string) (*Database, error) {
	db, err := sql.Open("postgres", addr)

	_, err = db.Exec(`CREATE SCHEMA IF NOT EXISTS melodious;`)
	if err != nil {
		return nil, err
	}
	log.Info("DB: check/create melodious schema")

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS melodious.accounts (
			id serial NOT NULL PRIMARY KEY,
			username varchar(32) NOT NULL UNIQUE,
			passhash varchar(64) NOT NULL,
			owner BOOLEAN NOT NULL,
			banned BOOLEAN NOT NULL DEFAULT false,
			ip inet NOT NULL DEFAULT '0.0.0.0'
		);`)
	if err != nil {
		return nil, err
	}
	log.Info("DB: check/create accounts table")

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS melodious.channels (
			id serial NOT NULL PRIMARY KEY,
			name varchar(32) NOT NULL UNIQUE,
			topic varchar(128) NOT NULL
		);`)
	if err != nil {
		return nil, err
	}
	log.Info("DB: check/create channels table")

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS melodious.messages (
			id serial NOT NULL PRIMARY KEY,
			chan_id int4 NOT NULL REFERENCES melodious.channels(id) ON DELETE CASCADE,
			author_id int4 NOT NULL REFERENCES melodious.accounts(id) ON DELETE CASCADE,
			message varchar(2048) NOT NULL,
			dt timestamp with time zone NOT NULL,
			pings varchar(32) []
		);`)
	if err != nil {
		return nil, err
	}
	log.Info("DB: check/create messages table")

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS melodious.groups (
			id serial NOT NULL PRIMARY KEY,
			name varchar(32) NOT NULL
		);`)
	if err != nil {
		return nil, err
	}
	log.Info("DB: check/create groups table")

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS melodious.group_holders (
			id serial NOT NULL PRIMARY KEY,
			group_id int4 NOT NULL REFERENCES melodious.groups(id) ON DELETE CASCADE,
			user_id int4 REFERENCES melodious.accounts(id) ON DELETE CASCADE,
			channel_id int4 REFERENCES melodious.channels(id) ON DELETE CASCADE,
			UNIQUE(group_id, user_id, channel_id)
		);`)
	if err != nil {
		return nil, err
	}
	log.Info("DB: check/create group_holders table")

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS melodious.group_flags (
			id serial NOT NULL PRIMARY KEY,
			group_id int4 NOT NULL REFERENCES melodious.groups(id) ON DELETE CASCADE,
			name varchar(32) NOT NULL,
			flag jsonb NOT NULL,
			UNIQUE(group_id, name)
		);`)
	if err != nil {
		return nil, err
	}
	log.Info("DB: check/create group_flags table")

	_, err = db.Exec(`
		CREATE OR REPLACE FUNCTION melodious.set_flag(group_name varchar(32), flag_name varchar(32), flag_data jsonb)
		RETURNS int4
		LANGUAGE plpgsql
		AS $$
		DECLARE
			fid int4 := NULL;
			gid int4 := NULL;
		BEGIN
			SELECT id INTO gid FROM melodious.groups WHERE name=group_name;
			SELECT id INTO fid FROM melodious.group_flags WHERE group_id=gid AND name=flag_name;
			IF fid IS NULL THEN
				INSERT INTO melodious.group_flags (group_id, name, flag) VALUES (gid, flag_name, flag_data) RETURNING id INTO fid;
				RETURN fid;
			ELSE
				UPDATE melodious.group_flags SET flag=flag_data WHERE id=fid;
				RETURN fid;
			END IF;
		END;
		$$;
	`)
	if err != nil {
		return nil, err
	}
	log.Info("DB: check/create set_flag function")

	_, err = db.Exec(`
		CREATE OR REPLACE PROCEDURE melodious.delete_flag(group_name varchar(32), flag_name varchar(32))
		LANGUAGE plpgsql
		AS $$
		DECLARE
			gid int4 := NULL;
		BEGIN
			SELECT id INTO gid FROM melodious.groups WHERE name=group_name;
			IF gid IS NOT NULL THEN
				DELETE FROM melodious.group_flags WHERE group_id=gid AND name=flag_name;
			END IF;
		END;
		$$;
	`)
	if err != nil {
		return nil, err
	}
	log.Info("DB: check/create delete_flag procedure")

	_, err = db.Exec(`
		CREATE OR REPLACE FUNCTION melodious.insert_group_holder(group_name varchar(32), user_name varchar(32), chan_name varchar(32))
		RETURNS int4
		LANGUAGE plpgsql
		AS $$
		DECLARE
			gid  int4 := NULL;
			uid  int4 := NULL;
			cid  int4 := NULL;
			ghid int4 := NULL;
		BEGIN
			SELECT id INTO gid FROM melodious.groups WHERE name=group_name;
			IF gid IS NULL THEN
				RAISE EXCEPTION 'no such group';
			END IF;
			IF user_name <> '' THEN
				SELECT id INTO uid FROM melodious.accounts WHERE username=user_name;
				IF uid IS NULL THEN
					RAISE EXCEPTION 'no such user';
				END IF;
			END IF;
			IF chan_name <> '' THEN
				SELECT id INTO cid FROM melodious.channels WHERE name=chan_name;
				IF cid IS NULL THEN
					RAISE EXCEPTION 'no such channel';
				END IF;
			END IF;
			INSERT INTO melodious.group_holders (group_id, user_id, channel_id) VALUES (gid, uid, cid) RETURNING id INTO ghid;
			RETURN ghid;
		END;
		$$;
	`)
	if err != nil {
		return nil, err
	}
	log.Info("DB: check/create insert_group_holder function")

	_, err = db.Exec(`
		CREATE OR REPLACE PROCEDURE melodious.delete_group_holders(group_name varchar(32), user_name varchar(32), chan_name varchar(32))
		LANGUAGE plpgsql
		AS $$
		DECLARE
			gid int4 := NULL;
			uid int4 := NULL;
			cid int4 := NULL;
		BEGIN
			IF group_name <> '' THEN
				SELECT id INTO gid FROM melodious.groups WHERE name=group_name;
				IF gid IS NULL THEN
					RAISE EXCEPTION 'no such group';
				END IF;
			END IF;
			IF user_name <> '' THEN
				SELECT id INTO uid FROM melodious.accounts WHERE username=user_name;
				IF uid IS NULL THEN
					RAISE EXCEPTION 'no such user';
				END IF;
			END IF;
			IF chan_name <> '' THEN
				SELECT id INTO cid FROM melodious.channels WHERE name=chan_name;
				IF cid IS NULL THEN
					RAISE EXCEPTION 'no such channel';
				END IF;
			END IF;
			
			IF group_name = '' AND user_name = '' AND chan_name = '' THEN
				DELETE FROM melodious.group_holders WHERE group_id=gid;
			ELSIF group_name <> '' AND user_name <> '' AND chan_name = '' THEN
				DELETE FROM melodious.group_holders WHERE group_id=gid AND user_id=uid;
			ELSIF group_name <> '' AND user_name = '' AND chan_name <> '' THEN
				DELETE FROM melodious.group_holders WHERE group_id=gid AND channel_id=cid;
			ELSIF group_name <> '' AND user_name <> '' AND chan_name <> '' THEN
				DELETE FROM melodious.group_holders WHERE group_id=gid AND user_id=uid AND channel_id=cid;
			ELSIF group_name = '' AND user_name = '' AND chan_name = '' THEN
				RAISE EXCEPTION 'cannot delete all group holders in a single request';
			ELSIF group_name = '' AND user_name <> '' AND chan_name = '' THEN
				DELETE FROM melodious.group_holders WHERE user_id=uid;
			ELSIF group_name = '' AND user_name = '' AND chan_name <> '' THEN
				DELETE FROM melodious.group_holders WHERE channel_id=cid;
			ELSIF group_name = '' AND user_name <> '' AND chan_name <> '' THEN
				DELETE FROM melodious.group_holders WHERE user_id=uid AND channel_id=cid;
			END IF;
		END;
		$$;
	`)
	if err != nil {
		return nil, err
	}
	log.Info("DB: check/create delete_group_holders function")

	_, err = db.Exec(`
		CREATE OR REPLACE FUNCTION melodious.query_flags(user_name varchar(32), chan_name varchar(32), igroup_name varchar(32), iflag_name varchar(32), flagcheck bool)
		RETURNS TABLE (
			group_holders jsonb [],
			flag_id int4,
			flag_name varchar(32),
			flag jsonb
		)
		LANGUAGE plpgsql
		AS $$
		DECLARE
			uid int4 := NULL;
			cid int4 := NULL;
			gid int4 := NULL;
		BEGIN
			IF user_name <> '' THEN
				SELECT id INTO uid FROM melodious.accounts WHERE username=user_name;
				IF uid IS NULL THEN
					RAISE EXCEPTION 'no such user';
				END IF;
			END IF;
			IF chan_name <> '' THEN
				SELECT id INTO cid FROM melodious.channels WHERE name=chan_name;
				IF cid IS NULL THEN
					RAISE EXCEPTION 'no such channel';
				END IF;
			END IF;
			IF igroup_name <> '' THEN
				SELECT id INTO gid FROM melodious.groups WHERE name=igroup_name;
				IF gid IS NULL THEN
					RAISE EXCEPTION 'go such group';
				END IF;
			END IF;

			IF igroup_name = '' THEN
				IF iflag_name = '' THEN
					IF user_name = '' AND chan_name = '' THEN
						RETURN QUERY SELECT
							ARRAY_AGG(JSONB_SET(ROW_TO_JSON(gh)::JSONB, '{group_name}'::TEXT[], TO_JSONB((SELECT name FROM melodious.groups WHERE id=gh.group_id LIMIT 1)))) AS group_holders,
							gf.id flag_id,
							gf.name flag_name,
							gf.flag flag
						FROM melodious.group_holders gh
						INNER JOIN melodious.group_flags gf
						ON gh.group_id = gf.group_id
						WHERE (NOT flagcheck) OR (flagcheck AND gh.user_id IS NULL AND gh.channel_id IS NULL)
						GROUP BY gf.id, gf.name;
					ELSIF user_name <> '' AND chan_name <> '' THEN
						RETURN QUERY SELECT
							ARRAY_AGG(JSONB_SET(ROW_TO_JSON(gh)::JSONB, '{group_name}'::TEXT[], TO_JSONB((SELECT name FROM melodious.groups WHERE id=gh.group_id LIMIT 1)))) AS group_holders,
							gf.id flag_id,
							gf.name flag_name,
							gf.flag flag
						FROM melodious.group_holders gh
						INNER JOIN melodious.group_flags gf
						ON gh.group_id = gf.group_id
						WHERE (gh.user_id = uid AND gh.channel_id = cid)
							 OR (flagcheck AND gh.user_id = uid AND gh.channel_id IS NULL)
							 OR (flagcheck AND gh.user_id IS NULL AND gh.channel_id = cid)
							 OR (flagcheck AND gh.user_id IS NULL AND gh.channel_id IS NULL)
						GROUP BY gf.id, gf.name;
					ELSIF user_name <> '' AND chan_name = '' THEN
						RETURN QUERY SELECT
							ARRAY_AGG(JSONB_SET(ROW_TO_JSON(gh)::JSONB, '{group_name}'::TEXT[], TO_JSONB((SELECT name FROM melodious.groups WHERE id=gh.group_id LIMIT 1)))) AS group_holders,
							gf.id flag_id,
							gf.name flag_name,
							gf.flag flag
						FROM melodious.group_holders gh
						INNER JOIN melodious.group_flags gf
						ON gh.group_id = gf.group_id
						WHERE (NOT flagcheck AND gh.user_id = uid)
							 OR (flagcheck AND gh.user_id = uid AND gh.channel_id IS NULL)
							 OR (flagcheck AND gh.user_id IS NULL AND gh.channel_id IS NULL)
						GROUP BY gf.id, gf.name;
					ELSIF user_name = '' AND chan_name <> '' THEN
						RETURN QUERY SELECT
							ARRAY_AGG(JSONB_SET(ROW_TO_JSON(gh)::JSONB, '{group_name}'::TEXT[], TO_JSONB((SELECT name FROM melodious.groups WHERE id=gh.group_id LIMIT 1)))) AS group_holders,
							gf.id flag_id,
							gf.name flag_name,
							gf.flag flag
						FROM melodious.group_holders gh
						INNER JOIN melodious.group_flags gf
						ON gh.group_id = gf.group_id
						WHERE (NOT flagcheck AND gh.channel_id = cid)
							 OR (flagcheck AND gh.user_id IS NULL AND gh.channel_id = cid)
							 OR (flagcheck AND gh.user_id IS NULL AND gh.channel_id IS NULL)
						GROUP BY gf.id, gf.name;
					END IF;
				ELSE
					IF user_name = '' AND chan_name = '' THEN
						RETURN QUERY SELECT
							ARRAY_AGG(JSONB_SET(ROW_TO_JSON(gh)::JSONB, '{group_name}'::TEXT[], TO_JSONB((SELECT name FROM melodious.groups WHERE id=gh.group_id LIMIT 1)))) AS group_holders,
							gf.id flag_id,
							gf.name flag_name,
							gf.flag flag
						FROM melodious.group_holders gh
						INNER JOIN melodious.group_flags gf
						ON gh.group_id = gf.group_id
						WHERE (gf.name = iflag_name)
							AND (
								(NOT flagcheck)
								OR (flagcheck AND gh.user_id IS NULL AND gh.channel_id IS NULL)
							)
						GROUP BY gf.id, gf.name;
					ELSIF user_name <> '' AND chan_name <> '' THEN
						RETURN QUERY SELECT
							ARRAY_AGG(JSONB_SET(ROW_TO_JSON(gh)::JSONB, '{group_name}'::TEXT[], TO_JSONB((SELECT name FROM melodious.groups WHERE id=gh.group_id LIMIT 1)))) AS group_holders,
							gf.id flag_id,
							gf.name flag_name,
							gf.flag flag
						FROM melodious.group_holders gh
						INNER JOIN melodious.group_flags gf
						ON gh.group_id = gf.group_id
						WHERE gf.name = iflag_name
						  AND (
								(gh.user_id = uid AND gh.channel_id = cid)
							 	OR (flagcheck AND gh.user_id = uid AND gh.channel_id IS NULL)
								OR (flagcheck AND gh.user_id IS NULL AND gh.channel_id = cid)
								OR (flagcheck AND gh.user_id IS NULL AND gh.channel_id IS NULL)
							)
						GROUP BY gf.id, gf.name;
					ELSIF user_name <> '' AND chan_name = '' THEN
						RETURN QUERY SELECT
							ARRAY_AGG(JSONB_SET(ROW_TO_JSON(gh)::JSONB, '{group_name}'::TEXT[], TO_JSONB((SELECT name FROM melodious.groups WHERE id=gh.group_id LIMIT 1)))) AS group_holders,
							gf.id flag_id,
							gf.name flag_name,
							gf.flag flag
						FROM melodious.group_holders gh
						INNER JOIN melodious.group_flags gf
						ON gh.group_id = gf.group_id
						WHERE gf.name = iflag_name
							AND (
								(NOT flagcheck AND gh.user_id = uid)
								OR (flagcheck AND gh.user_id = uid AND gh.channel_id IS NULL)
								OR (flagcheck AND gh.user_id IS NULL AND gh.channel_id IS NULL)
							)
						GROUP BY gf.id, gf.name;
					ELSIF user_name = '' AND chan_name <> '' THEN
						RETURN QUERY SELECT
							ARRAY_AGG(JSONB_SET(ROW_TO_JSON(gh)::JSONB, '{group_name}'::TEXT[], TO_JSONB((SELECT name FROM melodious.groups WHERE id=gh.group_id LIMIT 1)))) AS group_holders,
							gf.id flag_id,
							gf.name flag_name,
							gf.flag flag
						FROM melodious.group_holders gh
						INNER JOIN melodious.group_flags gf
						ON gh.group_id = gf.group_id
						WHERE gf.name = iflag_name
							AND (
								(NOT flagcheck AND gh.channel_id = cid)
								OR (flagcheck AND gh.channel_id = cid AND gh.user_id IS NULL)
								OR (flagcheck AND gh.channel_id IS NULL AND gh.user_id IS NULL)
							)
						GROUP BY gf.id, gf.name;
					END IF;
				END IF;
			ELSE
				IF iflag_name = '' THEN
					IF user_name = '' AND chan_name = '' THEN
						RETURN QUERY SELECT
							ARRAY_AGG(JSONB_SET(ROW_TO_JSON(gh)::JSONB, '{group_name}'::TEXT[], TO_JSONB(igroup_name))) AS group_holders,
							gf.id flag_id,
							gf.name flag_name,
							gf.flag flag
						FROM melodious.group_holders gh
						INNER JOIN melodious.group_flags gf
						ON gh.group_id = gf.group_id
						WHERE gh.group_id = gid
							AND (
								(NOT flagcheck)
								OR (flagcheck AND gh.user_id IS NULL AND gh.channel_id IS NULL)
							)
						GROUP BY gf.id, gf.name;
					ELSIF user_name <> '' AND chan_name <> '' THEN
						RETURN QUERY SELECT
							ARRAY_AGG(JSONB_SET(ROW_TO_JSON(gh)::JSONB, '{group_name}'::TEXT[], TO_JSONB(igroup_name))) AS group_holders,
							gf.id flag_id,
							gf.name flag_name,
							gf.flag flag
						FROM melodious.group_holders gh
						INNER JOIN melodious.group_flags gf
						ON gh.group_id = gf.group_id
						WHERE gh.group_id = gid
							AND (
								(gh.channel_id = cid AND gh.user_id = uid)
								OR (flagcheck AND gh.channel_id = cid AND gh.user_id IS NULL)
								OR (flagcheck AND gh.channel_id IS NULL AND gh.user_id = uid)
								OR (flagcheck AND gh.channel_id IS NULL AND gh.user_id IS NULL)
							)
						GROUP BY gf.id, gf.name;
					ELSIF user_name <> '' AND chan_name = '' THEN
						RETURN QUERY SELECT
							ARRAY_AGG(JSONB_SET(ROW_TO_JSON(gh)::JSONB, '{group_name}'::TEXT[], TO_JSONB(igroup_name))) AS group_holders,
							gf.id flag_id,
							gf.name flag_name,
							gf.flag flag
						FROM melodious.group_holders gh
						INNER JOIN melodious.group_flags gf
						ON gh.group_id = gf.group_id
						WHERE gh.group_id = gid
							AND (
								(NOT flagcheck AND gh.user_id = uid)
								OR (flagcheck AND gh.user_id = uid AND gh.channel_id IS NULL)
								OR (flagcheck AND gh.user_id IS NULL AND gh.channel_id IS NULL)
							)
						GROUP BY gf.id, gf.name;
					ELSIF user_name = '' AND chan_name <> '' THEN
						RETURN QUERY SELECT
							ARRAY_AGG(JSONB_SET(ROW_TO_JSON(gh)::JSONB, '{group_name}'::TEXT[], TO_JSONB(igroup_name))) AS group_holders,
							gf.id flag_id,
							gf.name flag_name,
							gf.flag flag
						FROM melodious.group_holders gh
						INNER JOIN melodious.group_flags gf
						ON gh.group_id = gf.group_id
						WHERE gh.group_id = gid
							AND (
								(NOT flagcheck AND gh.channel_id = cid)
								OR (flagcheck AND gh.channel_id = cid AND gh.user_id IS NULL)
								OR (flagcheck AND gh.channel_id IS NULL AND gh.user_id IS NULL) 
							)
						GROUP BY gf.id, gf.name;
					END IF;
				ELSE
					IF user_name = '' AND chan_name = '' THEN
						RETURN QUERY SELECT
							ARRAY_AGG(JSONB_SET(ROW_TO_JSON(gh)::JSONB, '{group_name}'::TEXT[], TO_JSONB(igroup_name))) AS group_holders,
							gf.id flag_id,
							gf.name flag_name,
							gf.flag flag
						FROM melodious.group_holders gh
						INNER JOIN melodious.group_flags gf
						ON gh.group_id = gf.group_id
						INNER JOIN melodious.groups g
						ON g.id = gh.group_id
						WHERE (gf.name = iflag_name AND gh.group_id = gid)
							AND (
								(NOT flagcheck)
								OR (flagcheck AND gh.channel_id IS NULL AND gh.user_id IS NULL)
							)
						GROUP BY gf.id, gf.name;
					ELSIF user_name <> '' AND chan_name <> '' THEN
						RETURN QUERY SELECT
							ARRAY_AGG(JSONB_SET(ROW_TO_JSON(gh)::JSONB, '{group_name}'::TEXT[], TO_JSONB(igroup_name))) AS group_holders,
							gf.id flag_id,
							gf.name flag_name,
							gf.flag flag
						FROM melodious.group_holders gh
						INNER JOIN melodious.group_flags gf
						ON gh.group_id = gf.group_id
						WHERE (gf.name = iflag_name AND gh.group_id = gid)
							AND (
								(gh.user_id = uid AND gh.channel_id = cid)
								OR (flagcheck AND gh.user_id = uid AND gh.channel_id IS NULL)
								OR (flagcheck AND gh.user_id IS NULL AND gh.channel_id = cid)
								OR (flagcheck AND gh.user_id IS NULL AND gh.channel_id IS NULL)
							)
						GROUP BY gf.id, gf.name;
					ELSIF user_name <> '' AND chan_name = '' THEN
						RETURN QUERY SELECT
							ARRAY_AGG(JSONB_SET(ROW_TO_JSON(gh)::JSONB, '{group_name}'::TEXT[], TO_JSONB(igroup_name))) AS group_holders,
							gf.id flag_id,
							gf.name flag_name,
							gf.flag flag
						FROM melodious.group_holders gh
						INNER JOIN melodious.group_flags gf
						ON gh.group_id = gf.group_id
						WHERE (gh.group_id = gid AND gf.name = iflag_name)
							AND (
								(NOT flagcheck AND gh.user_id = uid)
								OR (flagcheck AND gh.user_id = uid AND gh.channel_id IS NULL)
								OR (flagcheck AND gh.user_id IS NULL AND gh.channel_id IS NULL)
							)
						GROUP BY gf.id, gf.name;
					ELSIF user_name = '' AND chan_name <> '' THEN
						RETURN QUERY SELECT
							ARRAY_AGG(JSONB_SET(ROW_TO_JSON(gh)::JSONB, '{group_name}'::TEXT[], TO_JSONB(igroup_name))) AS group_holders,
							gf.id flag_id,
							gf.name flag_name,
							gf.flag flag
						FROM melodious.group_holders gh
						INNER JOIN melodious.group_flags gf
						ON gh.group_id = gf.group_id
						WHERE (gf.name = iflag_name AND gh.group_id = gid)
							AND (
								(NOT flagcheck AND gh.channel_id = cid)
								OR (flagcheck AND gh.channel_id = cid AND gh.user_id IS NULL)
								OR (flagcheck AND gh.channel_id IS NULL AND gh.user_id IS NULL)
							)
						GROUP BY gf.id, gf.name;
					END IF;
				END IF;
			END IF;
		END;
		$$;
	`)
	if err != nil {
		return nil, err
	}
	log.Info("DB: check/create query_flag(5) function")

	_, err = db.Exec(`
		CREATE OR REPLACE FUNCTION melodious.query_flags(user_name varchar(32), chan_name varchar(32), igroup_name varchar(32), iflag_name varchar(32))
		RETURNS TABLE (
			group_holders jsonb [],
			flag_id int4,
			flag_name varchar(32),
			flag jsonb
		)
		LANGUAGE plpgsql
		AS $$
		BEGIN
			RETURN QUERY SELECT melodious.query_flags(user_name, chan_name, igroup_name, iflag_name, false);
		END;
		$$;
	`)
	if err != nil {
		return nil, err
	}
	log.Info("DB: check/create query_flag(4) function")

	dbi := &Database{
		mel: mel,
		db:  db,
	}

	go func() {
		dhe, err := time.ParseDuration(mel.Config.DeleteHistoryEvery)
		if err != nil {
			log.WithField("err", err).Fatal("cannot parse duration")
		}
		shf := mel.Config.StoreHistoryFor
		for {
			func() {
				var err error
				defer log.WithFields(log.Fields{
					"storing-for":    shf,
					"deleting-every": dhe,
				}).Trace("deleting old messages").Stop(&err)
				err = dbi.DeleteOldMessages(shf)
			}()
			time.Sleep(dhe)
		}
	}()

	return dbi, nil
}
