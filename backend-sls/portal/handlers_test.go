package portal

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jackstenglein/rest_api_creator/backend-sls/errors"
)

type handlerFunc func(events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

func portalFuncMock(email string, password string, cookie string, err error) portalFunc {
	return func(gotEmail string, gotPassword string) (string, error) {
		if gotEmail != email || gotPassword != password {
			return "", errors.NewServer("Incorrect input to signup mock")
		}
		return cookie, err
	}
}

func handlerRequest(email string, password string) events.APIGatewayProxyRequest {
	json, _ := json.Marshal(&portalRequest{Email: email, Password: password})
	return events.APIGatewayProxyRequest{Body: string(json)}
}

func handlerResponse(cookie string, err string, status int) events.APIGatewayProxyResponse {
	json, _ := json.Marshal(&portalResponse{Error: err})
	var headers = map[string]string{
		"Set-Cookie":                       fmt.Sprintf("session=%s;HttpOnly;", cookie),
		"Access-Control-Allow-Origin":      "http://jackstenglein-rest-api-creator.s3-website-us-east-1.amazonaws.com",
		"Access-Control-Allow-Credentials": "true",
	}
	return events.APIGatewayProxyResponse{Headers: headers, Body: string(json), StatusCode: status}
}

var portalHandlers = []struct {
	name     string
	function handlerFunc
}{
	{
		name:     "Signup",
		function: HandleSignupRequest,
	},
	{
		name:     "Login",
		function: HandleLoginRequest,
	},
}

var handlerTests = []struct {
	name string

	// Input
	request events.APIGatewayProxyRequest

	// Mock data
	mockFunc portalFunc

	// Expected output
	wantResponse events.APIGatewayProxyResponse
	wantErr      error
}{
	{
		name:         "ClientError",
		request:      handlerRequest("test@example.com", "1234567"),
		mockFunc:     portalFuncMock("test@example.com", "1234567", "", errors.Wrap(errors.NewClient("Password is too short"), "Invalid password")),
		wantResponse: handlerResponse("", "Password is too short", 400),
	},
	{
		name:         "ServerError",
		request:      handlerRequest("test@example.com", "1234567"),
		mockFunc:     portalFuncMock("test@example.com", "1234567", "", errors.Wrap(errors.NewServer("Random error"), "Failed to create user")),
		wantResponse: handlerResponse("", "Failed to create user", 500),
	},
	{
		name:         "SuccessfulInvocation",
		request:      handlerRequest("test@example.com", "1234567"),
		mockFunc:     portalFuncMock("test@example.com", "1234567", "cookievalue", nil),
		wantResponse: handlerResponse("cookievalue", "", 200),
	},
}

func TestHandleRequest(t *testing.T) {
	for _, test := range handlerTests {
		for _, handler := range portalHandlers {
			t.Run(fmt.Sprintf("%s/%s", test.name, handler.name), func(t *testing.T) {
				// Setup
				signupFunc = test.mockFunc
				loginFunc = test.mockFunc
				defer func() {
					signupFunc = handleSignup
					loginFunc = handleLogin
				}()

				// Execute
				response, err := handler.function(test.request)

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
}
