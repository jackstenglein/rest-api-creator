package getproject

import (
	"github.com/jackstenglein/rest_api_creator/backend/auth"
	"github.com/jackstenglein/rest_api_creator/backend/dao"
	"github.com/jackstenglein/rest_api_creator/backend/errors"
)

// getProjectDatabase wraps the database methods required to perform the getProject
// action. This interface is used to perform dependency injection in unit tests.
type getProjectDatabase interface {
	auth.UserGetter
	GetProject(string, string) (*dao.Project, error)
}

// verifyCookie points to the function used to check the validity of the cookie.
// This variable should be changed only to perform dependency injection in unit tests.
var verifyCookie = auth.VerifyCookie

// db is the object that implements the required database methods defined in getProjectDatabase.
// This variable should be changed only to perform dependency injection in unit tests.
var db getProjectDatabase = dao.Dynamo

// getProject returns the project with the given id associated with the user email specified
// in cookie. If the specified email does not have a project with the id or another error
// occurs, getProject returns a nil pointer along with the error.
func getProject(id string, cookie string) (*dao.Project, error) {
	if id == "" {
		return nil, errors.NewClient("Parameter `id` is required")
	}

	email, err := verifyCookie(cookie, db)
	if err != nil {
		return nil, errors.NewClient("Not authenticated")
	}

	project, err := db.GetProject(email, id)
	return project, errors.Wrap(err, "Failed to get project")
}
