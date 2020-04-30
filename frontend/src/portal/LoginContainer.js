import React from 'react';
import Login from './Login';
import produce from 'immer';

class LoginContainer extends React.Component {

  constructor(props) {
    super(props);
    this.state = {
      email: "",
      password: ""
    };
    this.changeEmail = this.changeEmail.bind(this);
    this.changePassword = this.changePassword.bind(this);
    this.login = this.login.bind(this);
  }

  changeEmail(event) {
    const nextState = produce(this.state, draftState => {
      draftState.email = event.target.value
    });
    this.setState(nextState);
  }

  changePassword(event) {
    const nextState = produce(this.state, draftState => {
      draftState.password = event.target.value
    });
    this.setState(nextState);
  }

  login() {
    console.log(this.state);
  }

  render() {
    const isValid = this.state.email.length > 0 && this.state.password.length > 0;
    return (
      <Login 
        {...this.state}
        isValid={isValid}
        changeEmail={this.changeEmail}
        changePassword={this.changePassword}
        login={this.login}
      />
    )
  }
}

export default LoginContainer;
