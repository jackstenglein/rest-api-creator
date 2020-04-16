package actions

import (
	"golang.org/x/crypto/bcrypt"
	"github.com/rest_api_creator/backend-sls/authentication"
	"github.com/rest_api_creator/backend-sls/dao"
	"github.com/rest_api_creator/backend-sls/errors"
	"strings"
)

type SignupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignupResponse struct {
	Error string `json:"error,omitempty"`
}

type SignupAction struct {
	store dao.DataStore
	auth authentication.Authenticator
}

func validateEmail(email string) bool {
	if len(email) == 0 || !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return false
	}
	return true
}

func validatePassword(password string) errors.ApiError {
	if len(password) < 8 {
		return errors.NewUserError("Password is too short")
	}
	return nil
}

func DefaultSignupAction() *SignupAction {
	return &SignupAction{dao.NewDynamoStore(), authentication.NewSessionAuthenticator()}
}

func NewSignupAction(store dao.DataStore, auth authentication.Authenticator) *SignupAction {
	return &SignupAction{store, auth}
}

func (action *SignupAction) Signup(request SignupRequest) (SignupResponse, int) {
	ok := validateEmail(request.Email)
	if !ok {
		return SignupResponse{"Invalid email"}, 400
	}

	aerr := validatePassword(request.Password)
	if aerr != nil {
		return SignupResponse{aerr.Error()}, aerr.StatusCode()
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(request.Password), 14)
	if err != nil {
		return SignupResponse{"Failed to hash password"}, 500
	}

	token, err := action.auth.GenerateToken()
	if err != nil {
		return SignupResponse{"Failed to create auth token"}, 500
	}

	aerr = action.store.CreateUser(request.Email, string(bytes), token)
	if aerr != nil {
		return SignupResponse{aerr.Error()}, aerr.StatusCode()
	}

	return SignupResponse{""}, 200
}
