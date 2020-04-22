package signup

import (
	"fmt"
	"strings"

	"github.com/jackstenglein/rest_api_creator/backend-sls/auth"
	"github.com/jackstenglein/rest_api_creator/backend-sls/dao"
	"github.com/jackstenglein/rest_api_creator/backend-sls/errors"
	"golang.org/x/crypto/bcrypt"
)

type signupDatabase interface {
	CreateUser(string, string, string) error
}

var generateToken = auth.GenerateToken
var generateCookie = auth.GenerateCookie
var db signupDatabase = dao.Dynamo

func validateEmail(email string) bool {
	if len(email) == 0 || !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return false
	}
	return true
}

func validatePassword(password string) error {
	if len(password) < 8 {
		return errors.NewClient("Password is too short")
	}
	return nil
}

func signup(email string, password string) (string, error) {
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
