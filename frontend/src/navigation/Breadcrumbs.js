import React from 'react';
import Breadcrumb from 'react-bootstrap/Breadcrumb';

// Breadcrumbs returns a Bootstrap breadcrumb with the given path. props should have a `path` field,
// containing an array of segment objects. Each segment object should have the following fields:
//    - name
//    - link
//    - active
const Breadcrumbs = props => (
  <Breadcrumb bsPrefix="breadcrumb crud">
    { props.path.map(segment => 
      <Breadcrumb.Item href={segment.link} active={segment.active}>{segment.name}</Breadcrumb.Item>
    )}
  </Breadcrumb>
)

export default Breadcrumbs;
