package auth

import (
	"net/http"

	"git.sr.ht/~mjorgensen/jphotos/db"
)

// A Role represents a user's maximum permission level
type Role int

// Permissions work as follows:
// - Invalid:   Deactivated or otherwise invalid users
// - User:      Normal users
// - SuperUser: Has all server permissions
const (
	RoleInvalid Role = iota
	RoleUser
	RoleSuperUser
)

func (r Role) String() string {
	switch r {
	case RoleSuperUser:
		return "administrator"
	case RoleUser:
		return "user"
	default:
		return "invalid"
	}
}

// An Authorization represents a users' current authorization and authentication level
type Authorization struct {
	User    db.User
	session string
	Role    Role
}

// Revoke invalidates a session token in the database
func (a Authorization) Revoke(db db.Store) {
	db.RevokeSession(a.session)
}

// for testing: create a new sub-DB interface for auth - possibly in DB package?
// with just the stuff needed by the Auth package for the database - only used
// internally though

// Get returns a user's Authorization based on a session and a minimum required
// Role. If the session is valid, extend it.
func Get(r *http.Request, minimumRole Role, db db.Store) (*Authorization, error) {
	username, session, err := verifySessionCookie(r, db)
	if err != nil {
		return nil, err
	}

	user, err := db.GetUserByName(username)
	if err != nil {
		return nil, err
	}

	return &Authorization{
		User:    *user,
		session: session,
		Role:    RoleUser, //TODO - verify role from DB
	}, nil

}
