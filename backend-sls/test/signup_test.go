package test

import (
	gomock "github.com/golang/mock/gomock"
	"testing"
	"github.com/rest_api_creator/backend-sls/actions"
	"github.com/rest_api_creator/backend-sls/errors"
	"github.com/rest_api_creator/backend-sls/mock"
)

func testEmptyEmail(t *testing.T, action *actions.SignupAction) {
	request := actions.SignupRequest{Email: "", Password: "12345678"}
	response, status := action.Signup(request)
	if status != 400 {
		t.Errorf("Status = %d; want 400", status)
	}
	if response.Error != "Invalid email" {
		t.Errorf("Error = %s; want 'Invalid email'", response.Error)
	}
}

func testEmptyPassword(t *testing.T, action *actions.SignupAction) {
	request := actions.SignupRequest{Email: "test@example.com", Password: ""}
	response, status := action.Signup(request)
	if status != 400 {
		t.Errorf("Status = %d; want 400", status)
	}
	if response.Error != "Password is too short" {
		t.Errorf("Error = %s; want 'Password is too short'", response.Error)
	}
}

func testEmailMissingAt(t *testing.T, action *actions.SignupAction) {
	request := actions.SignupRequest{Email: "testexample.com", Password: "12345678"}
	response, status := action.Signup(request)
	if status != 400 {
		t.Errorf("Status = %d; want 400", status)
	}
	if response.Error != "Invalid email" {
		t.Errorf("Error = %s; want 'Invalid email'", response.Error)
	}
}

func testEmailMissingDot(t *testing.T, action *actions.SignupAction) {
	request := actions.SignupRequest{Email: "test@examplecom", Password: "12345678"}
	response, status := action.Signup(request)
	if status != 400 {
		t.Errorf("Status = %d; want 400", status)
	}
	if response.Error != "Invalid email" {
		t.Errorf("Error = %s; want 'Invalid email'", response.Error)
	}
}

func testEmailInUse(t *testing.T, action *actions.SignupAction, mockDataStore *mock.MockDataStore) {
	request := actions.SignupRequest{Email: "test@example.com", Password: "12345678"}
	mockDataStore.EXPECT().CreateUser("test@example.com", "12345678").Return(errors.NewUserError("Email already in use")).Times(1)
	response, status := action.Signup(request)
	if status != 400 {
		t.Errorf("Status = %d; want 400", status)
	}
	if response.Error != "Email already in use" {
		t.Errorf("Error = %s; want 'Email already in use'", response.Error)
	}
}

func testSuccess(t *testing.T, action *actions.SignupAction, mockDataStore *mock.MockDataStore) {
	request := actions.SignupRequest{Email: "test@example.com", Password: "12345678"}
	mockDataStore.EXPECT().CreateUser("test@example.com", "12345678").Return(nil).Times(1)
	response, status := action.Signup(request)
	if status != 200 {
		t.Errorf("Status = %d; want 200", status)
	}
	if response.Error != "" {
		t.Errorf("Error = %s; want ''", response.Error)
	}
}

func TestSignup(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockDataStore := mock.NewMockDataStore(mockCtrl)

	action := actions.NewSignupAction(mockDataStore)
	t.Run("EmptyEmail", func(t *testing.T) {testEmptyEmail(t, action)})
	t.Run("EmptyPassword", func(t *testing.T) {testEmptyPassword(t, action)})
	t.Run("EmailMissingAt", func(t *testing.T) {testEmailMissingAt(t, action)})
	t.Run("EmailMissingDot", func(t *testing.T) {testEmailMissingDot(t, action)})
	t.Run("EmailInUse", func(t *testing.T) {testEmailInUse(t, action, mockDataStore)})
	t.Run("Success", func(t *testing.T) {testSuccess(t, action, mockDataStore)})
}
