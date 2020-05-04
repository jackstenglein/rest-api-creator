package getuser

import (
	"github.com/jackstenglein/rest_api_creator/backend-sls/auth"
	"github.com/jackstenglein/rest_api_creator/backend-sls/dao"
	"github.com/jackstenglein/rest_api_creator/backend-sls/errors"
)

// getUserDatabase wraps the database methods required to perform the getUser action.
// This allows for dependency injection of the database.
type getUserDatabase interface {
	auth.UserGetter
	GetUser(string) (*dao.User, error)
}

// verifyCookieFunc wraps the function type used to check the validity of the user's cookie.
// This allows for dependency injection of the function.
type verifyCookieFunc func(string, auth.UserGetter) (string, error)

func getUser(cookie string, verifyCookie verifyCookieFunc, db getUserDatabase) (*dao.User, error) {
	if cookie == "" {
		return nil, errors.NewClient("Not authenticated")
	}

	email, err := verifyCookie(cookie, db)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to verify cookie")
	}

	user, err := db.GetUser(email)
	return user, errors.Wrap(err, "Failed to get user from database")
}
