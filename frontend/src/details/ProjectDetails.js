import React from 'react';
import Button from 'react-bootstrap/Button';
import Spinner from 'react-bootstrap/Spinner';
import ErrorAlert from '../portal/ErrorAlert';
import { STATUS_PENDING } from '../redux/modules/network.js';

const DownloadButton = props => {
  if (props.downloadStatus === undefined || props.downloadStatus === "") {
    return (<Button variant="primary" onClick={props.download}>Download Code</Button>);
  }

  return (
    <Button variant="primary" className="mr-2" disabled>
      <Spinner as="span" animation="border" size="sm" role="status" aria-hidden="true"/>
      {` ${props.downloadStatus}`}
    </Button>
  );
}

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
      <ErrorAlert error={props.downloadError}/>
      <DownloadButton {...props}/>
    </div>
  )
}

export default ProjectDetails;
