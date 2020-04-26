import React  from 'react';
import Form from 'react-bootstrap/Form';
import Col from 'react-bootstrap/Col';

const ID_DESCRIPTION = "The id attribute uniquely identifies database records of this"
    + " object type. There can only be one record with a given id. Each time an"
    + " object of this type is created, it is automatically given a new id. The id can"
    + " be used to read specific object instances from the database.";
const ATTR_DESC_PLACEHOLDER = "Add a description of what this attribute represents"
    + " and it will be included in your auto-generated documentation.";

// IdAttribute returns a disabled attribute row with hard-coded values referring to the id attribute.
const IdAttribute = () => (
  <div>
    <Form.Row>
      <Form.Group as={Col}>
          <Form.Label>Attribute name</Form.Label>
          <Form.Control value="id" disabled />
      </Form.Group>
        <Form.Group as={Col}>
            <Form.Label>Data type</Form.Label>
            <Form.Control value="Integer" disabled />
        </Form.Group>
        <Form.Group as={Col}>
            <Form.Label>Default value</Form.Label>
            <Form.Control value="Auto-Increment" disabled />
        </Form.Group>
    </Form.Row>
    <Form.Group>
        <Form.Label>Description</Form.Label>
        <Form.Control as="textarea" rows="2" value={ID_DESCRIPTION} disabled />
    </Form.Group>
  </div>
);

// AttributeRow returns a single row of the ObjectAttributeEditor. props should have a `values` field,
// containing the attribute object, and an `errors` field, containing the error object.
const AttributeRow = props => (
  <div>
    <hr />
    <Form.Row> 
      <Form.Group as={Col}>
        <Form.Label>Attribute name</Form.Label>
        <Form.Control
          value={props.values.name}
          // onChange={(event) => onChange(index, {"name": event.target.value})}
          placeholder="Enter name"
          isInvalid={props.errors.name !== undefined && props.errors.name.length > 0}
        />
        <Form.Control.Feedback type="invalid">{props.errors.name}</Form.Control.Feedback>
      </Form.Group>

      <Form.Group as={Col}>
        <Form.Label>Data type</Form.Label>
        <Form.Control
          as="select"
          value={props.values.type}
          // onChange={(event) => onChange(index, {"type": event.target.value})}
          isInvalid={props.errors.type !== undefined && props.errors.type.length > 0}
        >
          <option>Choose...</option>
          <option>Text</option>
          <option>Integer</option>
        </Form.Control>
        <Form.Control.Feedback type="invalid">{props.errors.type}</Form.Control.Feedback>
      </Form.Group>

      <Form.Group as={Col}>
        <Form.Label>Default value</Form.Label>
        <Form.Control
            value={props.values.defaultValue}
            // onChange={(event) => onChange(index, {"default": event.target.value})}
            placeholder="Optional value"
            isInvalid={props.errors.defaultValue !== undefined && props.errors.defaultValue.length > 0}
        />
        <Form.Control.Feedback type="invalid">{props.errors.defaultValue}</Form.Control.Feedback>
      </Form.Group>
    </Form.Row>

    <Form.Group>
      <Form.Label>Description</Form.Label>
      <Form.Control
          as="textarea"
          rows="2"
          value={props.values.description}
          // onChange={(event) => onChange(index, {"description": event.target.value})}
          placeholder={ATTR_DESC_PLACEHOLDER}
      >
      </Form.Control>
    </Form.Group>
  </div>
)

// ObjectAttributeEditor returns the attribute section of the ObjectEditor view. props should have
// a `values` field, containing an array of attribute objects, and an `errors` field, containing an array 
// of error objects.
const ObjectAttributeEditor = props => (
  <div>
    <h4>Attributes</h4>
    <IdAttribute />
    { props.values.map((value, idx) => (
      <AttributeRow key={idx} values={value} errors={props.errors[idx]}/>
    )) }
  </div>
)

export default ObjectAttributeEditor;
