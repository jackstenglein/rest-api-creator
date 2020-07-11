package deploy

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jackstenglein/rest_api_creator/backend/auth"
	"github.com/jackstenglein/rest_api_creator/backend/errors"
)

type deployFunc func(string, string, auth.VerifyCookieFunc, deployDatabase, deployer) (string, error)

func deployMock(wantCookie string, wantProjectID string, url string, err error) deployFunc {
	return func(cookie string, projectID string, _ auth.VerifyCookieFunc, _ deployDatabase, _ deployer) (string, error) {
		if cookie != wantCookie || projectID != projectID {
			return "", errors.NewServer("Incorrect parameters passed to mock")
		}
		return url, err
	}
}

func handlerRequest(cookie string, projectID string) events.APIGatewayProxyRequest {
	parameters := map[string]string{
		"pid": projectID,
	}
	headers := map[string]string{
		"Cookie": cookie,
	}
	return events.APIGatewayProxyRequest{PathParameters: parameters, Headers: headers}
}

func handlerResponse(url string, err string, status int) events.APIGatewayProxyResponse {
	json, _ := json.Marshal(&deployResponse{URL: url, Error: err})
	return events.APIGatewayProxyResponse{
		Body:       string(json),
		Headers:    map[string]string{"Access-Control-Allow-Origin": os.Getenv("CORS_ORIGIN"), "Access-Control-Allow-Credentials": "true"},
		StatusCode: status,
	}
}

var handleDeployTests = []struct {
	name string

	// Input
	request events.APIGatewayProxyRequest

	// Mock data
	deployMock deployFunc

	// Expected output
	wantResponse events.APIGatewayProxyResponse
	wantErr      error
}{
	{
		name:         "DeployProjectFailure",
		request:      handlerRequest("session=cookievalue", "projectId"),
		deployMock:   deployMock("cookievalue", "projectId", "", errors.NewServer("Failed database call")),
		wantResponse: handlerResponse("", "Failed database call", 500),
	},
	{
		name:         "SuccessfulInvocation",
		request:      handlerRequest("session=cookievalue", "projectId"),
		deployMock:   deployMock("cookievalue", "projectId", "example.com", nil),
		wantResponse: handlerResponse("example.com", "", 200),
	},
}

func TestHandleDeploy(t *testing.T) {
	for _, test := range handleDeployTests {
		t.Run(test.name, func(t *testing.T) {
			// Setup
			deploy = test.deployMock
			defer func() {
				deploy = deployProject
			}()

			// Execute
			response, err := HandleDeploy(test.request)

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
