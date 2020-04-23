package getproject

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jackstenglein/rest_api_creator/backend-sls/dao"
	"github.com/jackstenglein/rest_api_creator/backend-sls/errors"
)

type getProjectFunc func(string, string) (*dao.Project, error)

func getProjectMock(wantProjectID string, wantCookie string, output *dao.Project, err error) getProjectFunc {
	return func(gotProjectID string, gotCookie string) (*dao.Project, error) {
		if gotProjectID != wantProjectID || gotCookie != wantCookie {
			return nil, errors.NewServer("Incorrect parameters passed to mock")
		}
		return output, err
	}
}

func handlerRequest(projectID string, cookie string) events.APIGatewayProxyRequest {
	parameters := map[string]string{
		"id": projectID,
	}
	headers := map[string]string{
		"Cookie": cookie,
	}
	return events.APIGatewayProxyRequest{PathParameters: parameters, Headers: headers}
}

func handlerResponse(project *dao.Project, err string, status int) events.APIGatewayProxyResponse {
	json, _ := json.Marshal(&getProjectResponse{Project: project, Error: err})
	return events.APIGatewayProxyResponse{Body: string(json), StatusCode: status}
}

var handlerTests = []struct {
	name string

	// Input
	request events.APIGatewayProxyRequest

	// Mock data
	mockProjectID string
	mockCookie    string
	mockProject   *dao.Project
	mockErr       error

	// Expected output
	wantResponse events.APIGatewayProxyResponse
	wantErr      error
}{
	{
		name:          "GetProjectError",
		request:       handlerRequest("default", "session=cookie;HttpOnly;SameSite=strict;Secure"),
		mockProjectID: "default",
		mockCookie:    "cookie",
		mockErr:       errors.NewServer("Failed to get project"),
		wantResponse:  handlerResponse(nil, "Failed to get project", 500),
	},
	{
		name:          "SucccessfulInvocation",
		request:       handlerRequest("default", "session=cookie;HttpOnly;SameSite=strict;Secure"),
		mockProjectID: "default",
		mockCookie:    "cookie",
		mockProject:   &dao.Project{ID: "default", Name: "ProjectName"},
		wantResponse:  handlerResponse(&dao.Project{ID: "default", Name: "ProjectName"}, "", 200),
	},
}

func TestHandleRequest(t *testing.T) {
	for _, test := range handlerTests {
		t.Run("", func(t *testing.T) {
			// Setup
			actionFunc = getProjectMock(test.mockProjectID, test.mockCookie, test.mockProject, test.mockErr)
			defer func() {
				actionFunc = getProject
			}()

			// Execute
			response, err := HandleRequest(test.request)

			// Verify
			if !reflect.DeepEqual(response, test.wantResponse) {
				t.Errorf("Got response %v; want %v", response, test.wantResponse)
			}
			if !errors.Equal(err, test.wantErr) {
				t.Errorf("Got error %v; want %v", err, test.wantErr)
			}
		})
	}
}
