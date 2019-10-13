// Package db defines the requirements for our database, via the Store
// interface, and implements it for postgres.
package db

import (
	"database/sql"
	"errors"
	"time"
)

// A User is a view into the details for a given user
type User struct {
	id, Username string
	Hash         []byte
}

// A Session is a view into a session
type Session struct {
	User    User
	Expires time.Time
}

// A Group  holds information necessary for roles
type Group struct {
	Name, UUID string
}

// A GroupMember is a member of a group
type GroupMember struct {
	Username string
	Admin    bool
}

var (
	// ErrNotFound is returned when the requested value isn't found
	ErrNotFound = errors.New("DB: Not Found")
)

// IsExpired returns true is a session has expired
func (s Session) IsExpired() bool {
	return time.Now().After(s.Expires)
}

// Query executes a raw query against the DB and returns the result
func (pg *PGStore) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return pg.conn.Query(query, args...)
}

// A Store provides the methods required to access the database.
type Store interface {
	ExecuteSchema(filename string) error
	AddUser(username string, hash []byte) error
	GetUserByName(username string) (*User, error)
	UserAddSession(user User, session string, expires time.Time) error

	AddAlbum(name string) error
	GetAlbums() ([]Album, error)
	GetAlbum(slug string) (*Album, error)
	GetAlbumPhotos(id string) ([]Photo, error)
	GetAlbumSlugByID(id string) (string, error)

	AddPhoto(p Photo, albumID string) error
	GetPhotoByID(id string) (*Photo, error)

	GetGroupsForUser(u User) ([]Group, error)
	GetGroupByID(id string) (Group, []GroupMember, error)

	// SessionGet returns a valid session if one exists.
	// Guranteed to not return expired sessinos.
	// If a valid session is found, extend it! I don't recommend passing
	// in a time that's past, though.
	SessionGet(session string, newExpiration time.Time) (*Session, error)
	RevokeSession(session string) error
}
