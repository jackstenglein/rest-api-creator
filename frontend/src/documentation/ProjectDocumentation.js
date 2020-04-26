import React from 'react';

// ProjectDocumentation returns the JSX for the project section of the documentation view. props should have a
// `project` field, containing an object with the name and description of the documented project.
const ProjectDocumentation = props => (
  <div>
    <h2>{props.project.name}</h2>
    <p>{props.project.description}</p>
  </div>
)

export default ProjectDocumentation;
