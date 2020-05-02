import React from 'react';
import { connect } from 'react-redux';
import Button from 'react-bootstrap/Button';
import ErrorAlert from '../portal/ErrorAlert';
import { STATUS_PENDING } from '../redux/modules/network.js';

// ProjectDetails returns the JSX for the project details page.
const ProjectDetails = props => {
  if (props.project.network.status === STATUS_PENDING) {
    return null;
  }

  if (props.project.network.error) {
    return <ErrorAlert error={props.project.network.error}/>
  }

  return (
    <div>
      <h2>{props.project.name}</h2>
      <p>{props.project.description}</p>
      <Button variant="primary" onClick={() => console.log("Download clicked")}>Download Code</Button>
    </div>
  )
}

export default ProjectDetails;
