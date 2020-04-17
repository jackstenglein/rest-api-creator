package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rest_api_creator/backend-sls/actions"
)

func HandleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var signupRequest actions.SignupRequest
	json.Unmarshal([]byte(request.Body), &signupRequest)
	action := actions.DefaultSignupAction()

	response, cookie, status := action.Signup(signupRequest)
	json, _ := json.Marshal(response)
	var headers = map[string]string{
		"Set-Cookie": fmt.Sprintf("session=%s;HttpOnly;SameSite=strict;Secure", cookie),
	}
	return events.APIGatewayProxyResponse{Headers: headers, Body: string(json), StatusCode: status}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
