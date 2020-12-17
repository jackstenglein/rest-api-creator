package sails

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/jackstenglein/rest_api_creator/backend/dao"
	"github.com/jackstenglein/rest_api_creator/backend/errors"
)

// Generate creates the Sails.js code for the given project. It expects a blank Sails.js project located
// at the path specified by rootDir.
func Generate(project *dao.Project, rootDir string) error {

	// Generate objects
	for _, object := range project.Objects {
		err := generateModel(object, rootDir)
		if err != nil {
			return errors.Wrap(err, "Failed to generate model for object "+object.Name)
		}
	}

	// Change database migration strategy configuration
	err := setMigrationStrategy(rootDir)
	return errors.Wrap(err, "Failed to set migration strategy")

	// TODO: generate endpoints
}

func generateModel(object *dao.Object, rootDir string) error {
	file, err := os.Create(rootDir + "/api/models/" + object.CodeName + ".js")
	if err != nil {
		return errors.Wrap(err, "Failed to create model file for object: "+object.CodeName)
	}
	defer file.Close()

	err = writeModel(file, object)
	return errors.Wrap(err, "Failed to write model file for object: "+object.CodeName)
}

func setMigrationStrategy(rootDir string) error {
	filename := rootDir + "/config/models.js"
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return errors.Wrap(err, "Failed to read config/models.js file")
	}

	const fileMode os.FileMode = 0644
	newContents := strings.Replace(string(contents), "// migrate", "migrate", 1)
	err = ioutil.WriteFile(filename, []byte(newContents), fileMode)
	return errors.Wrap(err, "Failed to write config/models.js file")
}
