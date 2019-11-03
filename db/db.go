// Package db defines the requirements for our database, via the Store
// interface, and implements it for postgres.
package db

import (
	"database/sql"
	"errors"
	"log"
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
	log.Print("Running query: ", query)
	return pg.conn.Query(query, args...)
}

// Exec executes a raw query against the DB, returns the result, and
// closes the connection
func (pg *PGStore) Exec(query string, args ...interface{}) error {
	txn, err := pg.conn.Begin()
	if err != nil {
		log.Printf("Currently there are %d connections.", pg.conn.Stats().OpenConnections)
		return err
	}
	_, err = txn.Exec(query, args...)
	if err != nil {
		return err
	}
	err = txn.Commit()
	if err != nil {
		return err
	}
	return nil
}

// A Store provides the methods required to access the database.
type Store interface {
	ExecuteSchema(filename string) error
	AddUser(username string, hash []byte) error
	GetUserByName(username string) (*User, error)
	UserAddSession(user User, session string, expires time.Time) error

	//
	// Album Methods
	// These are methods used to primarily access the albums tabl
	GetAllAlbums() ([]Album, error)
	GetAlbumBySlug(slug string) (*Album, error)
	GetAlbumPhotosByID(id string) ([]Photo, error)
	GetAlbumSlugByID(id string) (string, error)

	AddAlbum(name string) error

	RenameAlbumByID(id, newName string) error

	DeleteAlbumBySlug(slug string) error

	//
	// Photo Methods
	// These are methods used to primarily access the photos table
	AddPhoto(p Photo, albumID string) error
	GetPhotoByID(id string) (*Photo, error)
	GetAlbumIDByPhotoID(id string) (string, error)

	UpdatePhotoCaptionByID(id, newCaption string) error
	UpdatePhotoAlbum(photoID, albumID string) error

	DeletePhotoByID(id string) error

	//
	// Group Methods
	GetGroupsForUser(u User) ([]Group, error)
	GetGroupByID(id string) (Group, []GroupMember, error)

	// SessionGet returns a valid session if one exists.
	// Guranteed to not return expired sessinos.
	// If a valid session is found, extend it! I don't recommend passing
	// in a time that's past, though.
	SessionGet(session string, newExpiration time.Time) (*Session, error)
	RevokeSession(session string) error
}
