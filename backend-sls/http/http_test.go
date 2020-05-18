package http

import (
	"os"
	"reflect"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jackstenglein/rest_api_creator/backend-sls/errors"
)

type testResponse struct {
	Test  string `json:"test,omitempty"`
	Error string `json:"error,omitempty"`
}

func (response *testResponse) SetError(err string) {
	if response != nil {
		response.Error = err
	}
}

var gatewayResponseTests = []struct {
	name     string
	response apiResponse
	cookie   string
	err      error

	wantResponse events.APIGatewayProxyResponse
}{
	{
		name: "NilResponse",
		wantResponse: events.APIGatewayProxyResponse{
			Headers: map[string]string{
				"Access-Control-Allow-Origin":      os.Getenv("CORS_ORIGIN"),
				"Access-Control-Allow-Credentials": "true",
			},
			StatusCode: 500,
		},
	},
	{
		name:     "BodyResponse",
		response: &testResponse{Test: "testValue"},
		wantResponse: events.APIGatewayProxyResponse{
			Body: `{"test":"testValue"}`,
			Headers: map[string]string{
				"Access-Control-Allow-Origin":      os.Getenv("CORS_ORIGIN"),
				"Access-Control-Allow-Credentials": "true",
			},
			StatusCode: 200,
		},
	},
	{
		name:     "SetCookie",
		response: &testResponse{},
		cookie:   "cookievalue",
		wantResponse: events.APIGatewayProxyResponse{
			Body: "{}",
			Headers: map[string]string{
				"Set-Cookie":                       "session=cookievalue;HttpOnly;",
				"Access-Control-Allow-Origin":      os.Getenv("CORS_ORIGIN"),
				"Access-Control-Allow-Credentials": "true",
			},
			StatusCode: 200,
		},
	},
	{
		name:     "ServerError",
		response: &testResponse{},
		err:      errors.NewServer("This is the error message"),
		wantResponse: events.APIGatewayProxyResponse{
			Body: `{"error":"This is the error message"}`,
			Headers: map[string]string{
				"Access-Control-Allow-Origin":      os.Getenv("CORS_ORIGIN"),
				"Access-Control-Allow-Credentials": "true",
			},
			StatusCode: 500,
		},
	},
	{
		name:     "UserError",
		response: &testResponse{},
		err:      errors.NewClient("This is the error message"),
		wantResponse: events.APIGatewayProxyResponse{
			Body: `{"error":"This is the error message"}`,
			Headers: map[string]string{
				"Access-Control-Allow-Origin":      os.Getenv("CORS_ORIGIN"),
				"Access-Control-Allow-Credentials": "true",
			},
			StatusCode: 400,
		},
	},
}

func TestGatewayResponse(t *testing.T) {
	for _, test := range gatewayResponseTests {
		t.Run(test.name, func(t *testing.T) {
			response := GatewayResponse(test.response, test.cookie, test.err)
			if !reflect.DeepEqual(response, test.wantResponse) {
				t.Errorf("Got response %v; want %v", response, test.wantResponse)
			}
		})
	}
}
