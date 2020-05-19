package sails

import (
	"io"

	"github.com/jackstenglein/rest_api_creator/backend/dao"
	"github.com/jackstenglein/rest_api_creator/backend/errors"
)

// writeAttribute writes the given attribute to the given writer. It should only
// be called when the writer is in the middle of writing a model definition.
func writeAttribute(writer io.Writer, attribute *dao.Attribute) error {
	if attribute == nil {
		return errors.NewServer("Attribute cannot be nil")
	}

	_, err := io.WriteString(writer, "\t\t")
	if err != nil {
		return errors.Wrap(err, "Failed to write string")
	}
	_, err = io.WriteString(writer, attribute.CodeName)
	if err != nil {
		return errors.Wrap(err, "Failed to write string")
	}
	_, err = io.WriteString(writer, ": {\n\t\t\ttype: ")
	if err != nil {
		return errors.Wrap(err, "Failed to write string")
	}

	switch attribute.Type {
	case "Text":
		_, err = io.WriteString(writer, "'string',\n")
	case "Integer":
		_, err = io.WriteString(writer, "'number',\n")
	default:
		return errors.NewServer("Invalid attribute type: " + attribute.Type)
	}

	if err != nil {
		return errors.Wrap(err, "Failed to write string")
	}

	_, err = io.WriteString(writer, "\t\t},\n")
	return errors.Wrap(err, "Failed to write string")
}

// writeModel writes the model definition corresponding to the given object to
// the given writer.
func writeModel(writer io.Writer, object *dao.Object) error {
	if object == nil {
		return errors.NewServer("Object cannot be nil")
	}

	_, err := io.WriteString(writer, "// api/models/")
	if err != nil {
		return errors.Wrap(err, "Failed to write string")
	}
	_, err = io.WriteString(writer, object.CodeName)
	if err != nil {
		return errors.Wrap(err, "Failed to write string")
	}
	_, err = io.WriteString(writer, ".js\n\nmodule.exports = {\n\tattributes: {\n")
	if err != nil {
		return errors.Wrap(err, "Failed to write string")
	}

	for _, attribute := range object.Attributes {
		err = writeAttribute(writer, attribute)
		if err != nil {
			return errors.Wrap(err, "Failed to write attribute")
		}
	}

	_, err = io.WriteString(writer, "\t}\n}\n")
	return errors.Wrap(err, "Failed to write string")
}
