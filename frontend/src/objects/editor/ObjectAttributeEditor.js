import React  from 'react';
import Form from 'react-bootstrap/Form';
import Button from 'react-bootstrap/Button';
import Col from 'react-bootstrap/Col';
import Row from 'react-bootstrap/Row';
import { ATTRIBUTE_HELP } from '../../help/Help.js';

const ID_DESCRIPTION = "The id attribute uniquely identifies database records of this"
    + " object type. There can only be one record with a given id. Each time an"
    + " object of this type is created, it is automatically given a new id. The id can"
    + " be used to read and edit specific object instances from the database.";
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
      <Form.Control as="textarea" rows="3" value={ID_DESCRIPTION} disabled />
    </Form.Group>
  </div>
);

// AttributeInput returns the JSX for a form group matching the provided parameters. props should have the following fields:
//    - groupAs -- The component type to use for the form group
//    - name -- The form control name displayed to the end user
//    - as -- The type of form control to use
//    - rows -- The number of rows the form control should span
//    - value -- The value of the user's input
//    - onChange -- The function to call when the input changes
//    - placeholder -- The placeholder to display when there is no input
//    - error -- The error the input has generated
const AttributeInput = props => (
  <Form.Group as={props.groupAs}>
    <Form.Label>{props.name}</Form.Label>
    <Form.Control
      as={props.as}
      rows={props.rows}
      value={props.value}
      onChange={props.onChange}
      placeholder={props.placeholder}
      isInvalid={props.error !== undefined && props.error.length > 0}
    >
      {props.children}
    </Form.Control>
    <Form.Control.Feedback type="invalid">{props.error}</Form.Control.Feedback>
  </Form.Group>
)

// AttributeRow returns a single row of the ObjectAttributeEditor. props should have a `values` field,
// containing the attribute object, and an `errors` field, containing the error object.
const AttributeRow = props => (
  <div>
    <hr />
    <Row className="align-items-center">
      <Col xs="11">
        <Form.Row> 
          <AttributeInput 
            groupAs={Col}
            name="Attribute name"
            placeholder="Enter name"
            value={props.values.name} 
            onChange={props.onChange.bind(null, props.index, "name")}
            error={props.errors ? props.errors.name : undefined}
          />

          <AttributeInput 
            groupAs={Col}
            name="Data type"
            as="select"
            value={props.values.type}
            onChange={props.onChange.bind(null, props.index, "type")}
            error={props.errors ? props.errors.type : undefined}
          >
            <option>Choose...</option>
            <option>Text</option>
            <option>Integer</option>
          </AttributeInput>

          <AttributeInput
            groupAs={Col}
            name="Default Value"
            placeholder="Optional value"
            value={props.values.defaultValue}
            onChange={props.onChange.bind(null, props.index, "defaultValue")}
            error={props.errors ? props.errors.defaultValue : undefined}
          />
        </Form.Row>

        <AttributeInput
          name="Description"
          as="textarea"
          rows="2"
          value={props.values.description}
          onChange={props.onChange.bind(null, props.index, "description")}
          placeholder={ATTR_DESC_PLACEHOLDER}
        />
      </Col>
      <Col xs="1" className="d-flex justify-content-center">
        <button type="button" className="close text-danger opacity-5" aria-label="Close" onClick={props.remove.bind(null, props.index)}>
          <span aria-hidden="true">&times;</span>
        </button>
      </Col>
    </Row>
  </div>
)

// ObjectAttributeEditor returns the attribute section of the ObjectEditor view. props should have
// a `values` field, containing an array of attribute objects, and an `errors` field, containing an array 
// of error objects.
const ObjectAttributeEditor = props => (
  <div>
    <h4>
      Attributes <i className="material-icons clickable" onClick={() => props.toggleHelp(ATTRIBUTE_HELP)}>help_outline</i>
    </h4>
    <IdAttribute />
    { props.values && props.values.map((value, idx) => (
      <AttributeRow key={idx} index={idx} values={value} errors={props.errors[idx]} onChange={props.onChange} remove={props.remove}/>
    )) }
    <Button variant="primary" className="mb-2" onClick={props.add}>Add attribute</Button>
  </div>
)

export default ObjectAttributeEditor;
