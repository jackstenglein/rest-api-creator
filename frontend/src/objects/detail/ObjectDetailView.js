import React from 'react';
import ErrorAlert from '../../portal/ErrorAlert';
import { ObjectDocumentation } from '../../documentation/ObjectsDocumentation';
import Button from 'react-bootstrap/Button';
import ButtonToolbar from 'react-bootstrap/ButtonToolbar';
import Col from 'react-bootstrap/Col';
import { Link } from 'react-router-dom';
import Row from 'react-bootstrap/Row';
import EndpointsDocumentation from '../../documentation/EndpointsDocumentation';


const Toolbar = props => (
  <Row className="align-items-center justify-content-between">
    <Col xs="auto">
      <h2>Object Details</h2>
    </Col>
    <Col xs="auto">
      <ButtonToolbar>
        <Link to={`/app/objects/${props.object.id}/edit`}><Button variant="primary">Edit</Button></Link>
        <Button variant="danger" className="ml-2" onClick={props.delete}>Delete</Button>
      </ButtonToolbar>
    </Col>
  </Row>
)

const ObjectDetailView = props => {
  if (props.object === undefined) {
    return null;
  }
  
  return (
    <div>
      <Toolbar delete={props.delete} object={props.object}/>
      <br />
      <ErrorAlert error={props.error} />
      <ObjectDocumentation object={props.object} />
      <EndpointsDocumentation objects={[props.object]} />
    </div>
  );
}

export default ObjectDetailView;
