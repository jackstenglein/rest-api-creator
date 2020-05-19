package getuser

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jackstenglein/rest_api_creator/backend/dao"
	"github.com/jackstenglein/rest_api_creator/backend/errors"
)

type getUserMocker func(string, verifyCookieFunc, getUserDatabase) (*dao.User, error)

func getUserMock(wantCookie string, user *dao.User, err error) getUserMocker {
	return func(cookie string, _ verifyCookieFunc, _ getUserDatabase) (*dao.User, error) {
		if cookie != wantCookie {
			return nil, errors.NewServer("Incorrect input to get user mock.")
		}
		return user, err
	}
}

func handlerRequest(cookie string) events.APIGatewayProxyRequest {
	headers := map[string]string{
		"Cookie": cookie,
	}
	return events.APIGatewayProxyRequest{Headers: headers}
}

func handlerResponse(user *dao.User, err string, status int) events.APIGatewayProxyResponse {
	json, _ := json.Marshal(&getUserResponse{User: user, Error: err})
	return events.APIGatewayProxyResponse{
		Body:       string(json),
		Headers:    map[string]string{"Access-Control-Allow-Origin": os.Getenv("CORS_ORIGIN"), "Access-Control-Allow-Credentials": "true"},
		StatusCode: status,
	}
}

var handlerTests = []struct {
	name string

	// Input
	request events.APIGatewayProxyRequest

	// Mock data
	getUserMock getUserMocker

	// Expected output
	wantResponse events.APIGatewayProxyResponse
	wantErr      error
}{
	{
		name:         "ActionError",
		request:      handlerRequest("session=cookievalue"),
		getUserMock:  getUserMock("cookievalue", nil, errors.NewServer("DB failure")),
		wantResponse: handlerResponse(nil, "DB failure", 500),
	},
	{
		name:         "SuccessfulInvocation",
		request:      handlerRequest("session=cookievalue"),
		getUserMock:  getUserMock("cookievalue", &dao.User{Email: "test@example.com"}, nil),
		wantResponse: handlerResponse(&dao.User{Email: "test@example.com"}, "", 200),
	},
}

func TestHandleGetUser(t *testing.T) {
	for _, test := range handlerTests {
		t.Run(test.name, func(t *testing.T) {
			// Setup
			getUserFunc = test.getUserMock
			defer func() {
				getUserFunc = getUser
			}()

			// Execute
			response, err := HandleGetUser(test.request)

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
