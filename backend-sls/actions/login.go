package actions

import (
	"github.com/rest_api_creator/backend-sls/authentication"
	"github.com/rest_api_creator/backend-sls/dao"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest SignupRequest
type LoginResponse SignupResponse

type LoginAction struct {
	store dao.DataStore
	auth  authentication.Authenticator
}

func DefaultLoginAction() *LoginAction {
	return &LoginAction{dao.DefaultDynamoStore(), authentication.NewSessionAuthenticator()}
}

func NewLoginAction(store dao.DataStore, auth authentication.Authenticator) *LoginAction {
	return &LoginAction{store, auth}
}

func (action *LoginAction) Login(request LoginRequest) (LoginResponse, string, int) {
	if request.Email == "" || request.Password == "" {
		return LoginResponse{"Email and password are required"}, "", 400
	}

	user, aerr := action.store.GetUser(request.Email)
	if aerr != nil {
		return LoginResponse{aerr.Error()}, "", aerr.StatusCode()
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		return LoginResponse{"Incorrect email or password"}, "", 400
	}

	token, err := action.auth.GenerateToken()
	if err != nil {
		return LoginResponse{"Failed to create auth token"}, "", 500
	}

	cookie, err := action.auth.GenerateCookie(request.Email, token)
	if err != nil {
		return LoginResponse{"Failed to create cookie"}, "", 500
	}

	err = action.store.UpdateUserToken(request.Email, token)
	if err != nil {
		return LoginResponse{err.Error()}, "", 500
	}

	return LoginResponse{}, cookie, 200
}
