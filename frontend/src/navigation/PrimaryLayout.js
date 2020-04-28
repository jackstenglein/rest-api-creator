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

const PrimaryLayout = () => (
  <div className="primary-layout">
    <link
      rel="stylesheet"
      href="https://maxcdn.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css"
      integrity="sha384-ggOyR0iXCbMQv3Xipma34MD+dH/1fQ784/j6cY/iJTQUOhcWr7x9JvoRxT2MZw1T"
      crossorigin="anonymous"
    />

    <main className="height100">
      <Container fluid className="height100">
        <Row className="height100">
          <Col xs="auto" className="border-right p-0 mr-1">
            <ProjectNav />
          </Col>
          <Col>
            <Route path="/app/" exact component={ProjectDetails} />
            <Route path="/app/objects/" component={ObjectLayout} />
            <Route path="/app/endpoints/" component={Endpoints} />
            <Route path="/app/documentation" component={DocumentationContainer} />
          </Col>
        </Row>
      </Container>
    </main>
  </div>
)

export default PrimaryLayout;
