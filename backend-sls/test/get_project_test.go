package test

import (
	"reflect"
	"testing"

	"github.com/rest_api_creator/backend-sls/dao"

	gomock "github.com/golang/mock/gomock"
	"github.com/rest_api_creator/backend-sls/actions"
	"github.com/rest_api_creator/backend-sls/errors"
	"github.com/rest_api_creator/backend-sls/mock"
)

var getProjectActionTests = []struct {
	name      string
	projectId string
	cookie    string

	verifyCookieResult bool
	verifyCookieErr    error
	verifyCookieCalls  int
	getProjectResult   dao.Project
	getProjectErr      error
	getProjectCalls    int

	wantError   error
	wantProject dao.Project
}{
	{
		name:      "EmptyProject",
		wantError: errors.NewUserError("project is required"),
	},
	{
		name:      "InvalidCookieFormat",
		projectId: "projectId",
		cookie:    "incorrectcookietype",
		wantError: errors.NewUserError("Not authenticated"),
	},
	{
		name:              "VerifyCookieError",
		projectId:         "projectId",
		cookie:            "email#token#hmac",
		verifyCookieErr:   errors.New("Unable to decode hex string"),
		verifyCookieCalls: 1,
		wantError:         errors.New("Unable to decode hex string"),
	},
	{
		name:               "InvalidCookie",
		projectId:          "projectId",
		cookie:             "email#token#hmac",
		verifyCookieResult: false,
		verifyCookieCalls:  1,
		wantError:          errors.NewUserError("Not authenticated"),
	},
	{
		name:               "DynamoError",
		projectId:          "project",
		cookie:             "email#token#hmac",
		verifyCookieResult: true,
		verifyCookieCalls:  1,
		getProjectErr:      errors.New("Dynamo failed"),
		getProjectCalls:    1,
		wantError:          errors.New("Dynamo failed"),
	},
	{
		name:               "Success",
		projectId:          "project",
		cookie:             "email#token#hmac",
		verifyCookieResult: true,
		verifyCookieCalls:  1,
		getProjectResult:   dao.Project{Id: "project", Name: "Asdf"},
		getProjectCalls:    1,
		wantError:          nil,
		wantProject:        dao.Project{Id: "project", Name: "Asdf"},
	},
}

func TestGetProjectAction(t *testing.T) {
	for _, test := range getProjectActionTests {
		t.Run(test.name, func(t *testing.T) {
			// Create test objects
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			mockStore := mock.NewMockDataStore(mockCtrl)
			mockAuth := mock.NewMockAuthenticator(mockCtrl)
			action := actions.NewGetProjectAction(mockStore, mockAuth)

			// Setup using test input and mock data
			request := actions.GetProjectRequest{ProjectId: test.projectId, Cookie: test.cookie}
			mockAuth.EXPECT().VerifyCookie("email", "token", "hmac", mockStore).Return(test.verifyCookieResult, test.verifyCookieErr).Times(test.verifyCookieCalls)
			mockStore.EXPECT().GetProject("email", test.projectId).Return(test.getProjectResult, test.getProjectErr).Times(test.getProjectCalls)

			// Perform the test
			response := action.GetProject(request)

			// Verify results
			if !reflect.DeepEqual(response.Error, test.wantError) {
				t.Errorf("Got error %v; want %v", response.Error, test.wantError)
			}
			if !reflect.DeepEqual(response.Project, test.wantProject) {
				t.Errorf("Got error %v; want %v", response.Project, test.wantProject)
			}
		})
	}
}
