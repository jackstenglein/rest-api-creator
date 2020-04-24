package getproject

import (
	"github.com/jackstenglein/rest_api_creator/backend-sls/auth"
	"github.com/jackstenglein/rest_api_creator/backend-sls/dao"
	"github.com/jackstenglein/rest_api_creator/backend-sls/errors"
	"reflect"
	"testing"
)

type databaseMock struct {
	email   string
	id      string
	project *dao.Project
	err     error
}

func (mock *databaseMock) GetUserInfo(email string) (*dao.User, error) {
	return nil, nil
}

func (mock *databaseMock) GetProject(email string, id string) (*dao.Project, error) {
	if email != mock.email || id != mock.id {
		return nil, errors.NewServer("Incorrect parameters passed to mock")
	}
	return mock.project, mock.err
}

type verifyCookieFunc func(string, auth.UserGetter) (string, error)

func verifyCookieMock(wantCookie string, wantDB auth.UserGetter, email string, err error) verifyCookieFunc {
	return func(gotCookie string, gotDB auth.UserGetter) (string, error) {
		if gotCookie != wantCookie || !reflect.DeepEqual(gotDB, wantDB) {
			return "", errors.NewServer("Incorrect parameters passed to mock")
		}
		return email, err
	}
}

var getProjectTests = []struct {
	name string

	// Test input
	id     string
	cookie string

	// Mock data
	email     string
	cookieErr error
	project   *dao.Project
	dbErr     error

	// Expected results
	wantProject *dao.Project
	wantErr     error
}{
	{
		name:    "EmptyID",
		wantErr: errors.NewClient("Parameter `id` is required"),
	},
	{
		name:      "EmptyCookie",
		id:        "projectID",
		cookie:    "invalidCookie",
		cookieErr: errors.NewClient("Invalid cookie format"),
		wantErr:   errors.NewClient("Not authenticated"),
	},
	{
		name:    "DatabaseFailure",
		id:      "projectID",
		cookie:  "validCookie",
		email:   "test@example.com",
		dbErr:   errors.NewServer("DynamoDB error"),
		wantErr: errors.Wrap(errors.NewServer("DynamoDB error"), "Failed to get project"),
	},
	{
		name:        "SuccessfulInvocation",
		id:          "projectID",
		cookie:      "validCookie",
		email:       "test@example.com",
		project:     &dao.Project{ID: "projectID", Name: "Default"},
		wantProject: &dao.Project{ID: "projectID", Name: "Default"},
	},
}

func TestGetProject(t *testing.T) {
	for _, test := range getProjectTests {
		t.Run(test.name, func(t *testing.T) {
			// Setup
			dbMock := &databaseMock{test.email, test.id, test.project, test.dbErr}
			db = dbMock
			verifyCookie = verifyCookieMock(test.cookie, dbMock, test.email, test.cookieErr)
			defer func() {
				db = dao.Dynamo
				verifyCookie = auth.VerifyCookie
			}()

			// Execute
			gotProject, gotErr := getProject(test.id, test.cookie)

			// Verify
			if !reflect.DeepEqual(gotProject, test.wantProject) {
				t.Errorf("Got project %v; want %v", gotProject, test.wantProject)
			}
			if !errors.Equal(gotErr, test.wantErr) {
				t.Errorf("Got err %v; want %v", gotErr, test.wantErr)
			}
		})
	}
}
