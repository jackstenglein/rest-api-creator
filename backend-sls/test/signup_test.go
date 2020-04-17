package test

import (
	gomock "github.com/golang/mock/gomock"
	"github.com/rest_api_creator/backend-sls/actions"
	"github.com/rest_api_creator/backend-sls/errors"
	"github.com/rest_api_creator/backend-sls/mock"
	"testing"
)

var signupTests = []struct {
	// Test name and input
	name     string
	email    string
	password string

	// Mock data
	token               string
	generateTokenCalls  int
	createUserErr       errors.ApiError
	createUserCalls     int
	cookie              string
	generateCookieCalls int

	// Output
	wantStatus int
	wantError  string
	wantCookie string
}{
	{
		name:       "EmptyEmail",
		password:   "12345678",
		wantStatus: 400,
		wantError:  "Invalid email",
	},
	{
		name:       "EmptyPassword",
		email:      "test@example.com",
		wantStatus: 400,
		wantError:  "Password is too short",
	},
	{
		name:       "EmailMissingAt",
		email:      "testexample.com",
		password:   "12345678",
		wantStatus: 400,
		wantError:  "Invalid email",
	},
	{
		name:       "EmailMissingDot",
		email:      "test@examplecom",
		password:   "12345678",
		wantStatus: 400,
		wantError:  "Invalid email",
	},
	{
		name:               "EmailInUse",
		email:              "test@example.com",
		password:           "12345678",
		token:              "testToken",
		generateTokenCalls: 1,
		createUserErr:      errors.NewUserError("Email already in use"),
		createUserCalls:    1,
		wantStatus:         400,
		wantError:          "Email already in use",
	},
	{
		name:                "Success",
		email:               "test@example.com",
		password:            "12345678",
		token:               "testToken",
		generateTokenCalls:  1,
		cookie:              "testCookie",
		generateCookieCalls: 1,
		createUserCalls:     1,
		wantStatus:          200,
		wantError:           "",
		wantCookie:          "testCookie",
	},
}

func TestSignup(t *testing.T) {
	for _, test := range signupTests {
		t.Run(test.name, func(t *testing.T) {
			// Create test objects
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			mockDataStore := mock.NewMockDataStore(mockCtrl)
			mockAuth := mock.NewMockAuthenticator(mockCtrl)
			action := actions.NewSignupAction(mockDataStore, mockAuth)

			// Setup test using input
			request := actions.SignupRequest{Email: test.email, Password: test.password}
			mockAuth.EXPECT().GenerateToken().Return(test.token, nil).Times(test.generateTokenCalls)
			mockDataStore.EXPECT().CreateUser(test.email, gomock.Not(test.password), test.token).Return(test.createUserErr).Times(test.createUserCalls)
			mockAuth.EXPECT().GenerateCookie(test.email, test.token).Return(test.cookie, nil).Times(test.generateCookieCalls)

			// Perform the test
			response, cookie, status := action.Signup(request)

			// Verify the results
			if status != test.wantStatus {
				t.Errorf("Status = %d; want %d", status, test.wantStatus)
			}
			if response.Error != test.wantError {
				t.Errorf("Error = %s; want '%s'", response.Error, test.wantError)
			}
			if cookie != test.wantCookie {
				t.Errorf("Cookie = %s; want '%s'", cookie, test.wantCookie)
			}
		})
	}
}
