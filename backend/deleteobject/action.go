package deleteobject

import (
	"github.com/jackstenglein/rest_api_creator/backend/auth"
	"github.com/jackstenglein/rest_api_creator/backend/errors"
)

// deleteObjectDatabase wraps the database methods required to perform the deleteObject action.
// This allows for dependency injection of the database.
type deleteObjectDatabase interface {
	auth.UserGetter
	DeleteObject(string, string, string) error
}

// deleteObject deletes the given objectID from the given projectID. If the object or project does not exist, no error is returned.
func deleteObject(cookie string, projectID string, objectID string, verifyCookie auth.VerifyCookieFunc, db deleteObjectDatabase) error {
	if projectID == "" || objectID == "" {
		return errors.NewClient("Parameters `projectID` and `objectID` are both required")
	}

	email, err := verifyCookie(cookie, db)
	if err != nil {
		return errors.NewClient("Not authenticated")
	}

	err = db.DeleteObject(email, projectID, objectID)
	return errors.Wrap(err, "Failed to delete object in database")
}
