// Package deploy handles requests to the PUT /projects/{pid}/deploy REST API endpoint.
package deploy

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jackstenglein/rest_api_creator/backend/auth"
	"github.com/jackstenglein/rest_api_creator/backend/dao"
	"github.com/jackstenglein/rest_api_creator/backend/deploy/ec2"
	"github.com/jackstenglein/rest_api_creator/backend/http"
	"github.com/jackstenglein/rest_api_creator/backend/log"
)

// deployRequest contains the fields passed in the API JSON request body.
type deployRequest struct {
	URL string `json:"url"`
}

// deployResponse contains the fields returned in the API JSON response body.
type deployResponse struct {
	ID    string `json:"instanceId,omitempty"`
	URL   string `json:"url,omitempty"`
	Error string `json:"error,omitempty"`
}

func (response *deployResponse) SetError(err string) {
	if response != nil {
		response.Error = err
	}
}

// deploy points to the function used to perform the deployProject action. It should
// not be changed except in unit tests.
var deploy = deployProject

// HandleDeploy parses the request object from AWS APIGateway and passes it to the deployProject action.
// The request must contain a valid `Cookie` header and a `pid` path parameter. If the request succeeds
// and the EC2 instance launches within 5 seconds, the response body will have a `url` field. If the request
// succeeds but the EC2 instance is slow to launch, the response body will be empty. If the request fails,
// the response body will have an `error` field.
func HandleDeploy(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Get request parameters
	projectID := request.PathParameters["pid"]
	cookie := auth.ExtractCookie(request.Headers["Cookie"])
	var deployRequest deployRequest
	json.Unmarshal([]byte(request.Body), &deployRequest)

	// Perform the action
	instanceID, url, err := deploy(cookie, projectID, deployRequest, auth.VerifyCookie, dao.Dynamo, ec2.EC2)
	log.Error(err)

	// Return the response
	return http.GatewayResponse(&deployResponse{ID: instanceID, URL: url}, "", err), nil
}
