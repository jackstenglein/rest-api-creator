package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jackstenglein/rest_api_creator/backend-sls/actions"
	"github.com/jackstenglein/rest_api_creator/backend-sls/dao"
	apierrors "github.com/jackstenglein/rest_api_creator/backend-sls/errors"
)

type JsonGetProjectResponse struct {
	Project *dao.Project `json:"project,omitempty"`
	Error   string       `json:"error,omitempty"`
}

func HandleError(err error) (string, int) {
	if err == nil {
		return "", 200
	}

	var uerr *apierrors.UserError
	if errors.As(err, &uerr) {
		return uerr.Error(), uerr.StatusCode()
	}

	fmt.Printf("%+v", err)
	return err.Error(), 500
}

func HandleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	getProjectRequest := actions.GetProjectRequest{}
	getProjectRequest.Id = request.PathParameters["id"]

	cookie := request.Headers["Cookie"]
	fmt.Println("Cookie", cookie)
	startIndex := strings.Index(cookie, "session=") + len("session=")
	stopIndex := strings.Index(cookie, ";HttpOnly")
	getProjectRequest.Cookie = cookie[startIndex:stopIndex]

	action := actions.DefaultGetProjectAction()
	project, err := action.GetProject(getProjectRequest)
	errString, status := HandleError(err)

	json, _ := json.Marshal(JsonGetProjectResponse{Project: project, Error: errString})
	return events.APIGatewayProxyResponse{Body: string(json), StatusCode: status}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
