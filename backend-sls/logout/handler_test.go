package logout

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jackstenglein/rest_api_creator/backend-sls/errors"
)

type logoutMockFunc func(string, verifyCookieFunc, logoutDatabase) error

func logoutMock(wantCookie string, err error) logoutMockFunc {
	return func(cookie string, verifyCookie verifyCookieFunc, db logoutDatabase) error {
		if cookie != wantCookie {
			return errors.NewServer(fmt.Sprintf("Incorrect parameters passed to mock: got '%s'; want '%s'", cookie, wantCookie))
		}
		return err
	}
}

func handlerRequest(cookie string) events.APIGatewayProxyRequest {
	headers := map[string]string{
		"Cookie": cookie,
	}
	return events.APIGatewayProxyRequest{Headers: headers}
}

func handlerResponse(err string, status int) events.APIGatewayProxyResponse {
	json, _ := json.Marshal(&logoutResponse{Error: err})
	return events.APIGatewayProxyResponse{Body: string(json), StatusCode: status}
}

var handleLogoutTests = []struct {
	name string

	// Input
	request events.APIGatewayProxyRequest

	// Mock data
	logoutMock logoutMockFunc

	// Expected output
	wantResponse events.APIGatewayProxyResponse
	wantErr      error
}{
	{
		name:         "InvalidCookie",
		request:      handlerRequest(";HttpOnly;session=asdfasdfasdf"),
		logoutMock:   logoutMock("", errors.NewClient("Not authenticated")),
		wantResponse: handlerResponse("Not authenticated", 400),
	},
	{
		name:         "SuccessfulInvocation",
		request:      handlerRequest("session=cookievalue;HttpOnly;SameSite=strict;Secure"),
		logoutMock:   logoutMock("cookievalue", nil),
		wantResponse: handlerResponse("", 200),
	},
}

func TestHandleLogout(t *testing.T) {
	for _, test := range handleLogoutTests {
		t.Run(test.name, func(t *testing.T) {
			// Setup
			logoutFunc = test.logoutMock
			defer func() {
				logoutFunc = logout
			}()

			// Execute
			response, err := HandleLogout(test.request)

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
