package signup

import (
	"testing"

	"github.com/jackstenglein/rest_api_creator/backend-sls/auth"
	"github.com/jackstenglein/rest_api_creator/backend-sls/dao"
	"github.com/jackstenglein/rest_api_creator/backend-sls/errors"
)

type databaseMock struct {
	email        string
	plaintextPwd string
	token        string
	err          error
}

func (mock *databaseMock) CreateUser(email string, hashedPwd string, token string) error {
	if email != mock.email || hashedPwd == mock.plaintextPwd || token != mock.token {
		return errors.NewServer("Incorrect input to CreateUser mock")
	}
	return mock.err
}

type tokenGenerator func() (string, error)

func generateTokenMock(mockToken string, mockErr error) tokenGenerator {
	return func() (string, error) {
		return mockToken, mockErr
	}
}

type cookieGenerator func(string, string) (string, error)

func generateCookieMock(mockEmail string, mockToken string, mockCookie string, mockErr error) cookieGenerator {
	return func(email string, token string) (string, error) {
		if email != mockEmail || token != mockToken {
			return "", errors.NewServer("Incorrect input to GenerateCookie mock")
		}
		return mockCookie, mockErr
	}
}

var signupTests = []struct {
	name string

	// Input
	email    string
	password string

	// Mock data
	tokenFunc  tokenGenerator
	cookieFunc cookieGenerator
	mockDB     *databaseMock

	// Expected output
	wantCookie string
	wantErr    error
}{
	{
		name:    "EmptyEmail",
		wantErr: errors.NewClient("Invalid email: ''"),
	},
	{
		name:    "EmailMissingAt",
		email:   "testexample.com",
		wantErr: errors.NewClient("Invalid email: 'testexample.com'"),
	},
	{
		name:    "EmailMissingDot",
		email:   "test@examplecom",
		wantErr: errors.NewClient("Invalid email: 'test@examplecom'"),
	},
	{
		name:     "ShortPassword",
		email:    "test@example.com",
		password: "1234567",
		wantErr:  errors.Wrap(errors.NewClient("Password is too short"), "Invalid password"),
	},
	{
		name:      "GenerateTokenError",
		email:     "test@example.com",
		password:  "12345678",
		tokenFunc: generateTokenMock("", errors.NewServer("Token failure")),
		wantErr:   errors.Wrap(errors.NewServer("Token failure"), "Failed to create auth token"),
	},
	{
		name:       "GenerateCookieError",
		email:      "test@example.com",
		password:   "12345678",
		tokenFunc:  generateTokenMock("token", nil),
		cookieFunc: generateCookieMock("test@example.com", "token", "", errors.NewServer("Cookie failure")),
		wantErr:    errors.Wrap(errors.NewServer("Cookie failure"), "Failed to create cookie"),
	},
	{
		name:       "CreateUserError",
		email:      "test@example.com",
		password:   "12345678",
		tokenFunc:  generateTokenMock("token", nil),
		cookieFunc: generateCookieMock("test@example.com", "token", "cookie", nil),
		mockDB:     &databaseMock{"test@example.com", "12345678", "token", errors.NewClient("Email already exists")},
		wantErr:    errors.Wrap(errors.NewClient("Email already exists"), "Failed to create user"),
	},
	{
		name:       "SuccessfulInvocation",
		email:      "test@example.com",
		password:   "12345678",
		tokenFunc:  generateTokenMock("token", nil),
		cookieFunc: generateCookieMock("test@example.com", "token", "cookie", nil),
		mockDB:     &databaseMock{"test@example.com", "12345678", "token", nil},
		wantCookie: "cookie",
	},
}

func TestSignup(t *testing.T) {
	for _, test := range signupTests {
		t.Run(test.name, func(t *testing.T) {
			// Setup
			generateToken = test.tokenFunc
			generateCookie = test.cookieFunc
			db = test.mockDB
			defer func() {
				generateToken = auth.GenerateToken
				generateCookie = auth.GenerateCookie
				db = dao.Dynamo
			}()

			// Execute
			cookie, err := signup(test.email, test.password)

			// Verify
			if cookie != test.wantCookie {
				t.Errorf("Got cookie '%s'; want '%s'", cookie, test.wantCookie)
			}
			if !errors.Equal(err, test.wantErr) {
				t.Errorf("Got error %v; want %v", err, test.wantErr)
			}
		})
	}
}
