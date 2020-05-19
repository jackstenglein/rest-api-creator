package getuser

import (
	"reflect"
	"testing"

	"github.com/jackstenglein/rest_api_creator/backend/auth"
	"github.com/jackstenglein/rest_api_creator/backend/dao"
	"github.com/jackstenglein/rest_api_creator/backend/errors"
)

type databaseMock struct {
	email string
	user  *dao.User
	err   error
}

func (mock *databaseMock) GetUser(email string) (*dao.User, error) {
	return mock.user, mock.err
}

func (mock *databaseMock) GetUserInfo(email string) (*dao.User, error) {
	return nil, nil
}

func verifyCookieMock(mockCookie string, mockDB auth.UserGetter, mockEmail string, mockErr error) verifyCookieFunc {
	return func(cookie string, db auth.UserGetter) (string, error) {
		if cookie != mockCookie || !reflect.DeepEqual(db, mockDB) {
			return "", errors.NewServer("Incorrect input to VerifyCookie mock")
		}
		return mockEmail, mockErr
	}
}

var getUserTests = []struct {
	name string

	// Input
	cookie    string
	email     string
	db        *databaseMock
	verifyErr error

	// Expected output
	wantErr  error
	wantUser *dao.User
}{
	{
		name:    "EmptyCookie",
		wantErr: errors.NewClient("Not authenticated"),
	},
	{
		name:      "InvalidCookie",
		cookie:    "invalidcookie",
		verifyErr: errors.NewClient("Not authenticated"),
		wantErr:   errors.Wrap(errors.NewClient("Not authenticated"), "Failed to verify cookie"),
	},
	{
		name:      "DatabaseError",
		cookie:    "validcookie",
		email:     "test@example.com",
		db:        &databaseMock{"test@example.com", nil, errors.NewServer("DynamoDB failure")},
		verifyErr: nil,
		wantErr:   errors.Wrap(errors.NewServer("DynamoDB failure"), "Failed to get user from database"),
	},
	{
		name:      "SuccessfulInvocation",
		cookie:    "validcookie",
		email:     "test@example.com",
		db:        &databaseMock{"test@example.com", &dao.User{Email: "test@example.com"}, nil},
		verifyErr: nil,
		wantUser:  &dao.User{Email: "test@example.com"},
	},
}

func TestGetUser(t *testing.T) {
	for _, test := range getUserTests {
		t.Run(test.name, func(t *testing.T) {
			// Setup
			verifyCookie := verifyCookieMock(test.cookie, test.db, test.email, test.verifyErr)

			// Execute
			user, err := getUser(test.cookie, verifyCookie, test.db)

			// Verify
			if !reflect.DeepEqual(user, test.wantUser) {
				t.Errorf("Got user %v; want %v", user, test.wantUser)
			}
			if !errors.Equal(err, test.wantErr) {
				t.Errorf("Got err %v; want %v", err, test.wantErr)
			}
		})
	}
}
