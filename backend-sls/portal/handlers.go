// Package portal handles requests to the POST /signup and PUT /login REST API endpoints.
package portal

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jackstenglein/rest_api_creator/backend-sls/auth"
	"github.com/jackstenglein/rest_api_creator/backend-sls/dao"
	"github.com/jackstenglein/rest_api_creator/backend-sls/errors"
)

// portalRequest represents the HTTP request body when calling the signup or login APIs and is used for unmarshalling.
type portalRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// portalResponse represents the HTTP response body when calling the signup or login APIs and is used for marshalling.
type portalResponse struct {
	Error string `json:"error,omitempty"`
}

// portalFunc wraps the function signature for functions that perform portal actions.
type portalFunc func(email string, password string) (string, error)

// generateTokenFunc wraps the function signature for functions that generate auth tokens.
type generateTokenFunc func() (string, error)

// generateCookieFunc wraps the function signature for functions that generate cookies.
type generateCookieFunc func(string, string) (string, error)

// signupFunc points to the function used to perform the signup action.
// This variable should be changed only to perform dependency injection in unit tests.
var signupFunc = handleSignup

// loginFunc points to the function used to perform the login action.
// This variable should be changed only to perform dependency injection in unit tests.
var loginFunc = handleLogin

// HandleSignupRequest acts as a middle-man between AWS APIGateway and the signup action function. HandleSignupRequest
// unmarshals the request from APIGateway and forwards it to the signup action. HandleSignupRequest then marshals
// the response from the signup action and returns it to APIGateway. The request must contain the `email` and
// `password` body parameters. If the request succeeds, the response will have status 200 OK, the response body will be
// empty and the Set-Cookie header will contain the user's auth token. If the request fails, the response will have
// either a 400 or a 500 status, the body will have an `error` field, and the Set-Cookie header will contain an empty
// cookie. This function returns a non-nil error only if JSON marshaling of the response body fails.
func HandleSignupRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return handleRequest(request, signupFunc)
}

// HandleLoginRequest acts as a middle-man between AWS APIGateway and the login action function. HandleLoginRequest
// unmarshals the request from APIGateway and forwards it to the login action. HandleLoginRequest then marshals
// the response from the login action and returns it to APIGateway. The request must contain the `email` and
// `password` body parameters. If the request succeeds, the response will have status 200 OK, the response body
// will be empty and the Set-Cookie header will contain the user's auth token. If the request fails, the
// response will have either a 400 or a 500 status, the body will have an `error` field, and the Set-Cookie
// header will contain an empty cookie. This function returns a non-nil error only if JSON marshaling of the
// response body fails.
func HandleLoginRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return handleRequest(request, loginFunc)
}

// handleRequest is a helper for HandleSignupRequest and HandleLoginRequest. It parses the request object
// from AWS APIGateway and forwards it to the specified actionFunc. It then marshals the response from the
// actionFunc and returns the marshalled response to the caller.
func handleRequest(request events.APIGatewayProxyRequest, actionFunc portalFunc) (events.APIGatewayProxyResponse, error) {
	// Parse request
	var portalRequest portalRequest
	json.Unmarshal([]byte(request.Body), &portalRequest)

	// Execute action
	cookie, err := actionFunc(portalRequest.Email, portalRequest.Password)

	// Create response
	errString, status := errors.UserDetails(err)
	fmt.Println(errors.StackTrace(err))
	var headers = map[string]string{
		"Set-Cookie":                       fmt.Sprintf("session=%s;HttpOnly;", cookie),
		"Access-Control-Allow-Origin":      "http://jackstenglein-rest-api-creator.s3-website-us-east-1.amazonaws.com",
		"Access-Control-Allow-Credentials": "true",
	}
	json, err := json.Marshal(&portalResponse{Error: errString})
	if err != nil {
		fmt.Println(errors.StackTrace(errors.Wrap(err, "Failed to marshal signup response")))
	}
	return events.APIGatewayProxyResponse{Headers: headers, Body: string(json), StatusCode: status}, nil
}

// actionFunc for the signup action.
func handleSignup(email string, password string) (string, error) {
	return signup(email, password, auth.GenerateToken, auth.GenerateCookie, dao.Dynamo)
}

// actionFunc for the login action. This must be separated from handleSignup because signup and login
// take different interface types for the db parameter.
func handleLogin(email string, password string) (string, error) {
	return login(email, password, auth.GenerateToken, auth.GenerateCookie, dao.Dynamo)
}
