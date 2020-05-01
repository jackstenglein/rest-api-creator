// Package putobject handles requests to the PUT /projects/{pid}/objects REST API endpoint.
package putobject

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
	"github.com/jackstenglein/rest_api_creator/backend-sls/auth"
	"github.com/jackstenglein/rest_api_creator/backend-sls/dao"
	"github.com/jackstenglein/rest_api_creator/backend-sls/errors"
)

// putObjectResponse contains the fields returned in the API JSON response body.
type putObjectResponse struct {
	ID    string `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
}

// putObjectFun points to the function used to perform the putObject action. It
// should not be changed except in unit tests, when performing dependency injection.
var putObjectFunc = putObject

// HandlePutObject parses the request object from AWS APIGateway and passes it to the putObject action. The
// request must contain a valid `Cookie` header, a `pid` path parameter, and an object defintion in the body.
// If the request succeeds, the response will have a 200 status, and the body will be empty. If the request
// fails, the response will have either a 400 or a 500 status, and the body will have an `error` field detailing
// what went wrong. This function returns a non-nil error only if JSON marshaling of the response body fails.
func HandlePutObject(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Get request parameters
	cookie := auth.ExtractCookie(request.Headers["Cookie"])
	projectID := request.PathParameters["pid"]
	var object *dao.Object
	json.Unmarshal([]byte(request.Body), &object)

	// Perform the action
	id, err := putObjectFunc(cookie, projectID, object, auth.VerifyCookie, dao.Dynamo, uuid.NewRandom)

	// Handle the output
	errString, status := errors.UserDetails(err)
	fmt.Println(errors.StackTrace(err))
	json, err := json.Marshal(&putObjectResponse{ID: id, Error: errString})
	if err != nil {
		fmt.Println(errors.StackTrace(errors.Wrap(err, "Failed to marshal PutObject response")))
	}
	return events.APIGatewayProxyResponse{
		Body:       string(json),
		Headers:    map[string]string{"Access-Control-Allow-Origin": "http://localhost:3000", "Access-Control-Allow-Credentials": "true"},
		StatusCode: status,
	}, err
}
