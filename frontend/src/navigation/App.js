import React from 'react';
import { BrowserRouter, Route, Switch } from 'react-router-dom';
import PrimaryLayout from './PrimaryLayout.js';
import './App.css';
import Signup from '../portal/Signup.js';

function App() {
  return (
    <BrowserRouter>
      <link
        rel="stylesheet"
        href="https://maxcdn.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css"
        integrity="sha384-ggOyR0iXCbMQv3Xipma34MD+dH/1fQ784/j6cY/iJTQUOhcWr7x9JvoRxT2MZw1T"
        crossorigin="anonymous"
      />
      <Switch>
        {/* <Route path="/auth" component="" /> */}
        <Route path="/" exact component={Signup} />
        <Route path="/app" component={PrimaryLayout} />
      </Switch>
    </BrowserRouter>
  );
}
 
export default App;
