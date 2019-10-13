package auth

import (
	"errors"
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"

	"git.sr.ht/~mjorgensen/jphotos/db"
)

const (
	// SessionCookieName is the name of the session cookie expected for authentication
	SessionCookieName = "jphotos_session"

	// ExpirationTime is the duration a session cookie is valid for once issued
	ExpirationTime = 60 * time.Minute
)

// A Token holds a session token
type Token struct {
	Session string
	Expires time.Time
}

// Error types returned by functions in auth
var (
	ErrInvalidUsernameOrPassword = errors.New("Invalid username or password")
	ErrUnauthorized              = errors.New("Unauthorized")
)

// NewSession creates a new session and inserts it into the session table
// Possible errors:
// * ErrInvalidUsernameOrPassword
// * Errors from database access
func NewSession(username, password string, s db.Store) (*Token, error) {
	user, err := s.GetUserByName(username)
	if err != nil {
		if err == db.ErrNotFound {
			return nil, ErrInvalidUsernameOrPassword
		}
		return nil, err
	}

	valid := bcrypt.CompareHashAndPassword(user.Hash, []byte(password))
	if valid != nil {
		return nil, ErrInvalidUsernameOrPassword
	}

	// TODO this seems a little ridiculout to pull a dep for...
	sessionToken := uuid.NewV4().String()
	expires := time.Now().Add(ExpirationTime)

	s.UserAddSession(*user, sessionToken, expires)

	return &Token{sessionToken, expires}, nil
}

func verifySessionCookie(r *http.Request, s db.Store) (string, string, error) {
	c, err := r.Cookie(SessionCookieName)
	if err == http.ErrNoCookie {
		return "", "", ErrUnauthorized
	} else if err != nil {
		return "", "", err
	}

	session, err := s.SessionGet(c.Value, time.Now().Add(ExpirationTime))
	if err != nil {
		return "", "", ErrUnauthorized
	}

	return session.User.Username, c.Value, nil
}
