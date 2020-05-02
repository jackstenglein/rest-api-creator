import React from 'react';
import { Route } from 'react-router-dom';
import ObjectLayout from '../objects/ObjectLayout.js';
import ProjectNav from './ProjectNav.js';
import Col from 'react-bootstrap/Col';
import Container from 'react-bootstrap/Container';
import Row from 'react-bootstrap/Row';
import Documentation from '../documentation/Documentation.js';
import Endpoints from '../endpoints/Endpoints.js';
import ProjectDetails from '../details/ProjectDetails.js';
import { Switch, Redirect } from 'react-router-dom';
import ProjectRoute from './ProjectRoute.js';

const PrimaryLayout = () => (
  <div className="primary-layout">
    <main className="height100">
      <Container fluid className="height100">
        <Row className="height100">
          <Col xs="auto" className="border-right p-0 mr-1">
            <ProjectNav />
          </Col>
          <Col className="mt-2">
            <Switch>
              <Route path="/app/details" exact component={ProjectDetails} />
              <Route path="/app/objects/" component={ObjectLayout} />
              <Route path="/app/endpoints/" component={Endpoints} />
              <ProjectRoute path="/app/documentation" component={Documentation} />
              <Redirect to="/app/details" />
            </Switch>
          </Col>
        </Row>
      </Container>
    </main>
  </div>
)

export default PrimaryLayout;
