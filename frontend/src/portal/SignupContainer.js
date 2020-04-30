import React from 'react';
import Signup from './Signup';
import produce from 'immer';


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
      email: "test",
      password: "asf",
      confirmPassword: ""
    };

    this.changeEmail = this.changeEmail.bind(this);
    this.changePassword = this.changePassword.bind(this);
    this.changeConfirmPassword = this.changeConfirmPassword.bind(this);
    this.submit = this.submit.bind(this);
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

  submit() {
    console.log(this.state);
  }

  render() {
    const [isValid, errors] = validateSignup(this.state.email, this.state.password, this.state.confirmPassword);
    return (
      <Signup 
        {...this.state}
        isValid={isValid}
        errors={errors} 
        changeEmail={this.changeEmail} 
        changePassword={this.changePassword} 
        changeConfirmPassword={this.changeConfirmPassword} 
        submit={this.submit}
      />
    )
  }
}

export default SignupContainer;
