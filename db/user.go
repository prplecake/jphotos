package db

import (
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
)

var (
	// ErrUsernameExists is returned when the unique requirement of a username is violated
	ErrUsernameExists = errors.New("DB: Username exists")
	// ErrNotFound is returned when the requested value isn't found
	ErrNotFound = errors.New("DB: Not Found")
)

// Query executes a raw query against the DB and returns the result
func (pg *PGStore) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return pg.conn.Query(query, args...)
}

// AddUser adds a user to the users table of the database, hashing the password
// with bcrypt
func (pg *PGStore) AddUser(username string, hash []byte) error {
	now := time.Now()
	_, err := pg.conn.Query("INSERT INTO users (username, hash, created, last_login) VALUES ($1, $2, $3, $4)",
		username, string(hash), now, now)
	if err, ok := err.(*pq.Error); ok {
		if err.Code == "23505" {
			return ErrUsernameExists
		}
	}
	return err
}

// GetUserByName returns the DB user information for a user if that user exists
func (pg *PGStore) GetUserByName(username string) (*User, error) {

	rows, err := pg.Query("SELECT id, hash FROM users WHERE username = $1", username)
	if err != nil {
		return nil, err
	}

	if !rows.Next() {
		return nil, ErrNotFound
	}

	var hash []byte
	var id string

	err = rows.Scan(&id, &hash)
	if err != nil {
		return nil, err
	}

	if rows.Next() {
		panic("Database guarantee not met: multiple users with same name")
	}

	return &User{
		Username: username,
		id:       id,
		Hash:     hash,
	}, nil
}

// UserAddSession stores a session in the database
func (pg *PGStore) UserAddSession(user User, session string, expires time.Time) error {
	_, err := pg.Query("INSERT INTO sessions VALUES ($1, $2, $3)",
		user.id, session, expires)
	if err != nil {
		return err
	}
	_, err = pg.Query("UPDATE users SET last_login = NOW() WHERE id = $1", user.id)
	return err
}

// SessionGet checks the database for a session and returns it if found
// If the session is absent, an error is returned.
// SessionGet will not return an expired session.
func (pg *PGStore) SessionGet(session string, newExpiration time.Time) (*Session, error) {
	// TODO: Is this a bad thing to do?
	_, err := pg.Query("DELETE FROM sessions WHERE expires < NOW()")
	if err != nil {
		panic(err)
	}
	_, err = pg.Query(
		"UPDATE sessions"+
			" SET expires = $1"+
			" WHERE token = $2", newExpiration, session)
	if err != nil {
		return nil, err
	}

	rows, err := pg.Query(
		"SELECT username, id, expires FROM sessions"+
			" INNER JOIN users ON sessions.user_id = users.id"+
			" WHERE token = $1", session)
	if err != nil {
		return nil, err
	}
	if !rows.Next() {
		return nil, ErrNotFound
	}

	var user User
	var expires time.Time

	err = rows.Scan(&user.Username, &user.id, &expires)
	if err != nil {
		return nil, err
	}
	//TODO: Update session to expire later?
	return &Session{
		User:    user,
		Expires: expires,
	}, nil
}

// RevokeSession removes a user's session from the sessions table
func (pg *PGStore) RevokeSession(session string) error {
	_, err := pg.Query("DELETE FROM sessions WHERE token = $1", session)
	return err
}
