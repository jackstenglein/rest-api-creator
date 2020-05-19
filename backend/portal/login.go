package portal

import (
	"github.com/jackstenglein/rest_api_creator/backend/dao"
	"github.com/jackstenglein/rest_api_creator/backend/errors"
	"golang.org/x/crypto/bcrypt"
)

// loginDatabase wraps the database methods required to perform the login action.
type loginDatabase interface {
	GetUserInfo(string) (*dao.User, error)
	UpdateUserToken(string, string) error
}

// login performs the actual actions required to login a user. login checks the user's password against the
// hashed password stored in the database, it generates an auth token and cookie, and it updates the user's
// auth token in the database. If there are no errors, login returns the generated cookie. Otherwise, login
// returns the empty string and the error.
func login(email string, password string, generateToken generateTokenFunc, generateCookie generateCookieFunc, db loginDatabase) (string, error) {
	if email == "" || password == "" {
		return "", errors.NewClient("Email and password parameters are required")
	}

	user, err := db.GetUserInfo(email)
	if err != nil {
		return "", errors.Wrap(err, "Failed to get user")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.NewClient("Incorrect email or password")
	}

	token, err := generateToken()
	if err != nil {
		return "", errors.Wrap(err, "Failed to create auth token")
	}

	cookie, err := generateCookie(email, token)
	if err != nil {
		return "", errors.Wrap(err, "Failed to create cookie")
	}

	err = db.UpdateUserToken(email, token)
	if err != nil {
		return "", errors.Wrap(err, "Failed to update auth token")
	}

	return cookie, nil
}
