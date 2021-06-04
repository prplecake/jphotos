package auth

import (
	"golang.org/x/crypto/bcrypt"

	"git.sr.ht/~mjorgensen/jphotos/db"
)

// AddUser hashes the password and adds the new user to the database
// Returns an error if a user of the same name already exists or a DB
// error occurs.
func AddUser(username, password string, db db.Store) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return db.AddUser(username, hash)
}

// RemoveUser removes a user from the database.
func RemoveUser(username string, db db.Store) error {
	return db.RemoveUser(username)
}
