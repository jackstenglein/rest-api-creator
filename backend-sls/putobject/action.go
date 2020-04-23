package putobject

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jackstenglein/rest_api_creator/backend-sls/auth"
	"github.com/jackstenglein/rest_api_creator/backend-sls/dao"
	"github.com/jackstenglein/rest_api_creator/backend-sls/errors"
)

// verifyCookieFunc wraps the function type used to check the validity of the user's cookie.
// This allows for dependency injection of the function.
type verifyCookieFunc func(string, auth.UserGetter) (string, error)

// putObjectDatabase wraps the database methods required to perform the putObject action.
// This allows for dependency injection of the database.
type putObjectDatabase interface {
	auth.UserGetter
	UpdateItem(string, string, interface{}) error
}

// uuidFunc wraps the function type used to create UUIDs.
// This allows for dependency injection of the function.
type uuidFunc func() (uuid.UUID, error)

// putObject either creates or replaces the given object within the given project. If an error occurs,
// it is returned. If the object's `ID` field is empty, a UUID will be generated for the object and the
// object will be created in the project. If the object's `ID` field is not empty and an object with the
// ID value already exists in the project, the existing object will be replaced. If the object's `ID`
// field is not empty but no object with the ID value already exists, then the given object will be created
// as a new object.
func putObject(cookie string, projectID string, object *dao.Object, verifyCookie verifyCookieFunc, db putObjectDatabase, uuid uuidFunc) error {
	if cookie == "" || projectID == "" || object == nil {
		return errors.NewClient("Parameters `cookie`, `projectId` and `object` are required")
	}
	if object.Name == "" {
		return errors.NewClient("Object must have a `name` field")
	}

	if object.ID == "" {
		uuid, err := uuid()
		if err != nil {
			return errors.Wrap(err, "Failed to generate UUID")
		}
		object.ID = uuid.String()
	}

	email, err := verifyCookie(cookie, db)
	if err != nil {
		return errors.NewClient("Not authenticated")
	}

	path := fmt.Sprintf("Projects.%s.Objects.%s", projectID, object.ID)
	err = db.UpdateItem(email, path, object)
	return errors.Wrap(err, "Failed database call to put object")
}
