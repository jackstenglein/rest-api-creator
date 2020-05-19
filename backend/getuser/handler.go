// Package getuser handles requests to the GET /user REST API endpoint.
package getuser

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/jackstenglein/rest_api_creator/backend/auth"
	"github.com/jackstenglein/rest_api_creator/backend/dao"
	"github.com/jackstenglein/rest_api_creator/backend/http"
)

// getUserResponse contains the fields returned in the API JSON response body.
type getUserResponse struct {
	User  *dao.User `json:"user,omitempty"`
	Error string    `json:"error,omitempty"`
}

func (response *getUserResponse) SetError(err string) {
	if response != nil {
		response.Error = err
	}
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
	// Get request parameters
	cookie := auth.ExtractCookie(request.Headers["Cookie"])

	// Get the user
	user, err := getUserFunc(cookie, auth.VerifyCookie, dao.Dynamo)

	// Return the response
	response := &getUserResponse{User: user}
	return http.GatewayResponse(response, "", err), nil
}
