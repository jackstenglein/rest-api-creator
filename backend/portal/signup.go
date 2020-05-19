package portal

import (
	"fmt"
	"strings"

	"github.com/jackstenglein/rest_api_creator/backend/errors"
	"golang.org/x/crypto/bcrypt"
)

// signupDatabase wraps the database methods required to perform the signup action.
type signupDatabase interface {
	CreateUser(string, string, string) error
}

// validateEmail returns true if the email is valid and false otherwise.
func validateEmail(email string) bool {
	if len(email) == 0 || !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return false
	}
	return true
}

// validatePassword checks that the user's password is valid. If it is not valid, it returns
// an error containing the reason. Otherwise, it returns nil.
func validatePassword(password string) error {
	if len(password) < 8 {
		return errors.NewClient("Password is too short")
	}
	return nil
}

// signup performs the actual actions required to create a new user. signup hashes the user's password, generates
// an auth token and cookie, and stores the new user in the database. If there are no errors, signup returns the
// generated cookie. Otherwise, signup returns the empty string and the error. If a user with the given email
// already exists in the database, an error is returned.
func signup(email string, password string, generateToken generateTokenFunc, generateCookie generateCookieFunc, db signupDatabase) (string, error) {
	ok := validateEmail(email)
	if !ok {
		return "", errors.NewClient(fmt.Sprintf("Invalid email: '%s'", email))
	}

	err := validatePassword(password)
	if err != nil {
		return "", errors.Wrap(err, "Invalid password")
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", errors.Wrap(err, "Failed to hash password")
	}

	token, err := generateToken()
	if err != nil {
		return "", errors.Wrap(err, "Failed to create auth token")
	}

	cookie, err := generateCookie(email, token)
	if err != nil {
		return "", errors.Wrap(err, "Failed to create cookie")
	}

	err = db.CreateUser(email, string(bytes), token)
	if err != nil {
		return "", errors.Wrap(err, "Failed to create user")
	}

	return cookie, nil
}
