package sails

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/jackstenglein/rest_api_creator/backend-sls/dao"
)

var project = &dao.Project{
	Name:        "Default Project",
	Description: "Test project",
	Objects: map[string]*dao.Object{
		"testobject": &dao.Object{
			Name:        "TestObject",
			CodeName:    "TestObject",
			Description: "This is a description of testobject.",
			Attributes: []*dao.Attribute{
				{
					Name:     "testAttribute",
					CodeName: "testAttribute",
					Type:     "Integer",
				},
			},
		},
	},
}

const wantText = `// api/models/TestObject.js

module.exports = {
	attributes: {
		testAttribute: {
			type: 'number',
		},
	}
}
`

func TestGenerate(t *testing.T) {
	// Setup
	wd, err := os.Getwd()
	if err != nil {
		t.Errorf("Error getting working directory: %v", err)
	}
	err = os.Remove(filepath.Join(wd, "testdata", "api", "models", "TestObject.js"))
	if err != nil && !os.IsNotExist(err) {
		t.Errorf("Error removing model file: %v", err)
	}

	// Execute
	err = Generate(project, filepath.Join(wd, "testdata"))
	if err != nil {
		t.Errorf("Got error generating project: %v", err)
	}

	// Verify
	content, err := ioutil.ReadFile(filepath.Join(wd, "testdata", "api", "models", "TestObject.js"))
	if err != nil {
		t.Errorf("Failed to read model file: %v", err)
	}

	text := string(content)
	if text != wantText {
		t.Errorf("Got incorrect text from model file:\n%s", text)
	}
}
