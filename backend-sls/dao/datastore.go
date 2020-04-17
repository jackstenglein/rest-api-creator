package dao

import "github.com/rest_api_creator/backend-sls/errors"

type User struct {
	Email    string `dynamodbav:"Email" json:"email"`
	Password string `dynamodbav:"Password" json:"-"`
	Token    string `dynamodbav:"SessionToken" json:"-"`
}

type DataStore interface {
	CreateUser(string, string, string) errors.ApiError
	GetUser(string) (User, errors.ApiError)
	UpdateUserToken(string, string) error
}
