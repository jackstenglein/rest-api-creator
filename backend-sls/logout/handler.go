// Package logout handles requests to the PUT /logout REST API endpoint.
package logout

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jackstenglein/rest_api_creator/backend-sls/auth"
	"github.com/jackstenglein/rest_api_creator/backend-sls/dao"
	"github.com/jackstenglein/rest_api_creator/backend-sls/errors"
)

// logoutResponse contains the fields returned in the API JSON response body.
type logoutResponse struct {
	Error string `json:"error,omitempty"`
}

// logoutFunc points to the function used to perform the logout action. It should not be changed
// except in unit tests, when performing dependency injection.
var logoutFunc = logout

// HandleLogout parses the request object from AWS APIGateway and passes it to the logout action. The
// request must contain a valid `Cookie` header. If the request succeeds, the response will have a 200
// status, and the body will be empty. If the request fails, the response will have either a 400 or a
// 500 status, and the body will have an `error` field detailing what went wrong. This function returns
// a non-nil error only if JSON marshaling of the response body fails.
func HandleLogout(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Get request parameters
	cookie := auth.ExtractCookie(request.Headers["Cookie"])

	// Perform the action
	err := logoutFunc(cookie, auth.VerifyCookie, dao.Dynamo)

	// Handle the output
	errString, status := errors.UserDetails(err)
	fmt.Println(errors.StackTrace(err))
	json, err := json.Marshal(&logoutResponse{Error: errString})
	if err != nil {
		fmt.Println(errors.StackTrace(errors.Wrap(err, "Failed to marshal Logout response")))
	}
	return events.APIGatewayProxyResponse{Body: string(json), StatusCode: status}, err
}
