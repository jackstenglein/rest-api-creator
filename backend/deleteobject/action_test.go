package deleteobject

import (
	"reflect"
	"testing"

	"github.com/jackstenglein/rest_api_creator/backend/auth"
	"github.com/jackstenglein/rest_api_creator/backend/dao"
	"github.com/jackstenglein/rest_api_creator/backend/errors"
)

func verifyCookieMock(mockCookie string, mockDB auth.UserGetter, mockEmail string, mockErr error) auth.VerifyCookieFunc {
	return func(cookie string, db auth.UserGetter) (string, error) {
		if cookie != mockCookie || !reflect.DeepEqual(db, mockDB) {
			return "", errors.NewServer("Incorrect input to VerifyCookie mock")
		}
		return mockEmail, mockErr
	}
}

type databaseMock struct {
	email     string
	projectID string
	objectID  string
	err       error
}

func (mock *databaseMock) GetUserInfo(email string) (*dao.User, error) {
	return nil, nil
}

func (mock *databaseMock) DeleteObject(email string, projectID string, objectID string) error {
	if email != mock.email || projectID != mock.projectID || objectID != mock.objectID {
		return errors.NewServer("Incorrect input to DeleteObject mock.")
	}
	return mock.err
}

var deleteObjectTests = []struct {
	name string

	cookie    string
	projectID string
	objectID  string

	db        *databaseMock
	email     string
	verifyErr error

	wantErr error
}{
	{
		name:    "EmptyProjectID",
		wantErr: errors.NewClient("Parameters `projectID` and `objectID` are both required"),
	},
	{
		name:      "EmptyObjectID",
		projectID: "project",
		wantErr:   errors.NewClient("Parameters `projectID` and `objectID` are both required"),
	},
	{
		name:      "InvalidCookie",
		cookie:    "invalidCookie",
		projectID: "project",
		objectID:  "object",
		verifyErr: errors.NewClient("Invalid cookie format"),
		wantErr:   errors.NewClient("Not authenticated"),
	},
	{
		name:      "DatabaseFailure",
		cookie:    "validCookie",
		projectID: "project",
		objectID:  "object",
		db:        &databaseMock{email: "test@example.com", projectID: "project", objectID: "object", err: errors.NewServer("Database failure")},
		email:     "test@example.com",
		wantErr:   errors.Wrap(errors.NewServer("Database failure"), "Failed to delete object in database"),
	},
	{
		name:      "SuccessfulInvocation",
		cookie:    "validCookie",
		projectID: "project",
		objectID:  "object",
		db:        &databaseMock{email: "test@example.com", projectID: "project", objectID: "object"},
		email:     "test@example.com",
	},
}

func TestDeleteObject(t *testing.T) {
	for _, test := range deleteObjectTests {
		t.Run(test.name, func(t *testing.T) {
			// Setup
			verifyCookie := verifyCookieMock(test.cookie, test.db, test.email, test.verifyErr)

			// Execute
			err := deleteObject(test.cookie, test.projectID, test.objectID, verifyCookie, test.db)

			// Verify
			if !errors.Equal(err, test.wantErr) {
				t.Errorf("Got err `%v`; want `%v`", err, test.wantErr)
			}
		})
	}
}
