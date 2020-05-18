package putobject

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jackstenglein/rest_api_creator/backend-sls/dao"
	"github.com/jackstenglein/rest_api_creator/backend-sls/errors"
)

type putObjectMockFunc func(string, string, *dao.Object, verifyCookieFunc, putObjectDatabase) (string, error)

func putObjectMock(wantCookie string, wantProjectID string, wantObject *dao.Object, wantID string, err error) putObjectMockFunc {
	return func(cookie string, projectID string, object *dao.Object, verify verifyCookieFunc, db putObjectDatabase) (string, error) {
		if cookie != wantCookie || projectID != projectID || !reflect.DeepEqual(object, wantObject) {
			return "", errors.NewServer("Incorrect parameters passed to mock")
		}
		return wantID, err
	}
}

func handlerRequest(cookie string, projectID string, object *dao.Object) events.APIGatewayProxyRequest {
	parameters := map[string]string{
		"pid": projectID,
	}
	headers := map[string]string{
		"Cookie": cookie,
	}
	json, _ := json.Marshal(object)
	return events.APIGatewayProxyRequest{PathParameters: parameters, Headers: headers, Body: string(json)}
}

func handlerResponse(id string, err string, status int) events.APIGatewayProxyResponse {
	json, _ := json.Marshal(&putObjectResponse{ID: id, Error: err})
	return events.APIGatewayProxyResponse{
		Body:       string(json),
		Headers:    map[string]string{"Access-Control-Allow-Origin": os.Getenv("CORS_ORIGIN"), "Access-Control-Allow-Credentials": "true"},
		StatusCode: status,
	}
}

var handlePutObjectTests = []struct {
	name string

	// Input
	request events.APIGatewayProxyRequest

	// Mock data
	putObjectMock putObjectMockFunc

	// Expected output
	wantResponse events.APIGatewayProxyResponse
	wantErr      error
}{
	{
		name:          "PutObjectFailure",
		request:       handlerRequest("session=cookievalue", "projectId", nil),
		putObjectMock: putObjectMock("cookievalue", "projectId", nil, "", errors.NewServer("Failed database call")),
		wantResponse:  handlerResponse("", "Failed database call", 500),
	},
	{
		name:          "SuccessfulInvocation",
		request:       handlerRequest("session=cookievalue", "projectId", &dao.Object{ID: "objectId", Name: "objectName", Description: "desc"}),
		putObjectMock: putObjectMock("cookievalue", "projectId", &dao.Object{ID: "objectId", Name: "objectName", Description: "desc"}, "objectId", nil),
		wantResponse:  handlerResponse("objectId", "", 200),
	},
}

func TestHandlePutObject(t *testing.T) {
	for _, test := range handlePutObjectTests {
		t.Run(test.name, func(t *testing.T) {
			// Setup
			putObjectFunc = test.putObjectMock
			defer func() {
				putObjectFunc = putObject
			}()

			// Execute
			response, err := HandlePutObject(test.request)

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
