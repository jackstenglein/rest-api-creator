package getdownload

import (
	"github.com/jackstenglein/rest_api_creator/backend-sls/auth"
	"github.com/jackstenglein/rest_api_creator/backend-sls/codegen/s3"
	"github.com/jackstenglein/rest_api_creator/backend-sls/codegen/sails"
	"github.com/jackstenglein/rest_api_creator/backend-sls/codegen/zip"
	"github.com/jackstenglein/rest_api_creator/backend-sls/dao"
	"github.com/jackstenglein/rest_api_creator/backend-sls/errors"
)

// generateCodeDatabase wraps the database methods required to perform the generateCode
// action. This interface is used to perform dependency injection in unit tests.
type generateCodeDatabase interface {
	auth.UserGetter
	GetProject(string, string) (*dao.Project, error)
}

// cookieVerifier wraps the function type used to check the validity of the user's cookie.
// This allows for dependency injection of the function.
type cookieVerifier func(string, auth.UserGetter) (string, error)

// These variables wrap the different functions that generateCode relies upon. They should
// not be changed except for dependency injection within unit tests.
var download = s3.Download
var unzip = zip.Unzip
var generate = sails.Generate
var zipper = zip.Zip
var upload = s3.Upload
var presign = s3.Presign

// generateCode performs the following steps:
//		1. Download the blank project template from S3
//		2. Unzip the template
//		3. Generate the code for the given project
//		4. Zip the generated code
//		5. Upload the generated zip to S3
// 		6. Generate a pre-signed URL to download the generated zip from S3
// The pre-signed URL is returned, or an empty string if an error occurred.
func generateCode(projectID string, cookie string, verifyCookie cookieVerifier, db generateCodeDatabase) (string, error) {
	if cookie == "" {
		return "", errors.NewClient("Not authenticated")
	}
	if projectID == "" {
		return "", errors.NewClient("Parameter `pid` is required")
	}

	email, err := verifyCookie(cookie, db)
	if err != nil {
		return "", errors.Wrap(err, "Failed to verify cookie")
	}

	project, err := db.GetProject(email, projectID)
	if err != nil {
		return "", errors.Wrap(err, "Failed to get project from database")
	}

	// Download the project template from S3
	err = download("/tmp/blank-sails.zip", "templates/sails.zip")
	if err != nil {
		return "", errors.Wrap(err, "Failed to get project template from S3")
	}

	// Unzip the template
	err = unzip("/tmp/blank-sails.zip", "/tmp")
	if err != nil {
		return "", errors.Wrap(err, "Failed to unzip project template")
	}

	// Generate the code
	err = generate(project, "/tmp/blank-sails-project")
	if err != nil {
		return "", errors.Wrap(err, "Failed to generate project code")
	}

	// Zip the generated code
	err = zipper("/tmp/generated-project.zip", "/tmp/blank-sails-project")
	if err != nil {
		return "", errors.Wrap(err, "Failed to zip generated code")
	}

	// Upload generated zip to S3
	err = upload("/tmp/generated-project.zip", email+"/defaultProject.zip")
	if err != nil {
		return "", errors.Wrap(err, "Failed to upload generated zip to S3")
	}

	// Create the pre-signed URL to download the code
	url, err := presign(email + "/defaultProject.zip")
	return url, errors.Wrap(err, "Failed to generate pre-signed URL")
}
