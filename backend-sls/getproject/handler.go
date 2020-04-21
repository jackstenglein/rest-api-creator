// Package getproject handles requests to the GET /project/{id} API endpoint.
package getproject

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jackstenglein/rest_api_creator/backend-sls/dao"
	"github.com/jackstenglein/rest_api_creator/backend-sls/errors"
)

type getProjectResponse struct {
	Project *dao.Project `json:"project,omitempty"`
	Error   string       `json:"error,omitempty"`
}

func HandleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	projectID := request.PathParameters["id"]

	cookie := request.Headers["Cookie"]
	startIndex := strings.Index(cookie, "session=") + len("session=")
	stopIndex := strings.Index(cookie, ";HttpOnly")
	cookie = cookie[startIndex:stopIndex] // TODO: check bounds of stop and start indices

	project, err := getProject(projectID, cookie)
	errString, status := errors.UserDetails(err)
	fmt.Println(errors.StackTrace(err))

	json, _ := json.Marshal(&getProjectResponse{Project: project, Error: errString}) // TODO: handle this error
	return events.APIGatewayProxyResponse{Body: string(json), StatusCode: status}, nil
}
