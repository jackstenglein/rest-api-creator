package putobject

import (
	"fmt"
	"regexp"
	"strings"

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
	UpdateObject(string, string, *dao.Object) error
}

// uuidFunc wraps the function type used to create UUIDs.
// This allows for dependency injection of the function.
type uuidFunc func() (uuid.UUID, error)

// validAttribute checks that the given attribute has valid values. If the attribute is valid, its
// CodeName field is set. If the attribute is invalid, an error is returned. An attribute is invalid if:
//		- Its Name field has length 0
//		- Its Name field contains non-alphabetical characters
// 		- Its Type field is neither `Text` nor `Integer`
func validAttribute(attribute *dao.Attribute) error {
	isAlpha := regexp.MustCompile(`^[A-Za-z]+$`).MatchString

	if attribute == nil {
		return errors.NewClient("Attribute cannot be nil")
	}

	if len(attribute.Name) == 0 {
		return errors.NewClient("Attribute must have a `name` field")
	}

	if !isAlpha(attribute.Name) {
		return errors.NewClient(fmt.Sprintf("Attribute name `%s` contains non-alphabetical characters", attribute.Name))
	}

	switch attribute.Type {
	case "Text", "Integer":
	default:
		return errors.NewClient(fmt.Sprintf("Attribute type `%s` is not supported", attribute.Type))
	}

	attribute.CodeName = fmt.Sprintf("%s%s", strings.ToLower(attribute.Name[0:1]), attribute.Name[1:])
	return nil
}

// validObject checks that the given object has valid values. If the object is valid, its CodeName field
// is set. If the object is invalid, an error is returned. An obejct is invalid if:
// 		- Its Name field has length 0
// 		- Its Name field contains non-alphabetical characters
// 		- Any of its attributes are invalid
func validObject(object *dao.Object) error {
	isAlpha := regexp.MustCompile(`^[A-Za-z]+$`).MatchString

	if len(object.Name) == 0 {
		return errors.NewClient("Object must have a `name` field")
	}

	if !isAlpha(object.Name) {
		return errors.NewClient(fmt.Sprintf("Object name `%s` contains non-alphabetical characters", object.Name))
	}

	for _, attribute := range object.Attributes {
		err := validAttribute(attribute)
		if err != nil {
			return errors.Wrap(err, "Object must have valid attributes")
		}
	}

	object.CodeName = fmt.Sprintf("%s%s", strings.ToUpper(object.Name[0:1]), object.Name[1:])
	return nil
}

// putObject either creates or replaces the given object within the given project. The object's ID is returned.
// If an error occurs, it is returned. If the object's `ID` field is empty, a UUID will be generated for the object
// and the object will be created in the project. If the object's `ID` field is not empty and an object with the
// ID value already exists in the project, the existing object will be replaced. If the object's `ID`
// field is not empty but no object with the ID value already exists, then the given object will be created
// as a new object.
func putObject(cookie string, projectID string, object *dao.Object, verifyCookie verifyCookieFunc, db putObjectDatabase, uuid uuidFunc) (string, error) {
	if cookie == "" || projectID == "" || object == nil {
		return "", errors.NewClient("Parameters `cookie`, `projectId` and `object` are required")
	}

	err := validObject(object)
	if err != nil {
		return "", errors.Wrap(err, "Object is invalid")
	}

	if object.ID == "" {
		uuid, err := uuid()
		if err != nil {
			return "", errors.Wrap(err, "Failed to generate UUID")
		}
		object.ID = uuid.String()
	}

	email, err := verifyCookie(cookie, db)
	if err != nil {
		return "", errors.NewClient("Not authenticated")
	}

	err = db.UpdateObject(email, projectID, object)
	if err != nil {
		return "", errors.Wrap(err, "Failed database call to put object")
	}
	return object.ID, nil
}
