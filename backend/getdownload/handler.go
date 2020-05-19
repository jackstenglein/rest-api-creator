package getdownload

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/jackstenglein/rest_api_creator/backend/auth"
	"github.com/jackstenglein/rest_api_creator/backend/dao"
	"github.com/jackstenglein/rest_api_creator/backend/http"
)

type getDownloadResponse struct {
	URL   string `json:"url,omitempty"`
	Error string `json:"error,omitempty"`
}

func (response *getDownloadResponse) SetError(err string) {
	if response != nil {
		response.Error = err
	}
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
	response := &getDownloadResponse{URL: url}
	return http.GatewayResponse(response, "", err), nil
}
