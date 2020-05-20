package deleteobject

import (
	"os"
	"reflect"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jackstenglein/rest_api_creator/backend/auth"
	"github.com/jackstenglein/rest_api_creator/backend/errors"
)

type deleteObjectMockFunc func(string, string, string, auth.VerifyCookieFunc, deleteObjectDatabase) error

func deleteObjectMock(wantCookie string, wantPID string, wantOID string, err error) deleteObjectMockFunc {
	return func(cookie string, pid string, oid string, _ auth.VerifyCookieFunc, _ deleteObjectDatabase) error {
		if cookie != wantCookie || pid != wantPID || oid != wantOID {
			return errors.NewServer("Incorrect parameters passed to mock")
		}
		return err
	}
}

func handlerRequest(cookie string, pid string, oid string) events.APIGatewayProxyRequest {
	parameters := map[string]string{
		"pid": pid,
		"oid": oid,
	}
	headers := map[string]string{
		"Cookie": cookie,
	}
	return events.APIGatewayProxyRequest{PathParameters: parameters, Headers: headers}
}

func handlerResponse(err string, status int) events.APIGatewayProxyResponse {
	if err != "" {
		return events.APIGatewayProxyResponse{
			Body: `{"error":"` + err + `"}`,
			Headers: map[string]string{
				"Access-Control-Allow-Origin":      os.Getenv("CORS_ORIGIN"),
				"Access-Control-Allow-Credentials": "true",
			},
			StatusCode: status,
		}
	}
	return events.APIGatewayProxyResponse{
		Body: "{}",
		Headers: map[string]string{
			"Access-Control-Allow-Origin":      os.Getenv("CORS_ORIGIN"),
			"Access-Control-Allow-Credentials": "true",
		},
		StatusCode: status,
	}
}

var handleDeleteObjectTests = []struct {
	name string

	request          events.APIGatewayProxyRequest
	deleteObjectMock deleteObjectMockFunc

	wantResponse events.APIGatewayProxyResponse
	wantErr      error
}{
	{
		name:             "DeleteObjectFailure",
		request:          handlerRequest("session=cookievalue", "projectId", "objectId"),
		deleteObjectMock: deleteObjectMock("cookievalue", "projectId", "objectId", errors.Wrap(errors.NewClient("Invalid object ID"), "Failed database update")),
		wantResponse:     handlerResponse("Invalid object ID", 400),
	},
	{
		name:             "SuccessfulInvocation",
		request:          handlerRequest("session=cookievalue", "projectId", "objectId"),
		deleteObjectMock: deleteObjectMock("cookievalue", "projectId", "objectId", nil),
		wantResponse:     handlerResponse("", 200),
	},
}

func TestHandleDeleteObject(t *testing.T) {
	for _, test := range handleDeleteObjectTests {
		t.Run(test.name, func(t *testing.T) {
			// Setup
			deleteObjectFunc = test.deleteObjectMock
			defer func() {
				deleteObjectFunc = deleteObject
			}()

			// Execute
			response, err := HandleDeleteObject(test.request)

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
