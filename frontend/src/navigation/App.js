import React from 'react';
import { BrowserRouter, Switch } from 'react-router-dom';
import PrimaryLayout from './PrimaryLayout.js';
import './App.css';
import LoginContainer from '../portal/LoginContainer.js';
import SignupContainer from '../portal/SignupContainer.js';
import AuthorizedRoute from './AuthorizedRoute.js';

function App() {
  return (
    <BrowserRouter>
      <link
        rel="stylesheet"
        href="https://maxcdn.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css"
        integrity="sha384-ggOyR0iXCbMQv3Xipma34MD+dH/1fQ784/j6cY/iJTQUOhcWr7x9JvoRxT2MZw1T"
        crossOrigin="anonymous"
      />
      <Switch>
        <AuthorizedRoute path="/" exact shouldAuthenticate={false} component={SignupContainer} />
        <AuthorizedRoute path="/login" shouldAuthenticate={false} component={LoginContainer} />
        <AuthorizedRoute path="/app" shouldAuthenticate={true} component={PrimaryLayout} />
      </Switch>
    </BrowserRouter>
  );
}
 
export default App;
