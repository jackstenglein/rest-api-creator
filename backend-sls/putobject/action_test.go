package putobject

import (
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/jackstenglein/rest_api_creator/backend-sls/auth"
	"github.com/jackstenglein/rest_api_creator/backend-sls/dao"
	"github.com/jackstenglein/rest_api_creator/backend-sls/errors"
)

func verifyCookieMock(mockCookie string, mockDB auth.UserGetter, mockEmail string, mockErr error) verifyCookieFunc {
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
	object    *dao.Object
	err       error
}

func (mock *databaseMock) GetUserInfo(email string) (*dao.User, error) {
	return nil, nil
}

func (mock *databaseMock) UpdateObject(email string, projectID string, object *dao.Object) error {
	if email != mock.email || projectID != mock.projectID || !reflect.DeepEqual(object, mock.object) {
		return errors.NewServer("Incorrect input to UpdateObject mock")
	}
	return mock.err
}

var mockUUID = uuid.New()
var mockUUIDString = mockUUID.String()

func newUUIDMock(err error) uuidFunc {
	return func() (uuid.UUID, error) {
		return mockUUID, err
	}
}

var putObjectTests = []struct {
	name string

	// Input
	cookie    string
	projectID string
	object    *dao.Object

	// Mock data
	db        *databaseMock
	uuidFunc  uuidFunc
	email     string
	verifyErr error

	// Expected output
	wantErr error
}{
	{
		name:    "EmptyCookie",
		wantErr: errors.NewClient("Parameters `cookie`, `projectId` and `object` are required"),
	},
	{
		name:    "EmptyProjectID",
		cookie:  "cookie",
		wantErr: errors.NewClient("Parameters `cookie`, `projectId` and `object` are required"),
	},
	{
		name:      "NilObject",
		cookie:    "cookie",
		projectID: "projectId",
		wantErr:   errors.NewClient("Parameters `cookie`, `projectId` and `object` are required"),
	},
	{
		name:      "ObjectEmptyName",
		cookie:    "cookie",
		projectID: "projectId",
		object:    &dao.Object{},
		wantErr:   errors.NewClient("Object must have a `name` field"),
	},
	{
		name:      "UuidError",
		cookie:    "cookie",
		projectID: "projectId",
		object:    &dao.Object{Name: "New object", Description: "desc"},
		uuidFunc:  newUUIDMock(errors.NewServer("UUID failure")),
		wantErr:   errors.Wrap(errors.NewServer("UUID failure"), "Failed to generate UUID"),
	},
	{
		name:      "InvalidCookie",
		cookie:    "cookie",
		projectID: "projectId",
		object:    &dao.Object{ID: "id", Name: "name", Description: "desc"},
		db:        &databaseMock{},
		email:     "test@example.com",
		verifyErr: errors.NewClient("Invalid cookie"),
		wantErr:   errors.NewClient("Not authenticated"),
	},
	{
		name:      "DatabaseError",
		cookie:    "cookie",
		projectID: "projectId",
		object:    &dao.Object{ID: "id", Name: "name", Description: "desc"},
		db:        &databaseMock{"test@example.com", "projectId", &dao.Object{ID: "id", Name: "name", Description: "desc"}, errors.NewServer("DDB failure")},
		email:     "test@example.com",
		wantErr:   errors.Wrap(errors.NewServer("DDB failure"), "Failed database call to put object"),
	},
	{
		name:      "SuccessfulUpdate",
		cookie:    "cookie",
		projectID: "projectId",
		object:    &dao.Object{ID: "id", Name: "name", Description: "desc"},
		db:        &databaseMock{"test@example.com", "projectId", &dao.Object{ID: "id", Name: "name", Description: "desc"}, nil},
		email:     "test@example.com",
	},
	{
		name:      "SuccessfulCreate",
		cookie:    "cookie",
		projectID: "projectId",
		object:    &dao.Object{Name: "name", Description: "desc"},
		uuidFunc:  newUUIDMock(nil),
		db:        &databaseMock{"test@example.com", "projectId", &dao.Object{ID: mockUUIDString, Name: "name", Description: "desc"}, nil},
		email:     "test@example.com",
	},
}

func TestPutObject(t *testing.T) {
	for _, test := range putObjectTests {
		t.Run(test.name, func(t *testing.T) {
			// Setup
			verifyCookie := verifyCookieMock("cookie", test.db, test.email, test.verifyErr)

			// Execute
			err := putObject(test.cookie, test.projectID, test.object, verifyCookie, test.db, test.uuidFunc)

			// Verify
			if !errors.Equal(err, test.wantErr) {
				t.Errorf("Got err %v; want %v", err, test.wantErr)
			}
		})
	}
}
