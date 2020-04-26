import React from 'react';
import Nav from 'react-bootstrap/Nav';
import { NavLink } from 'react-router-dom';

const ProjectNav = () => (
  <Nav variant="pills" className="flex-column">
    <NavLink to="/app/details" eventKey="details" className="nav-link" activeClassName="active"><Nav.Item>Project Details</Nav.Item></NavLink>
    <NavLink to="/app/objects" eventKey="objects" className="nav-link" activeClassName="active"><Nav.Item>Objects</Nav.Item></NavLink>
    <NavLink to="/app/endpoints" eventKey="endpoints" className="nav-link" activeClassName="active"><Nav.Item>Endpoints</Nav.Item></NavLink>
    <NavLink to="/app/documentation" eventKey="documentation" className="nav-link" activeClassName="active"><Nav.Item>Documentation</Nav.Item></NavLink>
  </Nav>
)

export default ProjectNav;
