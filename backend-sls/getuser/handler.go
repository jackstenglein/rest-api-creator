// Package getuser handles requests to the GET /user REST API endpoint.
package getuser

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jackstenglein/rest_api_creator/backend-sls/auth"
	"github.com/jackstenglein/rest_api_creator/backend-sls/dao"
	"github.com/jackstenglein/rest_api_creator/backend-sls/errors"
)

// getUserResponse contains the fields returned in the API JSON response body.
type getUserResponse struct {
	User  *dao.User `json:"user,omitempty"`
	Error string    `json:"error,omitempty"`
}

// getUserFunc points to the function used to perform the getUser action. It should only be
// changed in unit tests.
var getUserFunc = getUser

// HandleGetUser parses the request object from AWS APIGateway and passes it to the getUser action.
// The request must contain a valid `Cookie` header. If the request succeeds, the response will have
// a 200 status, and the body will contain the user object. If the request fails, the response will
// have either a 400 or 500 status, and the body will have an `error` field detailing what went wrong.
// This function always returns a nil error.
func HandleGetUser(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	cookie := auth.ExtractCookie(request.Headers["Cookie"])
	user, err := getUserFunc(cookie, auth.VerifyCookie, dao.Dynamo)

	errString, status := errors.UserDetails(err)
	fmt.Println(errors.StackTrace(err))
	json, err := json.Marshal(&getUserResponse{User: user, Error: errString})
	if err != nil {
		fmt.Println(errors.StackTrace(errors.Wrap(err, "Failed to marshal GetUser response")))
	}

	return events.APIGatewayProxyResponse{
		Body:       string(json),
		Headers:    map[string]string{"Access-Control-Allow-Origin": "http://jackstenglein-rest-api-creator.s3-website-us-east-1.amazonaws.com", "Access-Control-Allow-Credentials": "true"},
		StatusCode: status,
	}, nil
}
