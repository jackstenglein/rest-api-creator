package logout

import (
	"reflect"
	"testing"

	"github.com/jackstenglein/rest_api_creator/backend-sls/auth"
	"github.com/jackstenglein/rest_api_creator/backend-sls/dao"
	"github.com/jackstenglein/rest_api_creator/backend-sls/errors"
)

type databaseMock struct {
	email string
	token string
	err   error
}

func (mock *databaseMock) GetUserInfo(email string) (*dao.User, error) {
	return nil, nil
}

func (mock *databaseMock) UpdateUserToken(email string, token string) error {
	if email != mock.email || token != mock.token {
		return errors.NewServer("Incorrect input to UpdateUserToken mock")
	}
	return mock.err
}

func verifyCookieMock(mockCookie string, mockDB auth.UserGetter, mockEmail string, mockErr error) verifyCookieFunc {
	return func(cookie string, db auth.UserGetter) (string, error) {
		if cookie != mockCookie || !reflect.DeepEqual(db, mockDB) {
			return "", errors.NewServer("Incorrect input to VerifyCookie mock")
		}
		return mockEmail, mockErr
	}
}

var logoutTests = []struct {
	name string

	// Input
	cookie string

	// Mock data
	db        *databaseMock
	email     string
	verifyErr error

	// Expected output
	wantErr error
}{
	{
		name:    "EmptyCookie",
		wantErr: errors.NewClient("Parameter `cookie` is required"),
	},
	{
		name:      "VerifyCookieError",
		cookie:    "invalidcookie",
		verifyErr: errors.NewClient("Invalid cookie"),
		wantErr:   errors.NewClient("Not authenticated"),
	},
	{
		name:      "UpdateTokenError",
		cookie:    "validcookie",
		email:     "test@example.com",
		db:        &databaseMock{"test@example.com", "", errors.NewServer("DynamoDB failure")},
		verifyErr: nil,
		wantErr:   errors.Wrap(errors.NewServer("DynamoDB failure"), "Failed to remove auth token"),
	},
	{
		name:      "SuccessfulInvocation",
		cookie:    "validcookie",
		email:     "test@example.com",
		db:        &databaseMock{"test@example.com", "", nil},
		verifyErr: nil,
		wantErr:   nil,
	},
}

func TestLogout(t *testing.T) {
	for _, test := range logoutTests {
		t.Run(test.name, func(t *testing.T) {
			// Setup
			verifyCookie := verifyCookieMock(test.cookie, test.db, test.email, test.verifyErr)

			// Execute
			err := logout(test.cookie, verifyCookie, test.db)

			// Verify
			if !errors.Equal(err, test.wantErr) {
				t.Errorf("Got err %v; want %v", err, test.wantErr)
			}
		})
	}
}
