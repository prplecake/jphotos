// Package db defines the requirements for our database, via the Store
// interface, and implements it for postgres.
package db

import (
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

// IsExpired returns true is a session has expired
func (s Session) IsExpired() bool {
	return time.Now().After(s.Expires)
}

// A Store provides the methods required to access the database.
type Store interface {
	ExecuteSchema(filename string) error
	AddUser(username string, hash []byte) error
	GetUserByName(username string) (*User, error)
	UserAddSession(user User, session string, expires time.Time) error

	GetGroupsForUser(u User) ([]Group, error)
	GetGroupByID(is string) (Group, []GroupMember, error)

	// SessionGet returns a valid session if one exists.
	// Guranteed to not return expired sessinos.
	// If a valid session is found, extend it! I don't recommend passing
	// in a time that's past, though.
	SessionGet(session string, newExpiration time.Time) (*Session, error)
	RevokeSession(session string) error
}
