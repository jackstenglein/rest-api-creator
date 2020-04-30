import React from 'react';
import Button from 'react-bootstrap/Button';
import Col from 'react-bootstrap/Col';
import Container from 'react-bootstrap/Container';
import Form from 'react-bootstrap/Form';
import Row from 'react-bootstrap/Row';
import { Link } from 'react-router-dom';

const Input = props => (
  <Form.Group as={Row}>
    <Form.Label column sm="2">{props.name}</Form.Label>
    <Col sm="10">
        <Form.Control 
          placeholder={props.placeholder}
          type={props.type}
          value={props.value}
          onChange={props.onChange}
          isInvalid={props.error !== undefined && props.error.length > 0}
        />
        <Form.Control.Feedback type="invalid">{props.error}</Form.Control.Feedback>
    </Col>
  </Form.Group>
)

const Login = props => (
  <Container className="mt-2">
    <h2>Login</h2>
    <Input name="Email" placeholder="email@example.com" value={props.email} onChange={props.changeEmail} error={props.errors.email}/>
    <Input name="Password" type="password" placeholder="Password" value={props.password} onChange={props.changePassword} error={props.errors.password}/>
    <Button onClick={props.submit} disabled={!props.isValid}>Login</Button>
    <p className="mt-4">Don't have an account? <Link to="/">Signup</Link></p>
  </Container>
)

export default Login;
