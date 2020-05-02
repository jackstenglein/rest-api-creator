import React from 'react';
import Button from 'react-bootstrap/Button';
import ButtonToolbar from 'react-bootstrap/ButtonToolbar';
import Col from 'react-bootstrap/Col';
import ErrorAlert from '../portal/ErrorAlert.js';
import Row from 'react-bootstrap/Row';
import ListGroup from 'react-bootstrap/ListGroup';
import { NavLink } from 'react-router-dom';

const listItemClassName="list-group-item list-group-item-action border-bottom"

const ObjectList = props => (
  <ListGroup variant="flush">
    {
      props.objects && Object.entries(props.objects).map(([id, object]) => 
        <NavLink key={id} to={`/app/objects/${id}/edit`} className={listItemClassName}>{object.name}</NavLink>
      )
    }
  </ListGroup>
)

const Toolbar = props => (
  <Row className="align-items-center justify-content-between">
    <Col xs="auto">
      <h2>Objects</h2>
    </Col>
    <Col xs="auto">
      <ButtonToolbar>
        <Button variant="secondary" className="mr-2" onClick={props.onRefresh}>Refresh</Button>
        <NavLink to="/app/objects/create"><Button variant="primary">Create</Button></NavLink>
      </ButtonToolbar>
    </Col>
  </Row>
)

const ObjectListView = props => (
  <div>
    <Toolbar onRefresh={props.onRefresh}/>
    <br />
    <ErrorAlert error={props.network.error} />
    <ObjectList objects={props.objects} />
  </div>
) 

export default ObjectListView;
