import React from 'react';
import Alert from 'react-bootstrap/Alert';


const ErrorAlert = props => {
  if (props.error === undefined || props.error.length === 0) {
    return null;
  }
  return (
    <Alert variant="danger">{props.error}</Alert>
  )
}

export default ErrorAlert;
