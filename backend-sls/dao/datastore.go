package dao

import "github.com/rest_api_creator/backend-sls/errors"

type DataStore interface {
	CreateUser(string, string) errors.ApiError
}
