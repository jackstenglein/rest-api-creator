
// validateAttributes checks an array of attributes and returns [true, errors] if the attributes are valid and
// [false, errors] otherwise. If the attributes are valid, errors is an array of empty error objects.
function validateAttributes(attributes) {
  if (attributes === undefined) {
    return [true, []]
  }

  const errors = [];
  let names = new Set([]);
  var attributesOk = true;
  attributes.forEach(attribute => {
    let [nameOk, nameError] = validateAttributeName(attribute.name);
    let [typeOk, typeError] = validateAttributeType(attribute.type);
    let [defaultOk, defaultError] = validateAttributeDefaultValue(attribute.type, attribute.defaultValue);
    const error = {
      name: nameError, 
      type: typeError,
      defaultValue: defaultError
    };
    if (names.has(attribute.name.toLowerCase())) {
      error.name = `Attribute ${attribute.name} already exists.`;
      attributesOk = false;
    }

    attributesOk = attributesOk && nameOk && typeOk && defaultOk;
    errors.push(error);
    names.add(attribute.name.toLowerCase());
  });
  return [attributesOk, errors];
}

// validateAttributeDefaultValue checks that the defaultValue of an attribute matches its type.
// It returns [true, ""] if the defaultValue is valid and [false, errorMessage] otherwise.
function validateAttributeDefaultValue(type, value) {
  if (type === 'Integer' && value && value.length > 0) {
    const defaultValue = Number(value.replace(/,/g, ''));
    if (Number.isNaN(defaultValue) || !Number.isInteger(defaultValue)) {
      return [false, "Please specify an integer."];
    }
  }
  return [true, ""];
}

// validateAttributeName returns [true, ""] if the name is valid and [false, errorMessage] otherwise.
function validateAttributeName(name) {
  if (name.length === 0) {
    return [false, "Name is a required field."];
  } else if (name.match(/[^a-z]/gi) !== null) {
    return [false, "Only a-z and A-Z are allowed for this field."];
  }
  return [true, ""];
}

// validateAttributeType returns [true, ""] if the attribute type is valid and [false, errorMessage] otherwise.
function validateAttributeType(type) {
  if (type.length === 0 || type === "Choose...") {
    return [false, "Type is a required field."];
  }
  return [true, ""];
}

// validateObjectName returns [true, ""] if the name is valid and [false, errorMessage] otherwise.
function validateObjectName(name, originalId, allObjects) {
  if (name.length === 0) {
    return [false, "Name is a required field."];
  } else if (name.match(/[^a-z]/gi) !== null) {
    return [false, "Only a-z and A-Z are allowed for this field."];
  }

  const newId = name.toLowerCase();
  if (newId === originalId) {
    // In this case, we are editing an existing object and have not changed its name. Therefore, its id should still be unique.
    return [true, ""]; 
  }

  const conflictingObject = allObjects ? allObjects[newId] : undefined;
  if (conflictingObject !== undefined) {
    // In this case, we are creating a new object or editing an existing object and changed its name such that
    // its id conflicts with another object.
    return [false, `An object named '${conflictingObject.name}' already exists.`];
  }

  return [true, ""];
}

// validateObject checks the provided object and returns a corresponding object of error messages. It returns
// [true, errors] if the object is valid and [false, errors] otherwise.
function validateObject(object, allObjects) {
  if (object === undefined) {
    return [true, {}]
  }

  let [nameOk, nameError] = validateObjectName(object.name, object.id, allObjects);
  let [attributesOk, attributeErrors] = validateAttributes(object.attributes);
  const errors = {
    name: nameError, 
    attributes: attributeErrors
  };
  return [nameOk && attributesOk, errors];
}

export default validateObject;
