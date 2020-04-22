package portal

import (
	"testing"

	"github.com/jackstenglein/rest_api_creator/backend-sls/dao"
	"github.com/jackstenglein/rest_api_creator/backend-sls/errors"
)

type loginDBMock struct {
	email          string
	token          string
	user           *dao.User
	getUserErr     error
	updateTokenErr error
}

func (mock *loginDBMock) GetUser(email string) (*dao.User, error) {
	if email != mock.email {
		return nil, errors.NewServer("Incorrect input to GetUser mock")
	}
	return mock.user, mock.getUserErr
}

func (mock *loginDBMock) UpdateUserToken(email string, token string) error {
	if email != mock.email || token != mock.token {
		return errors.NewServer("Incorrect input to UpdateUserToken mock")
	}
	return mock.updateTokenErr
}

var testUser = dao.User{Email: "test@example.com", Password: "$2a$14$MNkzNEv8Su7mHfLPIdWoU.t5lElbvlnDka11w27zgfy6Sw44zZsku"}

var loginTests = []struct {
	name string

	// Input
	email          string
	password       string
	db             loginDatabase
	generateToken  generateTokenFunc
	generateCookie generateCookieFunc

	// Expected output
	wantCookie string
	wantErr    error
}{
	{
		name:    "EmptyEmail",
		wantErr: errors.NewClient("Email and password parameters are required"),
	},
	{
		name:    "EmptyPassword",
		email:   "test@example.com",
		wantErr: errors.NewClient("Email and password parameters are required"),
	},
	{
		name:     "GetUserError",
		email:    "test@example.com",
		password: "12345678",
		db:       &loginDBMock{"test@example.com", "", nil, errors.NewServer("DB failure"), nil},
		wantErr:  errors.Wrap(errors.NewServer("DB failure"), "Failed to get user"),
	},
	{
		name:     "IncorrectPassword",
		email:    "test@example.com",
		password: "incorrect",
		db:       &loginDBMock{"test@example.com", "", &testUser, nil, nil},
		wantErr:  errors.NewClient("Incorrect email or password"),
	},
	{
		name:          "GenerateTokenError",
		email:         "test@example.com",
		password:      "12345678",
		db:            &loginDBMock{"test@example.com", "", &testUser, nil, nil},
		generateToken: generateTokenMock("", errors.NewServer("GenerateToken failure")),
		wantErr:       errors.Wrap(errors.NewServer("GenerateToken failure"), "Failed to create auth token"),
	},
	{
		name:           "GenerateCookieError",
		email:          "test@example.com",
		password:       "12345678",
		db:             &loginDBMock{"test@example.com", "", &testUser, nil, nil},
		generateToken:  generateTokenMock("token", nil),
		generateCookie: generateCookieMock("test@example.com", "token", "", errors.NewServer("GenerateCookie failure")),
		wantErr:        errors.Wrap(errors.NewServer("GenerateCookie failure"), "Failed to create cookie"),
	},
	{
		name:           "UpdateTokenError",
		email:          "test@example.com",
		password:       "12345678",
		db:             &loginDBMock{"test@example.com", "token", &testUser, nil, errors.NewServer("UpdateUserToken failure")},
		generateToken:  generateTokenMock("token", nil),
		generateCookie: generateCookieMock("test@example.com", "token", "cookie", nil),
		wantErr:        errors.Wrap(errors.NewServer("UpdateUserToken failure"), "Failed to update auth token"),
	},
	{
		name:           "SuccessfulInvocation",
		email:          "test@example.com",
		password:       "12345678",
		db:             &loginDBMock{"test@example.com", "token", &testUser, nil, nil},
		generateToken:  generateTokenMock("token", nil),
		generateCookie: generateCookieMock("test@example.com", "token", "cookie", nil),
		wantCookie:     "cookie",
	},
}

func TestLogin(t *testing.T) {
	for _, test := range loginTests {
		t.Run(test.name, func(t *testing.T) {
			// Execute
			cookie, err := login(test.email, test.password, test.generateToken, test.generateCookie, test.db)

			// Verify
			if cookie != test.wantCookie {
				t.Errorf("Got cookie '%s'; want '%s'", cookie, test.wantCookie)
			}
			if !errors.Equal(err, test.wantErr) {
				t.Errorf("Got err %v; want %v", err, test.wantErr)
			}
		})
	}
}
