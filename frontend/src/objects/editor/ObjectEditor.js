import React from 'react';
import ObjectAttributeEditor from './ObjectAttributeEditor';
import ObjectInfoEditor from './ObjectInfoEditor';
import Col from 'react-bootstrap/Col';
import Button from 'react-bootstrap/Button';
import Row from 'react-bootstrap/Row';
import ButtonToolbar from 'react-bootstrap/ButtonToolbar';

const EditorToolbar = props => (
  <Row className="align-items-center justify-content-between">
    <Col xs="auto">
      <h2>Edit Object</h2>
    </Col>
    <Col xs="auto">
      <ButtonToolbar>
        <Button className="mr-2" onClick={props.onSave} disabled={!props.isValid}>Save</Button>
        <Button variant="danger" onClick={props.onCancel}>Cancel</Button>
      </ButtonToolbar>
    </Col>
  </Row>
)

const ObjectEditor = props => (
  <div>
    <EditorToolbar isValid={props.isValid} onSave={props.onSave} onCancel={props.onCancel}/>
    <br />
    <ObjectInfoEditor {...props}/>
    <ObjectAttributeEditor 
      values={props.values.attributes} 
      errors={props.errors.attributes} 
      onChange={props.onChangeHandlers.attribute}
      remove={props.removeAttribute}
      add={props.addAttribute}
    />
  </div>
)


export default ObjectEditor;
