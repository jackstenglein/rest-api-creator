package putobject

import (
	"fmt"
	"reflect"
	"testing"

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
		fmt.Printf("Got %v; want %v", object, mock.object)
		return errors.NewServer("Incorrect input to UpdateObject mock.")
	}
	return mock.err
}

var putObjectTests = []struct {
	name string

	// Input
	cookie    string
	projectID string
	object    *dao.Object

	// Mock data
	db        *databaseMock
	email     string
	verifyErr error

	// Expected output
	wantID  string
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
		wantErr:   errors.Wrap(errors.NewClient("Object must have a `name` field"), "Object is invalid"),
	},
	{
		name:      "ObjectInvalidName",
		cookie:    "cookie",
		projectID: "projectId",
		object:    &dao.Object{Name: "New Object"},
		wantErr:   errors.Wrap(errors.NewClient("Object name `New Object` contains non-alphabetical characters"), "Object is invalid"),
	},
	{
		name:      "NilAttribute",
		cookie:    "cookie",
		projectID: "projectId",
		object: &dao.Object{
			Name:        "NewObject",
			Description: "desc",
			Attributes: []*dao.Attribute{
				nil,
			},
		},
		wantErr: errors.Wrap(errors.Wrap(errors.NewClient("Attribute cannot be nil"), "Object must have valid attributes"), "Object is invalid"),
	},
	{
		name:      "AttributeEmptyName",
		cookie:    "cookie",
		projectID: "projectId",
		object: &dao.Object{
			Name:        "NewObject",
			Description: "desc",
			Attributes: []*dao.Attribute{
				{Name: ""},
			},
		},
		wantErr: errors.Wrap(errors.Wrap(errors.NewClient("Attribute must have a `name` field"), "Object must have valid attributes"), "Object is invalid"),
	},
	{
		name:      "AttributeInvalidName",
		cookie:    "cookie",
		projectID: "projectId",
		object: &dao.Object{
			Name:        "NewObject",
			Description: "desc",
			Attributes: []*dao.Attribute{
				{Name: "Invalid Name"},
			},
		},
		wantErr: errors.Wrap(errors.Wrap(errors.NewClient("Attribute name `Invalid Name` contains non-alphabetical characters"), "Object must have valid attributes"), "Object is invalid"),
	},
	{
		name:      "AttributeInvalidType",
		cookie:    "cookie",
		projectID: "projectId",
		object: &dao.Object{
			Name:        "NewObject",
			Description: "desc",
			Attributes: []*dao.Attribute{
				{Name: "ValidName"},
			},
		},
		wantErr: errors.Wrap(errors.Wrap(errors.NewClient("Attribute type `` is not supported"), "Object must have valid attributes"), "Object is invalid"),
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
		object:    &dao.Object{Name: "name", Description: "desc"},
		db:        &databaseMock{"test@example.com", "projectId", &dao.Object{ID: "name", Name: "name", CodeName: "Name", Description: "desc"}, errors.NewServer("DDB failure")},
		email:     "test@example.com",
		wantErr:   errors.Wrap(errors.NewServer("DDB failure"), "Failed database call to put object"),
	},
	{
		name:      "SuccessfulUpdate",
		cookie:    "cookie",
		projectID: "projectId",
		object: &dao.Object{
			ID:          "id",
			Name:        "name",
			Description: "desc",
			Attributes: []*dao.Attribute{
				{Name: "ValidName", Type: "Text"},
			},
		},
		db: &databaseMock{
			"test@example.com",
			"projectId",
			&dao.Object{
				ID:          "name",
				Name:        "name",
				CodeName:    "Name",
				Description: "desc",
				Attributes: []*dao.Attribute{
					{Name: "ValidName", Type: "Text", CodeName: "validName"},
				},
			},
			nil,
		},
		email:  "test@example.com",
		wantID: "name",
	},
	{
		name:      "SuccessfulCreate",
		cookie:    "cookie",
		projectID: "projectId",
		object:    &dao.Object{Name: "name", Description: "desc"},
		db:        &databaseMock{"test@example.com", "projectId", &dao.Object{ID: "name", Name: "name", CodeName: "Name", Description: "desc"}, nil},
		email:     "test@example.com",
		wantID:    "name",
	},
}

func TestPutObject(t *testing.T) {
	for _, test := range putObjectTests {
		t.Run(test.name, func(t *testing.T) {
			// Setup
			verifyCookie := verifyCookieMock("cookie", test.db, test.email, test.verifyErr)

			// Execute
			id, err := putObject(test.cookie, test.projectID, test.object, verifyCookie, test.db)

			// Verify
			if id != test.wantID {
				t.Errorf("Got id '%s'; want '%s'", id, test.wantID)
			}
			if !errors.Equal(err, test.wantErr) {
				t.Errorf("Got err %v; want %v", err, test.wantErr)
			}
		})
	}
}
