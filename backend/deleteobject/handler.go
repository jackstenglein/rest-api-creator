// Package deleteobject handles requests to the DELETE /projects/{pid}/objects/{oid} REST API endpoint.
package deleteobject

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/jackstenglein/rest_api_creator/backend/auth"
	"github.com/jackstenglein/rest_api_creator/backend/dao"
	"github.com/jackstenglein/rest_api_creator/backend/http"
	"github.com/jackstenglein/rest_api_creator/backend/log"
)

// deleteObjectResponse contains the fields returned in the API JSON response body.
type deleteObjectResponse struct {
	Error string `json:"error,omitempty"`
}

func (response *deleteObjectResponse) SetError(err string) {
	if response != nil {
		response.Error = err
	}
}

// deleteObjectFunc points to the function used to perform the deleteObject action. It
// should not be changed except in unit tests, when performing dependency injection.
var deleteObjectFunc = deleteObject

// HandleDeleteObject parses the request from AWS APIGateway and passes it to the deleteObject action. The
// request must contain a valid `Cookie` header, as well as `pid` and `oid` path parameters. If the request
// succeeds, the response will have a 200 status and an empty body. If the request fails, the response will
// have either a 400 or 500 status and an `error` field in the body.
func HandleDeleteObject(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Get request parameters
	cookie := auth.ExtractCookie(request.Headers["Cookie"])
	projectID := request.PathParameters["pid"]
	objectID := request.PathParameters["oid"]
	log.Info("Got request parameters. Cookie:", cookie, "projectID:", projectID, "objectID:", objectID)

	// Delete the object
	err := deleteObjectFunc(cookie, projectID, objectID, auth.VerifyCookie, dao.Dynamo)
	log.Error(err)

	// Handle the output
	return http.GatewayResponse(&deleteObjectResponse{}, "", err), nil
}
