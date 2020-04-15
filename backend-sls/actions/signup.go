package actions

import (
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

func validateEmail(email string) bool {
	if len(email) == 0 || !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return false
	}
	return true
}

func validatePassword(password string) errors.ApiError {
	if len(password) < 8 {
		return errors.NewUserError("Password too short")
	}
	return nil
}

func Signup(request SignupRequest) (SignupResponse, int) {
	ok := validateEmail(request.Email)
	if !ok {
		return SignupResponse{"Invalid email"}, 400
	}

	err := validatePassword(request.Password)
	if err != nil {
		return SignupResponse{err.Error()}, err.StatusCode()
	}

	err = dao.CreateUser(request.Email, request.Password)
	if err != nil {
		return SignupResponse{err.Error()}, err.StatusCode()
	}

	return SignupResponse{""}, 200
}
