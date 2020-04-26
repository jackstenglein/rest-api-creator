import React from 'react';
import ObjectAttributeEditor from './ObjectAttributeEditor';
import ObjectInfoEditor from './ObjectInfoEditor';
import Breadcrumbs from '../navigation/Breadcrumbs';
import Col from 'react-bootstrap/Col';
import Button from 'react-bootstrap/Button';
import Row from 'react-bootstrap/Row';
import ButtonToolbar from 'react-bootstrap/ButtonToolbar';

const ObjectEditor = props => (
  <div>
    <Row className="align-items-center justify-content-between">
      <Col xs="auto">
        <h2>Edit Object</h2>
      </Col>
      <Col xs="auto">
        <ButtonToolbar>
          <Button className="mr-2">Save</Button>
          <Button variant="danger">Cancel</Button>
        </ButtonToolbar>
      </Col>
    </Row>
    <br />
    <ObjectInfoEditor values={props.values} errors={props.errors}/>
    <ObjectAttributeEditor values={props.values.attributes} errors={props.errors.attributes}/>
  </div>
)


export default ObjectEditor;
