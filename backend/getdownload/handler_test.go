package getdownload

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jackstenglein/rest_api_creator/backend/errors"
)

type generateCodeFunc func(string, string, cookieVerifier, generateCodeDatabase) (string, error)

func generateCodeMock(wantProjectID string, wantCookie string, url string, err error) generateCodeFunc {
	return func(gotProjectID string, gotCookie string, _ cookieVerifier, _ generateCodeDatabase) (string, error) {
		if gotProjectID != wantProjectID || gotCookie != wantCookie {
			return "", errors.NewServer("Incorrect parameters passed to mock")
		}
		return url, err
	}
}

func handlerRequest(projectID string, cookie string) events.APIGatewayProxyRequest {
	parameters := map[string]string{
		"pid": projectID,
	}
	headers := map[string]string{
		"Cookie": cookie,
	}
	return events.APIGatewayProxyRequest{PathParameters: parameters, Headers: headers}
}

func handlerResponse(url string, err string, status int) events.APIGatewayProxyResponse {
	json, _ := json.Marshal(&getDownloadResponse{URL: url, Error: err})
	return events.APIGatewayProxyResponse{
		Body:       string(json),
		Headers:    map[string]string{"Access-Control-Allow-Origin": os.Getenv("CORS_ORIGIN"), "Access-Control-Allow-Credentials": "true"},
		StatusCode: status,
	}
}

var handlerTests = []struct {
	name string

	// Input
	request events.APIGatewayProxyRequest

	// Mock data
	mockProjectID string
	mockCookie    string
	mockURL       string
	mockErr       error

	// Expected output
	wantResponse events.APIGatewayProxyResponse
	wantErr      error
}{
	{
		name:          "GetProjectError",
		request:       handlerRequest("default", "session=cookie"),
		mockProjectID: "default",
		mockCookie:    "cookie",
		mockErr:       errors.NewServer("Failed to get project"),
		wantResponse:  handlerResponse("", "Failed to get project", 500),
	},
	{
		name:          "SucccessfulInvocation",
		request:       handlerRequest("default", "session=cookie"),
		mockProjectID: "default",
		mockCookie:    "cookie",
		mockURL:       "presigned-url.com",
		wantResponse:  handlerResponse("presigned-url.com", "", 200),
	},
}

func TestHandleRequest(t *testing.T) {
	for _, test := range handlerTests {
		t.Run("", func(t *testing.T) {
			// Setup
			actionFunc = generateCodeMock(test.mockProjectID, test.mockCookie, test.mockURL, test.mockErr)
			defer func() {
				actionFunc = generateCode
			}()

			// Execute
			response, err := HandleRequest(test.request)

			// Verify
			if !reflect.DeepEqual(response, test.wantResponse) {
				t.Errorf("Got response %v; want %v", response, test.wantResponse)
			}
			if !errors.Equal(err, test.wantErr) {
				t.Errorf("Got error %v; want %v", err, test.wantErr)
			}
		})
	}
}
