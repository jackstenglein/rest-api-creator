package actions

import (
	"github.com/pkg/errors"
	"github.com/rest_api_creator/backend-sls/authentication"
	"github.com/rest_api_creator/backend-sls/dao"
	apierrors "github.com/rest_api_creator/backend-sls/errors"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest SignupRequest

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

func (action *LoginAction) Login(request LoginRequest) (string, error) {
	if request.Email == "" || request.Password == "" {
		return "", apierrors.NewUserError("Email and password are required")
	}

	user, err := action.store.GetUser(request.Email)
	if err != nil {
		return "", errors.Wrap(err, "Failed to get user")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		return "", apierrors.NewUserError("Incorrect email or password")
	}

	token, err := action.auth.GenerateToken()
	if err != nil {
		return "", errors.Wrap(err, "Failed to create auth token")
	}

	cookie, err := action.auth.GenerateCookie(request.Email, token)
	if err != nil {
		return "", errors.Wrap(err, "Failed to create cookie")
	}

	err = action.store.UpdateUserToken(request.Email, token)
	if err != nil {
		return "", errors.Wrap(err, "Failed to update auth token")
	}

	return cookie, nil
}
