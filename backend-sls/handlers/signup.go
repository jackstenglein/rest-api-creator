package main

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/pkg/errors"
	"github.com/jackstenglein/rest_api_creator/backend-sls/actions"
	apierrors "github.com/jackstenglein/rest_api_creator/backend-sls/errors"
)

type JsonSignupResponse struct {
	Error string `json:"error,omitempty"`
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
	var signupRequest actions.SignupRequest
	json.Unmarshal([]byte(request.Body), &signupRequest)
	action := actions.DefaultSignupAction()

	cookie, err := action.Signup(signupRequest)
	errString, status := HandleError(err)

	json, _ := json.Marshal(JsonSignupResponse{errString})
	var headers = map[string]string{
		"Set-Cookie": fmt.Sprintf("session=%s;HttpOnly;SameSite=strict;Secure", cookie),
	}
	return events.APIGatewayProxyResponse{Headers: headers, Body: string(json), StatusCode: status}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
