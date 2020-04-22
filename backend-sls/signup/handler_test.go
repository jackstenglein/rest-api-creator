package signup

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jackstenglein/rest_api_creator/backend-sls/errors"
)

type signupFunction func(string, string) (string, error)

func signupMock(email string, password string, cookie string, err error) signupFunction {
	return func(gotEmail string, gotPassword string) (string, error) {
		if gotEmail != email || gotPassword != password {
			return "", errors.NewServer("Incorrect input to signup mock")
		}
		return cookie, err
	}
}

func handlerRequest(email string, password string) events.APIGatewayProxyRequest {
	json, _ := json.Marshal(&signupRequest{Email: email, Password: password})
	return events.APIGatewayProxyRequest{Body: string(json)}
}

func handlerResponse(cookie string, err string, status int) events.APIGatewayProxyResponse {
	json, _ := json.Marshal(&signupResponse{Error: err})
	var headers = map[string]string{
		"Set-Cookie": fmt.Sprintf("session=%s;HttpOnly;SameSite=strict;Secure", cookie),
	}
	return events.APIGatewayProxyResponse{Headers: headers, Body: string(json), StatusCode: status}
}

var handlerTests = []struct {
	name string

	// Input
	request events.APIGatewayProxyRequest

	// Mock data
	mockSignup signupFunction

	// Expected output
	wantResponse events.APIGatewayProxyResponse
	wantErr      error
}{
	{
		name:         "ClientError",
		request:      handlerRequest("test@example.com", "1234567"),
		mockSignup:   signupMock("test@example.com", "1234567", "", errors.Wrap(errors.NewClient("Password is too short"), "Invalid password")),
		wantResponse: handlerResponse("", "Password is too short", 400),
	},
	{
		name:         "ServerError",
		request:      handlerRequest("test@example.com", "1234567"),
		mockSignup:   signupMock("test@example.com", "1234567", "", errors.Wrap(errors.NewServer("Random error"), "Failed to create user")),
		wantResponse: handlerResponse("", "Failed to create user", 500),
	},
	{
		name:         "SuccessfulInvocation",
		request:      handlerRequest("test@example.com", "1234567"),
		mockSignup:   signupMock("test@example.com", "1234567", "cookievalue", nil),
		wantResponse: handlerResponse("cookievalue", "", 200),
	},
}

func TestHandleRequest(t *testing.T) {
	for _, test := range handlerTests {
		t.Run(test.name, func(t *testing.T) {
			// Setup
			signupFunc = test.mockSignup
			defer func() {
				signupFunc = signup
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
