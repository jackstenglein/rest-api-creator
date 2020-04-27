import React  from 'react';
import Form from 'react-bootstrap/Form';

const OBJECT_DESC_PLACEHOLDER = "Add a description of what this object represents"
    + " and it will be included in your auto-generated documentation.";

// ObjectInfoEditor returns the info section of the ObjectEditor view. props should have a `values`
// field, containing the object definition, and an `errors` field, containing the form errors.
const ObjectInfoEditor = (props) => (
  <>
    <h4>Details</h4>
    <Form.Group>
    <Form.Label>Object name</Form.Label>
    <Form.Control
      onChange={props.onChangeHandlers.name}
      value={props.values.name}
      placeholder="Enter name"
      isInvalid={props.errors.name !== undefined && props.errors.name.length > 0}
    />
    <Form.Control.Feedback type="invalid">{props.errors.name}</Form.Control.Feedback>
    </Form.Group>
    <Form.Group>
      <Form.Label>Description</Form.Label>
      <Form.Control
        as="textarea"
        rows="3"
        placeholder={OBJECT_DESC_PLACEHOLDER}
        onChange={props.onChangeHandlers.description}
        value={props.values.description}
      />
    </Form.Group>
  </>
)

export default ObjectInfoEditor;
