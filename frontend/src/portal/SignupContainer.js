import React from 'react';
import Signup from './Signup';
import produce from 'immer';
import { connect } from 'react-redux';
import * as userInfo from '../redux/modules/userInfo.js';
import { signup } from '../api/api.js';
import { Redirect } from 'react-router-dom';

function validateSignup(email, password, confirmPassword) {
  var isValid = true;
  const errors = {};
  if (!email.includes("@") || !email.includes(".")) {
    isValid = false;
    errors.email = "Invalid email";
  }
  if (password.length < 8) {
    isValid = false;
    errors.password = "Password must be 8 or more characters.";
  }
  if (confirmPassword !== password) {
    isValid = false;
    errors.confirmPassword = "Passwords do not match.";
  }
  return [isValid, errors];
}


class SignupContainer extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      email: "",
      password: "",
      confirmPassword: ""
    };

    this.changeEmail = this.changeEmail.bind(this);
    this.changePassword = this.changePassword.bind(this);
    this.changeConfirmPassword = this.changeConfirmPassword.bind(this);
    this.signup = this.signup.bind(this);
  }

  changeEmail(event) {
    const nextState = produce(this.state, draftState => {
      draftState.email = event.target.value;
    })
    this.setState(nextState);
  }

  changePassword(event) {
    const nextState = produce(this.state, draftState => {
      draftState.password = event.target.value;
    })
    this.setState(nextState);
  }

  changeConfirmPassword(event) {
    const nextState = produce(this.state, draftState => {
      draftState.confirmPassword = event.target.value;
    })
    this.setState(nextState);
  }

  async signup() {
    var response = await signup(this.state.email, this.state.password);
    if (response === undefined) {
      this.setState(produce(this.state, draft => {draft.submitError = "Unable to contact server."}));
    }  else if (response.error !== undefined) {
      this.setState(produce(this.state, draft => {draft.submitError = response.error}));
    } else {
      this.props.signupSuccess(this.state.email);
    }
  }

  render() {
    if (this.props.authenticated) {
      return (<Redirect to="/app/details"/>)
    }

    const [isValid, errors] = validateSignup(this.state.email, this.state.password, this.state.confirmPassword);
    return (
      <Signup 
        {...this.state}
        isValid={isValid}
        errors={errors} 
        changeEmail={this.changeEmail} 
        changePassword={this.changePassword} 
        changeConfirmPassword={this.changeConfirmPassword} 
        submit={this.signup}
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
    signupSuccess: email => {
      dispatch(userInfo.loginSuccess(email))
    }
  };
}

export default connect(mapStateToProps, mapDispatchToProps)(SignupContainer);
