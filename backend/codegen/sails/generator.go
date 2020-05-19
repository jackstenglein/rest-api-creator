package sails

import (
	"os"

	"github.com/jackstenglein/rest_api_creator/backend/dao"
	"github.com/jackstenglein/rest_api_creator/backend/errors"
)

func generateModel(object *dao.Object, rootDir string) error {
	file, err := os.Create(rootDir + "/api/models/" + object.CodeName + ".js")
	if err != nil {
		return errors.Wrap(err, "Failed to create model file for object: "+object.CodeName)
	}
	defer file.Close()

	err = writeModel(file, object)
	return errors.Wrap(err, "Failed to write model file for object: "+object.CodeName)
}

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

	// TODO: generate endpoints

	return nil
}
