import React from 'react';
import { connect } from 'react-redux';
import Button from 'react-bootstrap/Button';

// ProjectDetails returns the JSX for the project details page.
const ProjectDetails = props => (
  <div>
    <h2>{props.project.name}</h2>
    <p>{props.project.description}
    </p>
    <Button variant="primary" onClick={() => console.log("Download clicked")}>Download Code</Button>
  </div>
)

const mapStateToProps = state => {
  return {
    project: state.projects.defaultProject
  }
}

export default connect(mapStateToProps)(ProjectDetails);
