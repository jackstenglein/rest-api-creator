package sails

import (
	"strings"
	"testing"

	"github.com/jackstenglein/rest_api_creator/backend-sls/dao"
	"github.com/jackstenglein/rest_api_creator/backend-sls/errors"
)

var writeModelTests = []struct {
	name       string
	object     *dao.Object
	wantErr    error
	wantString string
}{
	{
		name:    "NilObject",
		wantErr: errors.NewServer("Object cannot be nil"),
	},
	{
		name: "NilAttribute",
		object: &dao.Object{
			ID:          "testobject",
			Name:        "testObject",
			CodeName:    "TestObject",
			Description: "This is the description of testObject.",
			Attributes:  []*dao.Attribute{nil},
		},
		wantErr: errors.Wrap(errors.NewServer("Attribute cannot be nil"), "Failed to write attribute"),
	},
	{
		name: "InvalidAttributeType",
		object: &dao.Object{
			ID:          "testobject",
			Name:        "testObject",
			CodeName:    "TestObject",
			Description: "This is the description of testObject.",
			Attributes: []*dao.Attribute{
				{
					Name:     "testAttribute",
					CodeName: "testAttribute",
					Type:     "InvalidType",
				},
			},
		},
		wantErr: errors.Wrap(errors.NewServer("Invalid attribute type: InvalidType"), "Failed to write attribute"),
	},
	{
		name: "NoAttributes",
		object: &dao.Object{
			ID:          "testobject",
			Name:        "testObject",
			CodeName:    "TestObject",
			Description: "This is the description of testObject.",
		},
		wantString: "// api/models/TestObject.js\n" +
			"\n" +
			"module.exports = {\n" +
			"\tattributes: {\n" +
			"\t}\n" +
			"}\n",
	},
	{
		name: "SingleAttribute",
		object: &dao.Object{
			ID:          "testobject",
			Name:        "testObject",
			CodeName:    "TestObject",
			Description: "This is the description of testObject.",
			Attributes: []*dao.Attribute{
				{
					Name:     "testAttribute",
					CodeName: "testAttribute",
					Type:     "Text",
				},
			},
		},
		wantString: "// api/models/TestObject.js\n" +
			"\n" +
			"module.exports = {\n" +
			"\tattributes: {\n" +
			"\t\ttestAttribute: {\n" +
			"\t\t\ttype: 'string',\n" +
			"\t\t},\n" +
			"\t}\n" +
			"}\n",
	},
	{
		name: "MultipleAttributes",
		object: &dao.Object{
			ID:          "testobject",
			Name:        "testObject",
			CodeName:    "TestObject",
			Description: "This is the description of testObject.",
			Attributes: []*dao.Attribute{
				{
					Name:     "testAttribute",
					CodeName: "testAttribute",
					Type:     "Text",
				},
				{
					Name:     "testAttributeTwo",
					CodeName: "testAttributeTwo",
					Type:     "Integer",
				},
			},
		},
		wantString: "// api/models/TestObject.js\n" +
			"\n" +
			"module.exports = {\n" +
			"\tattributes: {\n" +
			"\t\ttestAttribute: {\n" +
			"\t\t\ttype: 'string',\n" +
			"\t\t},\n" +
			"\t\ttestAttributeTwo: {\n" +
			"\t\t\ttype: 'number',\n" +
			"\t\t},\n" +
			"\t}\n" +
			"}\n",
	},
}

func TestWriteModel(t *testing.T) {
	for _, test := range writeModelTests {
		t.Run(test.name, func(t *testing.T) {
			// Setup
			builder := &strings.Builder{}

			// Execute
			err := writeModel(builder, test.object)
			gotString := builder.String()

			// Verify
			if test.wantErr == nil && gotString != test.wantString {
				t.Errorf("Got string: \n %s\n want: %s", gotString, test.wantString)
			}
			if !errors.Equal(err, test.wantErr) {
				t.Errorf("Got err %v; want %v", err, test.wantErr)
			}
		})
	}
}
