
// validateAttributes checks an array of attributes and returns an array of corresponding error objects.
function validateAttributes(attributes) {
  const errors = [];
  let names = new Set([]);
  attributes.forEach(attribute => {
    const error = {
      name: validateName(attribute.name), 
      type: validateAttributeType(attribute.type),
      defaultValue: validateAttributeDefaultValue(attribute.type, attribute.defaultValue)
    };
    if (names.has(attribute.name.toLowerCase())) {
      error.name = `Attribute ${attribute.name} already exists.`;
    }

    errors.push(error);
    names.add(attribute.name.toLowerCase());
  });
  return errors;
}

// validateAttributeDefaultValue checks that the defaultValue of an attribute matches its type.
// It returns either an empty string or an error message if applicable.
function validateAttributeDefaultValue(type, value) {
  if (type === 'Integer' && value.length > 0) {
    const defaultValue = Number(value.replace(/,/g, ''));
    if (Number.isNaN(defaultValue) || !Number.isInteger(defaultValue)) {
      return "Please specify an integer.";
    }
  }
  return "";
}

// validateAttributeType returns the empty string if the attribute type is valid and an error message otherwise.
function validateAttributeType(type) {
  if (type.length === 0 || type === "Choose...") {
    return "Type is a required field.";
  }
  return "";
}

// validateName returns an empty string if the name is valid and an error message otherwise.
function validateName(name) {
  if (name.length === 0) {
    return "Name is a required field.";
  } else if (name.match(/[^a-z]/gi) !== null) {
    return "Only a-z and A-Z are allowed for this field.";
  }
  return "";
}

// validateObject checks the provided object and returns a corresponding object of error messages.
function validateObject(object) {
  const errors = {
    name: validateName(object.name), 
    attributes: validateAttributes(object.attributes)
  };
  return errors;
}

export default validateObject;
