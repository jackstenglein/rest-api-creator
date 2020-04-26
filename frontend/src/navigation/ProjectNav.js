import React from 'react';
import { NavLink } from 'react-router-dom';
import ListGroup from 'react-bootstrap/ListGroup';

const itemClassName="list-group-item list-group-item-action pl-5 pr-5"

const ProjectNav = () => (
  <ListGroup variant="flush" className="text-center sticky-top">
    <NavLink to="/app/details" className={itemClassName} activeClassName="active">Project Details</NavLink>
    <NavLink to="/app/objects" className={itemClassName} activeClassName="active">Objects</NavLink>
    <NavLink to="/app/endpoints" className={itemClassName} activeClassName="active">Endpoints</NavLink>
    <NavLink to="/app/documentation" className={itemClassName} activeClassName="active">Documentation</NavLink>
  </ListGroup>
)

export default ProjectNav;
