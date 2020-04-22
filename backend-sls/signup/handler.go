package signup

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jackstenglein/rest_api_creator/backend-sls/errors"
)

type signupResponse struct {
	Error string `json:"error,omitempty"`
}

type signupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// signupFunc points to the function used to perform the signup action.
// This variable should be changed only to perform dependency injection in unit tests.
var signupFunc = signup

// HandleRequest parses the request object from AWS APIGateway and returns a response object containing
// the status code and the error message, if present. The request must contain the `email` and `password`
// body parameters. If the request succeeds, the response will have status 200 OK, the response body
// will be empty and the Set-Cookie header will contain the user's auth token. If the request fails, the
// response will have either a 400 or a 500 status, the body will have an `error` field, and the Set-Cookie
// header will contain an empty cookie. This function returns a non-nil error only if JSON marshaling of the
// response body fails.
func HandleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Parse request
	var signupRequest signupRequest
	json.Unmarshal([]byte(request.Body), &signupRequest)

	// Execute action
	cookie, err := signupFunc(signupRequest.Email, signupRequest.Password)

	// Create response
	errString, status := errors.UserDetails(err)
	fmt.Println(errors.StackTrace(err))
	var headers = map[string]string{
		"Set-Cookie": fmt.Sprintf("session=%s;HttpOnly;SameSite=strict;Secure", cookie),
	}
	json, err := json.Marshal(&signupResponse{Error: errString})
	if err != nil {
		fmt.Println(errors.StackTrace(errors.Wrap(err, "Failed to marshal signup response")))
	}
	return events.APIGatewayProxyResponse{Headers: headers, Body: string(json), StatusCode: status}, nil
}
