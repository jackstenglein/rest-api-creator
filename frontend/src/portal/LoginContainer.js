import React from 'react';
import Login from './Login';
import produce from 'immer';
import { login } from '../api/api.js';
import * as userInfo from '../redux/modules/userInfo.js';
import { connect } from 'react-redux';
import { Redirect } from 'react-router-dom';

// TODO: add loading state
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

  async login() {
    var response = await login(this.state.email, this.state.password);
    if (response === undefined) {
      this.setState(produce(this.state, draft => {draft.error = "Unable to contact server."}));
    }  else if (response.error !== undefined) {
      this.setState(produce(this.state, draft => {draft.error = response.error}));
    } else {
      this.props.loginSuccess(this.state.email);
    }
  }

  render() {
    if (this.props.authenticated) {
      return (<Redirect to="/app/details"/>)
    }

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

const mapStateToProps = state => {
  return {
    authenticated: state.userInfo.authenticated
  };
}

const mapDispatchToProps = dispatch => {
  return {
    loginSuccess: email => {
      dispatch(userInfo.loginSuccess(email))
    }
  };
}

export default connect(mapStateToProps, mapDispatchToProps)(LoginContainer);
