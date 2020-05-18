package http

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jackstenglein/rest_api_creator/backend-sls/errors"
)

type apiResponse interface {
	SetError(string)
}

func headers(cookie string) map[string]string {
	if len(cookie) > 0 {
		return map[string]string{
			"Set-Cookie":                       fmt.Sprintf("session=%s;HttpOnly;", cookie),
			"Access-Control-Allow-Origin":      "http://jackstenglein-rest-api-creator.s3-website-us-east-1.amazonaws.com",
			"Access-Control-Allow-Credentials": "true",
		}
	}

	return map[string]string{
		"Access-Control-Allow-Origin":      "http://jackstenglein-rest-api-creator.s3-website-us-east-1.amazonaws.com",
		"Access-Control-Allow-Credentials": "true",
	}
}

// GatewayResponse does stuff
func GatewayResponse(response apiResponse, cookie string, err error) events.APIGatewayProxyResponse {

	errString, status := errors.UserDetails(err)
	response.SetError(errString)
	json, err := json.Marshal(response)
	// if err != nil {
	// 	fmt.Println(errors.StackTrace(errors.Wrap(err, "Failed to marshal GetProject response")))
	// }
	return events.APIGatewayProxyResponse{
		Body:       string(json),
		Headers:    headers(cookie),
		StatusCode: status,
	}
}
