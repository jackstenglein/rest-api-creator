import React from 'react';
import Login from './Login';

class LoginContainer extends React.Component {

  render() {
    return (
      <Login 
        email="test"
        password="asdf"
        errors={{}}
      />
    )
  }
}

export default LoginContainer;
