package logout

import (
	"github.com/jackstenglein/rest_api_creator/backend-sls/auth"
	"github.com/jackstenglein/rest_api_creator/backend-sls/errors"
)

// logoutDatabase wraps the database methods required to perform the logout action.
// This allows for dependency injection of the database.
type logoutDatabase interface {
	auth.UserGetter
	UpdateUserToken(string, string) error
}

// verifyCookieFunc wraps the function type used to check the validity of the user's cookie.
// This allows for dependency injection of the function.
type verifyCookieFunc func(string, auth.UserGetter) (string, error)

// logout removes the auth token from the user associated with the given cookie in the
// given database. It returns the error generated, if the cookie was invalid or the
// database update failed. In this case, the user should still be considered logged in.
func logout(cookie string, verifyCookie verifyCookieFunc, db logoutDatabase) error {
	if cookie == "" {
		return errors.NewClient("Parameter `cookie` is required")
	}

	email, err := verifyCookie(cookie, db)
	if err != nil {
		return errors.NewClient("Not authenticated")
	}

	err = db.UpdateUserToken(email, "")
	return errors.Wrap(err, "Failed to remove auth token")
}
