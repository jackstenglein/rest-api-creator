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
		name:          "InvalidCookie",
		request:       handlerRequest("default", ";HttpOnly;session=asdfasdfasdf"),
		mockProjectID: "default",
		mockErr:       errors.NewClient("Not authenticated"),
		wantResponse:  handlerResponse(nil, "Not authenticated", 400),
	},
	{
		name:          "InvalidCookie2",
		request:       handlerRequest("default", "session=;HttpOnly;"),
		mockProjectID: "default",
		mockErr:       errors.NewClient("Not authenticated"),
		wantResponse:  handlerResponse(nil, "Not authenticated", 400),
	},
	{
		name:          "InvalidCookie3",
		request:       handlerRequest("default", ";HttpOnly;"),
		mockProjectID: "default",
		mockErr:       errors.NewClient("Not authenticated"),
		wantResponse:  handlerResponse(nil, "Not authenticated", 400),
	},
	{
		name:          "SucccessfulInvocation",
		request:       handlerRequest("default", "session=asdfjklqwerty;HttpOnly;"),
		mockProjectID: "default",
		mockCookie:    "asdfjklqwerty",
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
