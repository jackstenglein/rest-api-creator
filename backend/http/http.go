// Package http provides methods for working with HTTP requests and responses.
package http

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jackstenglein/rest_api_creator/backend/errors"
)

type apiResponse interface {
	SetError(string)
}

func headers(cookie string) map[string]string {
	if len(cookie) > 0 {
		return map[string]string{
			"Set-Cookie":                       fmt.Sprintf("session=%s;HttpOnly;", cookie),
			"Access-Control-Allow-Origin":      os.Getenv("CORS_ORIGIN"),
			"Access-Control-Allow-Credentials": "true",
		}
	}

	return map[string]string{
		"Access-Control-Allow-Origin":      os.Getenv("CORS_ORIGIN"),
		"Access-Control-Allow-Credentials": "true",
	}
}

// GatewayResponse returns an APIGatewayResponse that contains the JSON representation of the given
// apiResponse in the body. GatewayResponse also adds CORS headers to the APIGatewayResponse and adds
// a Set-Cookie header if the given cookie is not the empty string.
func GatewayResponse(response apiResponse, cookie string, err error) events.APIGatewayProxyResponse {
	if response == nil {
		return events.APIGatewayProxyResponse{Headers: headers(""), StatusCode: 500}
	}

	errString, status := errors.UserDetails(err)
	response.SetError(errString)
	json, err := json.Marshal(response)
	// TODO: handle error from json marshalling

	return events.APIGatewayProxyResponse{
		Body:       string(json),
		Headers:    headers(cookie),
		StatusCode: status,
	}
}
