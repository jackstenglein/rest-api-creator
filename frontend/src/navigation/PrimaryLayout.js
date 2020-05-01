import React from 'react';
import { Route } from 'react-router-dom';
import ObjectLayout from '../objects/ObjectLayout.js';
import ProjectNav from './ProjectNav.js';
import Col from 'react-bootstrap/Col';
import Container from 'react-bootstrap/Container';
import Row from 'react-bootstrap/Row';
import DocumentationContainer from '../documentation/DocumentationContainer.js';
import Endpoints from '../endpoints/Endpoints.js';
import ProjectDetails from '../details/ProjectDetails.js';
import { Redirect } from 'react-router-dom';

const PrimaryLayout = () => (
  <div className="primary-layout">
    <main className="height100">
      <Container fluid className="height100">
        <Row className="height100">
          <Col xs="auto" className="border-right p-0 mr-1">
            <ProjectNav />
          </Col>
          <Col className="mt-2">
            <Route path="/app/details" exact component={ProjectDetails} />
            <Route path="/app/objects/" component={ObjectLayout} />
            <Route path="/app/endpoints/" component={Endpoints} />
            <Route path="/app/documentation" component={DocumentationContainer} />
            <Redirect to="/app/details" />
          </Col>
        </Row>
      </Container>
    </main>
  </div>
)

export default PrimaryLayout;
