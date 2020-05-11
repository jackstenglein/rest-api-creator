package getdownload

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jackstenglein/rest_api_creator/backend-sls/auth"
	"github.com/jackstenglein/rest_api_creator/backend-sls/dao"
	"github.com/jackstenglein/rest_api_creator/backend-sls/errors"
)

type getDownloadResponse struct {
	URL   string `json:"url,omitempty"`
	Error string `json:"error,omitempty"`
}

var actionFunc = generateCode

// HandleRequest parses the request object from AWS APIGateway and returns a response object containing a
// URL to download the generated code for the project. The project id must be passed in the `pid` path parameter,
// and the request must contain a valid `Cookie` header. If the request succeeds, the response will have a 200 status,
// and the body will have a `url` field. If the request fails, the response will have either a 400 or a 500 status,
// and the body will have an `error` field.
func HandleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Get request parameters
	projectID := request.PathParameters["pid"]
	cookie := auth.ExtractCookie(request.Headers["Cookie"])

	// Perform the action
	url, err := actionFunc(projectID, cookie, auth.VerifyCookie, dao.Dynamo)

	// Handle the output
	errString, status := errors.UserDetails(err)
	fmt.Println(errors.StackTrace(err))
	json, _ := json.Marshal(&getDownloadResponse{URL: url, Error: errString})
	return events.APIGatewayProxyResponse{
		Body:       string(json),
		Headers:    map[string]string{"Access-Control-Allow-Origin": "http://localhost:3000", "Access-Control-Allow-Credentials": "true"},
		StatusCode: status,
	}, nil
}
