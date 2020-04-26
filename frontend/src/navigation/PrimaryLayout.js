import React from 'react';
import { Route } from 'react-router-dom';
import ObjectLayout from '../objects/ObjectLayout.js';
import ProjectNav from './ProjectNav.js';
import Col from 'react-bootstrap/Col';
import Container from 'react-bootstrap/Container';
import Row from 'react-bootstrap/Row';
import Documentation from '../documentation/Documentation.js';
import Endpoints from '../endpoints/Endpoints.js';

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
            <Route path="/app/objects/" component={ObjectLayout} />
            <Route path="/app/endpoints/" component={Endpoints} />
            <Route path="/app/documentation" render={(props) => <Documentation {...props} project={TEST_PROJECT}/>} />
          </Col>
        </Row>
      </Container>
    </main>
  </div>
)

const TEST_PROJECT = {
  name: "Default Project",
  description: "This is a sample project to test static UI. This project will eventually be dynamically pulled from Redux.",
  objects: [
    {
      name: "User",
      description: "User represents the structure of an end user in the database.",
      attributes: [
        {
          name: "AttributeOne",
          type: "Integer",
          required: true,
          description: "AttributeOne is an integer attribute that stores the first attribute of the user."
        },
        {
          name: "AttributeTwo",
          type: "Text",
          required: false,
          description: "AttributeTwo is a text attribute that stores the second attribute of the user."
        }
      ]
    },
    {
      name: "Dog",
      description: "Dog represents the structure of a dog object in the database.",
      attributes: [
        {
          name: "AttributeOne",
          type: "Text",
          required: true,
          description: "AttributeOne is a text attribute that stores the first attribute of the dog. This attribute is required."
        },
        {
          name: "AttributeTwo",
          type: "Text",
          required: true,
          description: "AttributeTwo is a text attribute that stores the second attribute of the dog."
        }
      ]
    }
  ]
}

export default PrimaryLayout;
