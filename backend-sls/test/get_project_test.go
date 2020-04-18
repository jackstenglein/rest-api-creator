package test

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/rest_api_creator/backend-sls/dao"

	gomock "github.com/golang/mock/gomock"
	"github.com/rest_api_creator/backend-sls/actions"
	"github.com/rest_api_creator/backend-sls/mock"
)

var getProjectActionTests = []struct {
	name      string
	projectId string
	cookie    string

	verifyCookieResult bool
	verifyCookieErr    error
	verifyCookieCalls  int
	getProjectResult   *dao.Project
	getProjectErr      error
	getProjectCalls    int

	wantError   string
	wantProject *dao.Project
}{
	{
		name:      "EmptyProject",
		wantError: "Parameter `id` is required",
	},
	{
		name:      "InvalidCookieFormat",
		projectId: "projectId",
		cookie:    "incorrectcookietype",
		wantError: "Not authenticated",
	},
	{
		name:              "VerifyCookieError",
		projectId:         "projectId",
		cookie:            "email#token#hmac",
		verifyCookieErr:   errors.New("Unable to decode hex string"),
		verifyCookieCalls: 1,
		wantError:         "Failed to verify cookie",
	},
	{
		name:               "InvalidCookie",
		projectId:          "projectId",
		cookie:             "email#token#hmac",
		verifyCookieResult: false,
		verifyCookieCalls:  1,
		wantError:          "Not authenticated",
	},
	{
		name:               "DynamoError",
		projectId:          "project",
		cookie:             "email#token#hmac",
		verifyCookieResult: true,
		verifyCookieCalls:  1,
		getProjectErr:      errors.New("Dynamo failed"),
		getProjectCalls:    1,
		wantError:          "Failed to get project",
	},
	{
		name:               "Success",
		projectId:          "project",
		cookie:             "email#token#hmac",
		verifyCookieResult: true,
		verifyCookieCalls:  1,
		getProjectResult:   &dao.Project{Id: "project", Name: "Asdf"},
		getProjectCalls:    1,
		wantError:          "",
		wantProject:        &dao.Project{Id: "project", Name: "Asdf"},
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
			request := actions.GetProjectRequest{Id: test.projectId, Cookie: test.cookie}
			mockAuth.EXPECT().VerifyCookie("email", "token", "hmac", mockStore).Return(test.verifyCookieResult, test.verifyCookieErr).Times(test.verifyCookieCalls)
			mockStore.EXPECT().GetProject("email", test.projectId).Return(test.getProjectResult, test.getProjectErr).Times(test.getProjectCalls)

			// Perform the test
			project, err := action.GetProject(request)

			// Verify results
			if !reflect.DeepEqual(project, test.wantProject) {
				t.Errorf("Got project %v; want %v", project, test.wantProject)
			}
			if err == nil {
				if test.wantError != "" {
					t.Errorf("Got error nil; want '%s'", test.wantError)
				}
			} else if !strings.Contains(err.Error(), test.wantError) {
				t.Errorf("Got error %s; want '%s'", err, test.wantError)
			}
		})
	}
}
