import React from 'react';
import Button from 'react-bootstrap/Button';

// ProjectDetails returns the JSX for the project details page.
const ProjectDetails = props => (
  <div>
    <h2>Default Project</h2>
    <p>This is the default project. Future versions of this website will allow you to create multiple projects. If you had created your own project,
      a description of it would go here. For now, this site allows you to define database objects and auto-generate the code to support a REST API
      based on those objects. Future versions will also allow you to automatically deploy the API to AWS and other cloud services. This page will 
      allow you to see the status of those deployed APIs.
    </p>
    <Button variant="primary" onClick={console.log("Download clicked")}>Download Code</Button>
  </div>
)

export default ProjectDetails;
