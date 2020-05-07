import React from 'react';
import Button from 'react-bootstrap/Button';
import ButtonToolbar from 'react-bootstrap/ButtonToolbar';
import Col from 'react-bootstrap/Col';
import ErrorAlert from '../portal/ErrorAlert.js';
import Row from 'react-bootstrap/Row';
import ListGroup from 'react-bootstrap/ListGroup';
import { NavLink } from 'react-router-dom';

const listItemClassName="list-group-item list-group-item-action border-bottom"

const ObjectList = props => {
  if (props.objects && Object.keys(props.objects).length > 0) {
    return (
      <ListGroup variant="flush">
        {
          Object.entries(props.objects).map(([id, object]) => 
            <NavLink key={id} to={`/app/objects/${id}/edit`} className={listItemClassName}>{object.name}</NavLink>
          )
        }
      </ListGroup>
    )
  }

  return <p>This project has no objects yet. Click "Create" to add some.</p>
}

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

const ObjectListView = props => {
  const { network, objects } = props.project;
  return (
  <div>
    <Toolbar onRefresh={props.refreshProject}/>
    <br />
    <ErrorAlert error={network.error} />
    <ObjectList objects={objects} />
  </div>
  )
}

export default ObjectListView;
