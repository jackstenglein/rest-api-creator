// Package getproject handles requests to the GET /projects/{pid} API endpoint.
package getproject

import (
	// "encoding/json"
	// "fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jackstenglein/rest_api_creator/backend-sls/auth"
	"github.com/jackstenglein/rest_api_creator/backend-sls/dao"
	// "github.com/jackstenglein/rest_api_creator/backend-sls/errors"
	"github.com/jackstenglein/rest_api_creator/backend-sls/http"
)

type getProjectResponse struct {
	Project *dao.Project `json:"project,omitempty"`
	Error   string       `json:"error,omitempty"`
}

func (response *getProjectResponse) SetError(err string) {
	if response != nil {
		response.Error = err
	}
}

var actionFunc = getProject

// HandleRequest parses the request object from AWS APIGateway and returns a response object containing
// the requested project. The project id must be passed in the `id` path parameter and the request must
// contain a valid `Cookie` header. If the request succeeds, the response will have a 200 status, and the
// body will have a `project` field. If the request fails, the response will have either a 400 or a 500
// status, and the body will have an `error` field detailing what went wrong. This function returns a
// non-nil error only if JSON marshaling of the response body fails.
func HandleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Get request parameters
	projectID := request.PathParameters["pid"]
	cookie := auth.ExtractCookie(request.Headers["Cookie"])

	// Perform the action
	project, err := actionFunc(projectID, cookie)

	// Return the response
	response := &getProjectResponse{Project: project}
	return http.GatewayResponse(response, "", err), nil
}
