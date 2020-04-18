package actions

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/rest_api_creator/backend-sls/authentication"
	"github.com/rest_api_creator/backend-sls/dao"
	apierrors "github.com/rest_api_creator/backend-sls/errors"
	"golang.org/x/crypto/bcrypt"
)

type SignupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignupAction struct {
	store dao.DataStore
	auth  authentication.Authenticator
}

func validateEmail(email string) bool {
	if len(email) == 0 || !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return false
	}
	return true
}

func validatePassword(password string) error {
	if len(password) < 8 {
		return apierrors.NewUserError("Password is too short")
	}
	return nil
}

func DefaultSignupAction() *SignupAction {
	return &SignupAction{dao.DefaultDynamoStore(), authentication.NewSessionAuthenticator()}
}

func NewSignupAction(store dao.DataStore, auth authentication.Authenticator) *SignupAction {
	return &SignupAction{store, auth}
}

func (action *SignupAction) Signup(request SignupRequest) (string, error) {

	ok := validateEmail(request.Email)
	if !ok {
		return "", apierrors.NewUserError("Invalid email")
	}

	err := validatePassword(request.Password)
	if err != nil {
		return "", errors.Wrap(err, "Invalid password")
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(request.Password), 14)
	if err != nil {
		return "", errors.Wrap(err, "Failed to hash password")
	}

	token, err := action.auth.GenerateToken()
	if err != nil {
		return "", errors.Wrap(err, "Failed to create auth token")
	}

	err = action.store.CreateUser(request.Email, string(bytes), token)
	if err != nil {
		return "", errors.Wrap(err, "Failed to create user")
	}

	cookie, err := action.auth.GenerateCookie(request.Email, token)
	if err != nil {
		// This error doesn't affect the creation of the user, so don't return it
		fmt.Printf("%+v", err)
	}

	return cookie, nil
}
