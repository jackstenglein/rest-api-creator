package test

import (
	gomock "github.com/golang/mock/gomock"
	"github.com/rest_api_creator/backend-sls/actions"
	"github.com/rest_api_creator/backend-sls/dao"
	"github.com/rest_api_creator/backend-sls/errors"
	"github.com/rest_api_creator/backend-sls/mock"
	"testing"
)

var loginTests = []struct {
	// test name/input
	name     string
	email    string
	password string

	// mock data
	getUserErr          errors.ApiError
	user                dao.User
	getUserCalls        int
	generateTokenErr    error
	token               string
	generateTokenCalls  int
	generateCookieErr   error
	cookie              string
	generateCookieCalls int
	updateUserErr       error
	updateUserCalls     int

	// output
	wantStatus int
	wantError  string
	wantCookie string
}{
	{
		name:       "EmptyEmail",
		email:      "",
		password:   "12345678",
		wantStatus: 400,
		wantError:  "Email and password are required",
		wantCookie: "",
	},
	{
		name:         "EmptyPassword",
		email:        "test@example.com",
		password:     "",
		getUserErr:   nil,
		user:         dao.User{},
		getUserCalls: 0,
		wantStatus:   400,
		wantError:    "Email and password are required",
		wantCookie:   "",
	},
	{
		name:         "GetUserServerError",
		email:        "test@example.com",
		password:     "12345678",
		getUserErr:   errors.NewServerError("Failed to get user"),
		user:         dao.User{},
		getUserCalls: 1,
		wantStatus:   500,
		wantError:    "Failed to get user",
		wantCookie:   "",
	},
	{
		name:         "IncorrectPassword",
		email:        "test@example.com",
		password:     "incorrect",
		getUserErr:   nil,
		user:         dao.User{Email: "test@example.com", Password: "$2a$14$MNkzNEv8Su7mHfLPIdWoU.t5lElbvlnDka11w27zgfy6Sw44zZsku"},
		getUserCalls: 1,
		wantStatus:   400,
		wantError:    "Incorrect email or password",
		wantCookie:   "",
	},
	{
		name:               "GenerateTokenFailure",
		email:              "test@example.com",
		password:           "12345678",
		user:               dao.User{Email: "test@example.com", Password: "$2a$14$MNkzNEv8Su7mHfLPIdWoU.t5lElbvlnDka11w27zgfy6Sw44zZsku"},
		getUserCalls:       1,
		generateTokenErr:   errors.NewServerError("Could not read random bytes"),
		generateTokenCalls: 1,
		wantStatus:         500,
		wantError:          "Failed to create auth token",
		wantCookie:         "",
	},
	{
		name:                "GenerateCookieFailure",
		email:               "test@example.com",
		password:            "12345678",
		user:                dao.User{Email: "test@example.com", Password: "$2a$14$MNkzNEv8Su7mHfLPIdWoU.t5lElbvlnDka11w27zgfy6Sw44zZsku"},
		getUserCalls:        1,
		token:               "testToken",
		generateTokenCalls:  1,
		generateCookieErr:   errors.NewServerError("Failed to write to HMAC struct"),
		generateCookieCalls: 1,
		wantStatus:          500,
		wantError:           "Failed to create cookie",
		wantCookie:          "",
	},
	{
		name:                "UpdateUserTokenFailure",
		email:               "test@example.com",
		password:            "12345678",
		user:                dao.User{Email: "test@example.com", Password: "$2a$14$MNkzNEv8Su7mHfLPIdWoU.t5lElbvlnDka11w27zgfy6Sw44zZsku"},
		getUserCalls:        1,
		token:               "testToken",
		generateTokenCalls:  1,
		generateCookieCalls: 1,
		updateUserErr:       errors.NewServerError("Failed to update item"),
		updateUserCalls:     1,
		wantStatus:          500,
		wantError:           "Failed to update item",
		wantCookie:          "",
	},
	{
		name:                "Success",
		email:               "test@example.com",
		password:            "12345678",
		user:                dao.User{Email: "test@example.com", Password: "$2a$14$MNkzNEv8Su7mHfLPIdWoU.t5lElbvlnDka11w27zgfy6Sw44zZsku"},
		getUserCalls:        1,
		token:               "testToken",
		generateTokenCalls:  1,
		cookie:              "testCookie",
		generateCookieCalls: 1,
		updateUserCalls:     1,
		wantStatus:          200,
		wantError:           "",
		wantCookie:          "testCookie",
	},
}

func TestLogin(t *testing.T) {
	for _, test := range loginTests {
		t.Run(test.name, func(t *testing.T) {
			// Create test objects
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			mockDataStore := mock.NewMockDataStore(mockCtrl)
			mockAuth := mock.NewMockAuthenticator(mockCtrl)
			action := actions.NewLoginAction(mockDataStore, mockAuth)

			// Setup test using input
			request := actions.LoginRequest{Email: test.email, Password: test.password}
			mockDataStore.EXPECT().GetUser(test.email).Return(test.user, test.getUserErr).Times(test.getUserCalls)
			mockAuth.EXPECT().GenerateToken().Return(test.token, test.generateTokenErr).Times(test.generateTokenCalls)
			mockAuth.EXPECT().GenerateCookie(test.email, test.token).Return(test.cookie, test.generateCookieErr).Times(test.generateCookieCalls)
			mockDataStore.EXPECT().UpdateUserToken(test.email, test.token).Return(test.updateUserErr).Times(test.updateUserCalls)

			// Perform the action
			response, cookie, status := action.Login(request)

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
