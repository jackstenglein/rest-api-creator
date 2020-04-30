import React from 'react';
import Col from 'react-bootstrap/Col';
import Container from 'react-bootstrap/Container';
import Form from 'react-bootstrap/Form';
import Row from 'react-bootstrap/Row';

const Signup = () => (
  <Container>

    <h2>Signup</h2>

    <Form.Group as={Row}>
      <Form.Label column sm="2">Email</Form.Label>
      <Col sm="10">
          <Form.Control 
            placeholder="email@example.com" 
          />
      </Col>
    </Form.Group>

    <Form.Group as={Row}>
      <Form.Label column sm="2">Password</Form.Label>
      <Col sm="10">
        <Form.Control type="password" placeholder="Password" />
      </Col>
    </Form.Group>

    <Form.Group as={Row}>
      <Form.Label column sm="2">Confirm Password</Form.Label>
      <Col sm="10">
        <Form.Control type="password" placeholder="Password" />
      </Col>
    </Form.Group>
  </Container>
)

export default Signup;
