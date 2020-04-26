import React from 'react';
import ProjectDocumentation from './ProjectDocumentation.js';
import ObjectsDocumentation from './ObjectsDocumentation.js';

const Documentation = props => (
  <div>
    <ProjectDocumentation project={props.project}/>
    <ObjectsDocumentation objects={props.project.objects}/>
  </div>
)

export default Documentation;
